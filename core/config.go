package core

import (
	"github.com/W-B-S/nbs-light-node/utils"
	"github.com/whyrusleeping/go-logging"
	"os"
	"encoding/json"
	"io"
	"fmt"
	"github.com/W-B-S/nbs-light-node/errors"
	"gx/ipfs/QmVf8hTAsLLFtn4WPCRNdnaF2Eag2qTBS6uR8AiHPZARXy/go-libp2p-peer"
	"crypto/rand"
	"encoding/base64"
	ci "gx/ipfs/Qme1knMqwt1hKZbc1BmQFmnm9f36nyQGwXxPGVpVJ9rMK5/go-libp2p-crypto"
)

const (	LogFileName 	= "NBS-Light-Node.log"
		ConfigFileName 	= ".config")

var NbsLog = logging.MustGetLogger("core/config")

var SystemConfig NodeConfig


type NodeConfig struct {
	PeerID	string `json:"peerId"`
	PrivateKey string `json:",omitempty"`
	EncryptedPrivateKey string `json:"privateKey,omitempty"`
	Swarm      []string
}

func init()  {
	SystemConfig = initSystemConfig()
}

func GetSysConfig() (NodeConfig) {
	return SystemConfig
}

func (config NodeConfig) LoadId() (peer.ID, error){

	peerId := config.PeerID
	if peerId == "" {
		return "", errors.New("identity was not set in config")
	}

	if len(peerId) == 0 {
		return "", errors.New("no peer ID in config! (was 'ipfs init' run?)")
	}

	decodedId, err := peer.IDB58Decode(peerId)

	if err != nil {
		return "", fmt.Errorf("peer ID invalid: %s", err)
	}

	return decodedId, nil
}

func getDefaultPath() string{
	return utils.GetHomeDir() + ConfigFileName
}


func (config *NodeConfig)createDefaultIdentity() error{

	privateKey, publicKey, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, rand.Reader)
	if err != nil {
		return err
	}

	pid, err := peer.IDFromPublicKey(publicKey)
	if err != nil {
		return err
	}

	privateKeyBytes, err := privateKey.Bytes()
	if err != nil {
		return err
	}

	config.PeerID = pid.Pretty()
	config.PrivateKey = base64.StdEncoding.EncodeToString(privateKeyBytes)

	//TODO::	save the private key by encrypted data
	//config.EncryptedPrivateKey = ?

	config.Swarm =[]string{
		"/ip4/0.0.0.0/tcp/4001",
		"/ip6/::/tcp/4001",
	}

	return nil
}

func (config *NodeConfig) syncToRepo(path string) error{

	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		NbsLog.Error("Failed to create config file:%s", err)
		return err
	}

	defer  logFile.Close()

	jsonByte, err := json.Marshal(config)
	if err != nil{
		NbsLog.Error("Failed to marshal default config to disk:%s", err)
		return err
	}

	_, err = logFile.Write(jsonByte)
	if err != nil{
		NbsLog.Error("Failed to write json content to disk:%s", err)
		return err
	}

	return nil
}

func createDefaultConfig(path string) (NodeConfig, error) {

	config := NodeConfig{}

	config.createDefaultIdentity()

	err := config.syncToRepo(path)
	if err != nil{
		fmt.Errorf("failed to save config to disk!:%s", err.Error())
		return config, err
	}

	return config, nil
}


func loadFromRepo(path string) (NodeConfig){

	config := NodeConfig{}

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

	json.Unmarshal([]byte(jsonContent), config)

	return config
}

func initSystemConfig() (NodeConfig){

	path := getDefaultPath()

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {

			config, err := createDefaultConfig(path)
			if err != nil{
				panic(err)
			}
			return config

		} else {
			panic(err)
		}
	}else{

		return loadFromRepo(path)
	}
}