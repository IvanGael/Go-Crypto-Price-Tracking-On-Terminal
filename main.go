package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	// "golang.org/x/term"
	"github.com/gorilla/websocket"
	"github.com/guptarohit/asciigraph"
	"github.com/rivo/tview"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

const (
	streamURL    = "wss://stream.binance.com/stream"
	graphCutSize = 10
)

func main() {
	width := 600
	graphSize := width - graphCutSize
	dataGraph := make([]float64, 0, graphSize)
	redColor := "red"
	greenColor := "green"
	price := float64(0)
	priceUp := true
	graph := ""
	var symbol, symbolBase string
	p := message.NewPrinter(language.Swedish)

	flag.StringVar(&symbol, "symbol", "btc", "symbol")
	flag.StringVar(&symbolBase, "symbolbase", "usdt", "symbol")
	flag.Parse()
	symbolText := fmt.Sprintf("%s / %s", symbol, symbolBase)

	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetWrap(false)

	done := make(chan struct{})
	data := make(chan float64, graphSize)

	go listen(symbol, symbolBase, done, data)

	go func() {
		for {
			select {
			case <-done:
				return
			case datum := <-data:
				if len(dataGraph) > 0 {
					price = datum
					priceUp = dataGraph[len(dataGraph)-1] < price
				}

				dataGraph = appendGraph(graphSize, dataGraph, datum)
				graph = asciigraph.Plot(dataGraph, asciigraph.Precision(5), asciigraph.Height(5))
				app.QueueUpdateDraw(func() {
					colorText := redColor
					if priceUp {
						colorText = greenColor
					}
					textView.SetText(p.Sprintf("   %s: [%s]%f[default]\n\n%v", symbolText, colorText, price, graph))
				})
			}
		}
	}()

	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		done <- struct{}{}
		panic(err)
	}
}

func appendGraph(size int, data []float64, newItem float64) []float64 {
	if len(data) < size {
		return append(data, newItem)
	}
	copy(data, data[1:])
	data[size-1] = newItem
	return data
}

type ReceiveMsg struct {
	Stream string `json:"stream"`
	Data   struct {
		P string `json:"p"`
	} `json:"data"`
}

func listen(symbol, symbolBase string, done chan struct{}, data chan float64) {
	u, err := url.Parse(streamURL)
	if err != nil {
		log.Fatal("Can't parse ", u)
	}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				panic(err)
			}
			var msg ReceiveMsg
			if err := json.Unmarshal(message, &msg); err != nil {
				panic(err)
			}

			p, err := strconv.ParseFloat(msg.Data.P, 64)
			if err != nil {
				continue
			}
			data <- p
		}
	}()

	err = c.WriteMessage(websocket.TextMessage, []byte(`{"method": "SUBSCRIBE","params": ["`+symbol+symbolBase+`@aggTrade"],"id": 1}`))
	if err != nil {
		log.Println("write:", err)
		return
	}

	for {
		select {
		case <-done:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-time.After(time.Second):
			}
			fmt.Println("Stop done")
			return
		}
	}
}
