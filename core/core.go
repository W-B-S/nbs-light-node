package core

import (
	"context"
	p2p "github.com/ipfs/go-ipfs/p2p"
	p2phost "gx/ipfs/QmQQGtcp6nVUrQjNsnU53YWV1q8fK1Kd9S7FEkYbRZzxry/go-libp2p-host"
	goprocess "gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"
	discovery "gx/ipfs/QmUEAR2pS7fP1GPseS3i8MWFyENs7oDp4CZrgn8FCjbsBu/go-libp2p/p2p/discovery"
	psrouter "gx/ipfs/QmVJvKzFi93cWDknZr1UUKyLTJWVRRSixqCqwQ9nLfF8ob/go-libp2p-pubsub-router"
	peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	ic "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	metrics "gx/ipfs/QmRg1gKTHzc3CZXSKzem8aR4E3TubFhbgXwfVuWnSK5CC5/go-metrics-interface"
	floodsub "gx/ipfs/QmRFEBGcNjtWPupwHA7zGHeGVLuUyE4ZRFi2MgtrPM6pfb/go-libp2p-floodsub"
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
		}
	}()
}

func setUpNode(ctx context.Context, node *NbsLightNode) error{



	return nil
}

