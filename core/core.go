package core

import (
	"context"
	"github.com/ipfs/go-ipfs/p2p"
  	peer 			"github.com/libp2p/go-libp2p-peer"
  	discovery 		"github.com/libp2p/go-libp2p/p2p/discovery"
	p2phost 		"github.com/libp2p/go-libp2p-host"
	floodsub 		"github.com/libp2p/go-floodsub"
	psrouter 		"github.com/libp2p/go-libp2p-pubsub-router"
	goprocess 		"github.com/jbenet/goprocess"
	ic 				"github.com/libp2p/go-libp2p-crypto"
	"github.com/ipfs/go-metrics-interface"
	"github.com/whyrusleeping/go-logging"
)

var log = logging.MustGetLogger("nbs/light-node")

type NbsLightNode struct {
	identity 	peer.ID
	privateKey  ic.PrivKey // the local node's private Key
	discovery  	discovery.Service

	peerHost    p2phost.Host
	floodSub 	*floodsub.PubSub
	PSRouter 	*psrouter.PubsubValueStore
	P2P      	*p2p.P2P

	proc 		goprocess.Process
	ctx  		context.Context
}

func NewLightNode(ctx context.Context) (*NbsLightNode, error) {
	ctx = metrics.CtxScope(ctx, "nbs-light-node")

	node := &NbsLightNode{
		ctx:	ctx,
	}

	if err := setUpNode(ctx, node); err != nil{
		return nil, err
	}

	log.Info("---Setup nbs light node success---")

	return node, nil
}

func (node *NbsLightNode) Run()  {
	log.Info("---Start running---")
	defer func() {
		select {
			case <-node.ctx.Done():
				log.Info("---Node finished---")
		default:
			log.Info("--------1--------")
		}
	}()
}

func setUpNode(ctx context.Context, node *NbsLightNode) error{
	return nil
}

