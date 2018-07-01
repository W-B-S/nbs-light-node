package core

import "github.com/W-B-S/nbs-light-node/utils"

const (LogFileName = "NBS-Light-Node.log")

type NodeConfig struct {
	PeerID	string `json:"peerid,omitempty"`
}

var DefaultConfig = NodeConfig{
	PeerID:"",
}

func (config *NodeConfig)getDefaultPath() string{

	return utils.GetHomeDir() + ".config"
}

func (config *NodeConfig)LoadFromDisk(){
	//path := config.getDefaultPath()

}