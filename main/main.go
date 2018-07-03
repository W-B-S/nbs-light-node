package main

import (
	"context"
	"github.com/W-B-S/nbs-light-node/core"
	"github.com/whyrusleeping/go-logging"
)

var NbsLog = logging.MustGetLogger("nbs/light-node")

func main()  {

	node := core.GetNodeInstance(context.Background())

	node.Run()
}