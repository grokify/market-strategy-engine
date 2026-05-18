# Installation

## Using Go Install

The easiest way to install Market Strategy Engine:

```bash
go install github.com/grokify/market-strategy-engine/cmd/mse@latest
```

This installs the `mse` binary to your `$GOPATH/bin` directory.

## Building from Source

Clone and build manually:

```bash
git clone https://github.com/grokify/market-strategy-engine.git
cd market-strategy-engine
go build -o mse ./cmd/mse
```

Move the binary to your PATH:

```bash
mv mse /usr/local/bin/
```

## Verify Installation

Check the installation:

```bash
mse version
```

## As a Library

Add to your Go project:

```bash
go get github.com/grokify/market-strategy-engine
```

Then import the packages you need:

```go
import (
    "github.com/grokify/market-strategy-engine/model"
    "github.com/grokify/market-strategy-engine/engine"
    "github.com/grokify/market-strategy-engine/validate"
    "github.com/grokify/market-strategy-engine/report"
)
```
