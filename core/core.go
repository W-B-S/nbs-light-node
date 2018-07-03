package core

import (
	"context"
	"github.com/ipfs/go-ipfs/p2p"
	p2pHost "gx/ipfs/QmQQGtcp6nVUrQjNsnU53YWV1q8fK1Kd9S7FEkYbRZzxry/go-libp2p-host"
	"gx/ipfs/QmSF8fPo3jgVBAy8fpdjjYqgG87dkJgUprRBHRd2tmfgpP/goprocess"
	"gx/ipfs/QmUEAR2pS7fP1GPseS3i8MWFyENs7oDp4CZrgn8FCjbsBu/go-libp2p/p2p/discovery"
	pubRouter "gx/ipfs/QmVJvKzFi93cWDknZr1UUKyLTJWVRRSixqCqwQ9nLfF8ob/go-libp2p-pubsub-router"
	"gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	ic "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmRg1gKTHzc3CZXSKzem8aR4E3TubFhbgXwfVuWnSK5CC5/go-metrics-interface"
	"gx/ipfs/QmRFEBGcNjtWPupwHA7zGHeGVLuUyE4ZRFi2MgtrPM6pfb/go-libp2p-floodsub"
	peerStore "gx/ipfs/QmZhsmorLpD9kmQ4ynbAu4vbKv2goMUnXazwGA4gnWHDjB/go-libp2p-peerstore"
	"gx/ipfs/QmUEAR2pS7fP1GPseS3i8MWFyENs7oDp4CZrgn8FCjbsBu/go-libp2p"
	"github.com/whyrusleeping/go-logging"
	"fmt"
)

var log = logging.MustGetLogger("nbs/core")

type NbsLightNode struct {
	ctx  		context.Context
	identity 	peer.ID

	privateKey  ic.PrivKey
	discovery  	discovery.Service

	peerStore  	peerStore.Peerstore
	peerHost    p2pHost.Host
	floodSub 	*floodsub.PubSub
	PSRouter 	*pubRouter.PubsubValueStore
	P2P      	*p2p.P2P

	process		goprocess.Process
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

	peerId, err := GetSysConfig().LoadId()
	if err != nil{
		log.Error("Failed to load node's identity", err)
		return err
	}

	node.identity = peerId
	node.peerStore = peerStore.NewPeerstore()

	var options []libp2p.Option

	peerHost, err := constructPeerHost(ctx, node.identity, node.peerStore, options...)

	if err != nil {
		return err
	}

	service, err := floodsub.NewFloodSub(ctx, peerHost)
	if err != nil {
		return err
	}

	node.floodSub = service

	return nil
}

func constructPeerHost(ctx context.Context, id peer.ID, ps peerStore.Peerstore, options ...libp2p.Option) (p2pHost.Host, error) {
	privateKey := ps.PrivKey(id)
	if privateKey == nil {
		return nil, fmt.Errorf("missing private key for node ID: %s", id.Pretty())
	}
	options = append([]libp2p.Option{libp2p.Identity(privateKey), libp2p.Peerstore(ps)}, options...)
	return libp2p.New(ctx, options...)
}
