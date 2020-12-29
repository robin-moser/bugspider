package host

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"strings"
)

type HostListProcessor struct {
	hostFile string
}

func (processor *HostListProcessor) DoProcess(host *Host) (bool, error) {

	hostFile := path.Join(".", processor.hostFile)
	path := path.Dir(hostFile)

	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		os.MkdirAll(path, os.FileMode(0755))
	}

	f, err := os.OpenFile(hostFile, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return false, err
	}
	defer f.Close()

	alreadyInFile, err := alreadyInFile(host.Hostname, hostFile)
	if err != nil {
		return false, err
	} else if !alreadyInFile {
		appendToHostFile(host, *f)
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

func appendToHostFile(host *Host, f os.File) {
	var hostArray []string

	hostArray = append(hostArray, host.Hostname)
	hostArray = append(hostArray, host.Source)
	hostArray = append(hostArray, host.Date.String())

	csvwriter := csv.NewWriter(&f)
	_ = csvwriter.Write(hostArray)
	csvwriter.Flush()
}

func MakeNewHostProcessor(hostFile string) *HostListProcessor {
	return &HostListProcessor{hostFile: hostFile}
}
