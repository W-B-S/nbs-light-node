package main

import (
	"github.com/whyrusleeping/go-logging"
	"os"
	"path/filepath"
	"os/user"
	"fmt"
	"context"
	"github.com/W-B-S/nbs-light-node/core"
)

const (logFileName = "NBS-Server.log")

var log = logging.MustGetLogger("nbs/light-node")

var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func initLogFile() error{

	pathToLogFile := logFileName
	pathToHomeDir := ""
	pathToAppDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil{
		return err
	}

	if isPlatformDarwin{
		usr, err := user.Current()
		if err != nil{
			return err
		}
		pathToHomeDir = usr.HomeDir
		pathToAppFolder := pathToHomeDir +"/.nbs"
		os.Mkdir(pathToAppFolder, os.ModePerm)
		pathToLogFile = pathToAppFolder + "/" + logFileName

	} else if isPlatformLinux{
		pathToLogFile = pathToAppDirectory + "/" + logFileName
	}

	logFile, err := os.OpenFile(pathToLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	backend1 := logging.NewLogBackend(logFile, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(logging.INFO, "")

	logging.SetBackend(backend1Leveled)

	return nil
}


func createLightNode(ctx context.Context) (*core.NbsLightNode, error){

	node, err := core.NewLightNode(ctx)
	if err != nil{
		return nil, err
	}

	return node, nil
}

func main()  {

	if err := initLogFile(); err != nil{
		fmt.Errorf("---nbs-light-node---:failed to initLogFile:%s", err)
		return
	}

	node, err := createLightNode(context.Background())
	if err != nil{
		log.Error("---Failed to setup light node---,error:%s", err)
	}

	node.Run()
}