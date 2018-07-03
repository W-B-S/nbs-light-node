package core

import (
	"context"
	p2pHost "gx/ipfs/QmQQGtcp6nVUrQjNsnU53YWV1q8fK1Kd9S7FEkYbRZzxry/go-libp2p-host"
	"gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	ic "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
	"gx/ipfs/QmRg1gKTHzc3CZXSKzem8aR4E3TubFhbgXwfVuWnSK5CC5/go-metrics-interface"
	"gx/ipfs/QmRFEBGcNjtWPupwHA7zGHeGVLuUyE4ZRFi2MgtrPM6pfb/go-libp2p-floodsub"
	peerStore "gx/ipfs/QmZhsmorLpD9kmQ4ynbAu4vbKv2goMUnXazwGA4gnWHDjB/go-libp2p-peerstore"
	"gx/ipfs/QmUEAR2pS7fP1GPseS3i8MWFyENs7oDp4CZrgn8FCjbsBu/go-libp2p"
	"github.com/whyrusleeping/go-logging"
	"fmt"
	"github.com/ipfs/go-ipfs/p2p"
	ma "gx/ipfs/QmUxSEGbv2nmYNnfXi7839wwQqTN3kwQeUxe8dTjZWZs7J/go-multiaddr"
	"sync"
	"encoding/base64"
)

type NbsLightNode struct {
	ctx  		context.Context
	identity 	peer.ID

	privateKey  ic.PrivKey
	peerStore  	peerStore.Peerstore
	peerHost    p2pHost.Host
	floodSub 	*floodsub.PubSub
	P2P      	*p2p.P2P
}

var log = logging.MustGetLogger("nbs/core")
var once sync.Once
var instanceNode *NbsLightNode

func GetNodeInstance(ctx context.Context) (*NbsLightNode){

	once.Do(func() {

		node, err := newNode(ctx)
		if err != nil{
			panic(err)
		}

		instanceNode = node
	})

	return instanceNode
}

func newNode(ctx context.Context) (*NbsLightNode, error) {

	ctx = metrics.CtxScope(ctx, "nbs-light-node")

	node := &NbsLightNode{
		ctx:	ctx,
	}

	if err := setupNode(ctx, node); err != nil{
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

func setupNode(ctx context.Context, node *NbsLightNode) error{

	peerId, err := GetSysConfig().LoadId()
	if err != nil{
		log.Error("Failed to load node's identity", err)
		return err
	}
	node.identity = peerId

	node.peerStore = peerStore.NewPeerstore()

	if err := node.loadPrivateKey(); err != nil{
		return err
	}

	if err := setupFloodSub(ctx, node); err != nil{
		return err
	}

	if err := startListening(node.peerHost); err != nil {
		return err
	}

	node.P2P = p2p.NewP2P(node.identity, node.peerHost, node.peerStore)

	return nil
}

func (node *NbsLightNode) loadPrivateKey() error {

	pkb, err := base64.StdEncoding.DecodeString(GetSysConfig().PrivateKey)
	if err != nil {
		return err
	}

	//TODO:: use password to decode private key

	sk, err := ic.UnmarshalPrivateKey(pkb)
	if err != nil {
		return err
	}

	id2, err := peer.IDFromPrivateKey(sk)
	if err != nil {
		return err
	}

	if id2 != node.identity {
		return fmt.Errorf("private key in config does not match id: %s != %s", node.identity, id2)
	}

	node.privateKey = sk
	node.peerStore.AddPrivKey(node.identity, node.privateKey)
	node.peerStore.AddPubKey(node.identity, sk.GetPublic())

	return nil
}

func setupFloodSub(ctx context.Context, node *NbsLightNode) error{

	//TODO:: no security conntections right now
	var options = []libp2p.Option{libp2p.NoSecurity, libp2p.Peerstore(node.peerStore)}

	privateKey :=  node.peerStore.PrivKey(node.identity)

	if privateKey == nil {
		return fmt.Errorf("missing private key for node ID: %s", node.identity.Pretty())
	}
	options = append(options, libp2p.Identity(privateKey))

	peerHost, err := libp2p.New(ctx, options...)

	service, err := floodsub.NewFloodSub(ctx, peerHost)
	if err != nil {
		return err
	}

	node.floodSub = service
	node.peerHost = peerHost

	return nil
}

func startListening(host p2pHost.Host) error{

	var listenAddress []ma.Multiaddr
	for _, addr := range GetSysConfig().Swarm {
		multiAddress, err := ma.NewMultiaddr(addr)
		if err != nil {
			return fmt.Errorf("failure to parse config.Addresses.Swarm: %s", addr)
		}
		listenAddress = append(listenAddress, multiAddress)
	}

	if err := host.Network().Listen(listenAddress...); err != nil {
		return err
	}

	address, err := host.Network().InterfaceListenAddresses()
	if err != nil {
		return err
	}

	log.Infof("Swarm listening at: %s", address)

	return nil
}
