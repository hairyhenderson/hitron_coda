![Build](https://github.com/hairyhenderson/hitron_coda/workflows/Build/badge.svg)

# Hitron CODA-4x8x Client

A Go client for the [Hitron CODA-4x8x](http://hitron-americas.com/products/service-providers/coda-4680-cable-modem-router/) DOCSIS 3.1 cable modem/router series.

This is tested on a Hitron CODA-4680 with firmware `7.1.1.2.2b9`, untested on
other models and releases.

The goal is to be able to perform the same actions through this client that are
available in the web UI that ships with the device.

## Status

This project is in active development, and not all APIs are supported yet.

The code is generated based on [`apilist.yaml`](./apilist.yaml), by running
`go generate`.

## Usage

```go
host := "192.168.0.1"
username := "cusadmin"
password := "mypassword"

// Instantiate a new *CableModem
cm, _ := New(host, username, password)

ctx := context.Background()

// Now login
_ = cm.Login(ctx)

// Now that we're logged in we can call APIs
info, _ := cm.RouterSysInfo(ctx)

fmt.Printf("Private LAN IP: %s\n", info.PrivLanIP)

// Output: Private LAN IP: 192.168.0.1
```

## License

[The MIT License](http://opensource.org/licenses/MIT)

Copyright (c) 2020-2021 Dave Henderson
