# bor-lab

A small, self-contained simulation of a **cosmos-sdk-style chain node**, written
in Go. It is the `bor-lab` repo of the Devcon offsec lab, and it is a nod to the
0xPolygon `bor` execution client.

The point of this repo is to look and behave like a real node codebase — modules,
a keeper, an application object, a begin/end block loop — without the weight of
importing the real Cosmos SDK. The SDK is expensive to build and to pull, so the
few pieces this simulation needs are hand-written stand-ins.

## simapp

This module simulates a cosmos-sdk based rollup node. The `app` package is the
`simapp` equivalent: it owns the store, registers modules, and drives the
ABCI-style lifecycle.

```
app.NewApp()            # build the app, register modules
app.BeginBlock(height)  # run every module's begin blocker
app.EndBlock(height)    # run every module's end blocker
```

Modules implement the `app.Module` interface:

```go
type Module interface {
	Name() string
	BeginBlock(height int64) error
	EndBlock(height int64) error
}
```

## Layout

```
cmd/borlabd/          # the node daemon (cobra command tree)
app/                  # application object, module registry, block loop
store/                # KVStore interface + in-memory implementation
x/transfer/           # the transfer module: keeper + MsgSend handler
x/transfer/types/     # MsgSend, Coin, and their validation
```

## Build and test

```sh
make build   # go build -o bin/borlabd ./cmd/borlabd
make test    # go test ./...
make lint    # go vet ./...
```

Or directly:

```sh
go build ./...
go test ./...
```

## Running the simulation node

```sh
go run ./cmd/borlabd version
go run ./cmd/borlabd start --blocks 10
go run ./cmd/borlabd tx send --from alice --to bob --amount 100ubor
```

`start` runs the app in simulation mode: it seeds two accounts, then produces
blocks while draining a simulated mempool.

## Dependencies

The only third-party dependency is `github.com/spf13/cobra` for the command
tree. Everything else is the standard library. `go.sum` is committed so builds
are reproducible offline after the module cache is populated.
