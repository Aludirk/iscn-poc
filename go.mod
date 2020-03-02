module github.com/likecoin/iscn-poc

go 1.13

replace github.com/ipfs/go-ipfs => ./go-ipfs

replace github.com/likecoin/iscn-ipld => ./go-ipfs/plugin/plugins/iscn-ipld

require (
	github.com/cosmos/cosmos-sdk v0.38.1
	github.com/ipfs/go-cid v0.0.4
	github.com/ipfs/go-ipfs v0.4.23
	github.com/ipfs/go-ipfs-config v0.0.3
	github.com/ipfs/interface-go-ipfs-core v0.0.8
	github.com/likecoin/iscn-ipld v0.0.0-00010101000000-000000000000
	github.com/tendermint/tendermint v0.33.0
	go.uber.org/dig v1.8.0 // indirect
	go.uber.org/multierr v1.4.0 // indirect
)
