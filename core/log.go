package core

import (
	"github.com/W-B-S/nbs-light-node/utils"
	"os"
	"github.com/whyrusleeping/go-logging"
	"fmt"
)
var format = logging.MustStringFormatter(
	"%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}",
)

func init()  {
	if err := initLogs(); err != nil{
		fmt.Printf("---nbs-light-node---:failed to initLogFile:%s", err)
	}
}

func initLogs() error{

	pathToLogFile := utils.GetHomeDir() + LogFileName

	logFile, err := os.OpenFile(pathToLogFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}

	backend1 := logging.NewLogBackend(logFile, "", 0)

	backend2Formatter := logging.NewBackendFormatter(backend1, format)
	backend1Leveled := logging.AddModuleLevel(backend2Formatter)
	backend1Leveled.SetLevel(logging.DEBUG, "")

	logging.SetBackend(backend1Leveled)

	return nil
}
