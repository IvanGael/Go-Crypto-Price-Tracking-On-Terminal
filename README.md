<h1 align="center">Go - Crypto coin price tracking on terminal.</h1>

<div align="center">
  <p><img src="https://storage.googleapis.com/buildship-vos7yw-europe-west1/uploaded-files/cpdp.gif"></p>
</div>



# Install

**Source**

1. Requirement

```
Go 1.16
```

2. Build binary

```
go get;
go build;
```

3. Run the build

```
main
```


# Usage

Type from commandline

```
gohowmuch --symbol=doge
```

Support *--symbol* and *--symbolbase*.

*--symbol* default value is "btc".

*--symbolbase* default value is "usdt".

Example:

```
gohowmuch --symbol=btc --symbolbase=usdt
gohowmuch --symbol=doge --symbolbase=btc
gohowmuch --symbol=shib --symbolbase=doge
```
