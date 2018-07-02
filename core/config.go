package core

import (
	"github.com/W-B-S/nbs-light-node/utils"
	"github.com/whyrusleeping/go-logging"
	"os"
	"encoding/json"
	"io"
	"fmt"
	"github.com/W-B-S/nbs-light-node/errors"
	peer "gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
)

const (	LogFileName 	= "NBS-Light-Node.log"
		ConfigFileName 	= ".config")

var NbsLog = logging.MustGetLogger("core/config")

var SystemConfig *NodeConfig


type NodeConfig struct {
	peerID	string `json:"peerid, omitempty"`
	privKey string `json:",omitempty"`
}

var defaultConfig = NodeConfig{
	peerID:"",
}

func init()  {
	SystemConfig = LoadFromDisk()
}

func GetSysConfig() (NodeConfig) {
	return *SystemConfig
}

func (config NodeConfig) LoadId() (peer.ID, error){

	cid := config.peerID
	if cid == "" {
		return "", errors.New("identity was not set in config")
	}

	if len(cid) == 0 {
		return "", errors.New("no peer ID in config! (was 'ipfs init' run?)")
	}

	id, err := peer.IDB58Decode(cid)

	if err != nil {
		return "", fmt.Errorf("peer ID invalid: %s", err)
	}

	return id, nil
}

func getDefaultPath() string{
	return utils.GetHomeDir() + ConfigFileName
}

//TODO::
func initSystemConfig(){

	defaultConfig.peerID = "100132"
}

func createDefaultConfig(path string)  {

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		NbsLog.Error("Failed to create config file:%s", err)
		return
	}

	defer  logFile.Close()

	initSystemConfig()

	jsonByte, err := json.Marshal(defaultConfig)
	if err != nil{
		NbsLog.Error("Failed to marshal default config to disk:%s", err)
		return
	}

	_, err = logFile.Write(jsonByte)
	if err != nil{
		NbsLog.Error("Failed to write json content to disk:%s", err)
		return
	}
}

func LoadFromDisk() (*NodeConfig){

	path := getDefaultPath()

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {

			createDefaultConfig(path)

			return &defaultConfig

		} else {
			panic(err)
		}

	}else{

		logFile, err := os.OpenFile(path, os.O_RDONLY, 0600)
		if err != nil{
			panic(err)
		}
		defer logFile.Close()

		jsonContent := ""
		buffer := make([]byte, 1024)

		for {
			number, err := logFile.Read(buffer)
			if err != nil && err != io.EOF{
				panic(err)
			}
			if 0 == number{
				break
			}

			jsonContent = jsonContent + string(buffer[:number])
		}

		config := &NodeConfig{}
		json.Unmarshal([]byte(jsonContent), config)

		return config
	}
}