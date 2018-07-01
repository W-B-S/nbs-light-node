package utils

import (
	"path/filepath"
	"os"
	"os/user"
)

var applicationHomeDir = ""

func init()  {
	homeDir, err := appHomeDir()
	if err != nil{
		homeDir = ""
	}
	applicationHomeDir = homeDir
}

func GetHomeDir()  string{
	return applicationHomeDir
}

func appHomeDir() (string, error) {

	if isPlatformDarwin{
		usr, err := user.Current()
		if err != nil{
			return "", err
		}
		pathToHomeDir := usr.HomeDir
		pathToAppFolder := pathToHomeDir +"/.nbs"
		os.Mkdir(pathToAppFolder, os.ModePerm)
		return pathToAppFolder + "/", nil

	} else if isPlatformLinux{
		pathToAppDirectory, err := filepath.Abs(filepath.Dir(os.Args[0]))

		if err != nil{
			return "", err
		}

		return pathToAppDirectory + "/", nil
	}

	return "" , nil
}