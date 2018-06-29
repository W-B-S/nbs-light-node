package core

import (
	"context"
	"github.com/ipfs/go-ipfs/p2p"
  	peer "github.com/libp2p/go-libp2p-peer"
  	discovery "github.com/libp2p/go-libp2p/p2p/discovery"
	p2phost "github.com/libp2p/go-libp2p-host"
	//floodsub "github.com/libp2p/go-floodsub"
	floodsub "gx/ipfs/QmRFEBGcNjtWPupwHA7zGHeGVLuUyE4ZRFi2MgtrPM6pfb/go-libp2p-floodsub"
	psrouter "gx/ipfs/QmVJvKzFi93cWDknZr1UUKyLTJWVRRSixqCqwQ9nLfF8ob/go-libp2p-pubsub-router"
	goprocess "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"
	)


type NbsLightNode struct {
	Identity peer.ID

	Discovery  	discovery.Service

	PeerHost    p2phost.Host
	Floodsub 	*floodsub.PubSub
	PSRouter 	*psrouter.PubsubValueStore
	P2P      	*p2p.P2P

	proc 		goprocess.Process
	ctx  		context.Context
}