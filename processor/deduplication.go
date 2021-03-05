package processor

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
)

const HostFile string = "output/hosts.csv"

var csvHostFile string = path.Join(".", HostFile)
var csvHostPath string = path.Dir(csvHostFile)

func ProcessDeduplication(currentHost *Host) (bool, error) {

	// create the specified dir structure if not existant
	_, err := os.Stat(csvHostPath)
	if os.IsNotExist(err) {
		log.Println("Creating directory ", csvHostPath)
		os.MkdirAll(csvHostPath, os.FileMode(0755))
	}

	alreadyInFile, err := alreadyInFile(currentHost.Hostname, csvHostFile)

	if err != nil {
		log.Fatal("test - ", err)
		return false, err
	}
	if !alreadyInFile {

		var hostArray []string

		hostArray = append(hostArray, currentHost.Hostname)
		hostArray = append(hostArray, currentHost.Source)
		hostArray = append(hostArray, currentHost.Date.Format("2006-01-02 15:04:05"))

		err := appendToFile(hostArray, csvHostFile)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func alreadyInFile(hostname string, filepath string) (bool, error) {

	find := []byte("\n" + hostname + ",")
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false, nil
	}

	if bytes.Contains(append([]byte("\n"), dat...), find) {
		return true, nil
	}

	return false, nil
}
