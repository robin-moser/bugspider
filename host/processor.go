package host

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

type Processor struct {
	processor string
	hostFile  string
}

// Process validates the Host and stores it with the specified storage provider.
// Right now, the only possible storage provider is a csv file
func (processor *Processor) Process(host *Host) (bool, error) {

	hostFile := path.Join(".", processor.hostFile)
	path := path.Dir(hostFile)

	// create the specified dir structure if not existant
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("Creating directory ", path)
		os.MkdirAll(path, os.FileMode(0755))
	}

	// create the specified file if not existant
	_, err = os.Stat(hostFile)
	if os.IsNotExist(err) {
		log.Println("Creating file ", hostFile)
		os.Create(hostFile)
	}

	alreadyInFile, err := alreadyInFile(host.Hostname, hostFile)
	if err != nil {
		return false, err
	}
	if !alreadyInFile {
		err := appendToHostFile(host, hostFile)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func alreadyInFile(hostname string, hostFile string) (bool, error) {
	f, err := os.OpenFile(hostFile, os.O_RDONLY, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), hostname+",") {
			return true, nil
		}
		if scanner.Text() == hostname {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	return false, nil
}

func appendToHostFile(host *Host, hostFile string) error {

	var hostArray []string

	hostArray = append(hostArray, host.Hostname)
	hostArray = append(hostArray, host.Source)
	hostArray = append(hostArray, host.Date.String())

	f, err := os.OpenFile(hostFile, os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	err = csvwriter.Write(hostArray)
	if err != nil {
		return err
	}

	csvwriter.Flush()
	return nil
}

func NewProcessor(processor string, hostFile string) *Processor {
	return &Processor{
		processor: processor,
		hostFile:  hostFile,
	}
}
