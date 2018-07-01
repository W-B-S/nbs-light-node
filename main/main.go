package main

import (
	"context"
	"github.com/W-B-S/nbs-light-node/core"
	"github.com/whyrusleeping/go-logging"
)

var NbsLog = logging.MustGetLogger("nbs/light-node")

func createLightNode(ctx context.Context) (*core.NbsLightNode, error){

	node, err := core.NewLightNode(ctx)
	if err != nil{
		return nil, err
	}

	return node, nil
}

func main()  {

	node, err := createLightNode(context.Background())
	if err != nil{
		NbsLog.Error("---Failed to setup light node---,error:%s", err)
	}

	node.Run()
}