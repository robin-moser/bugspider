package host

import (
	"bufio"
	"fmt"
	"os"
	"path"
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
	f.WriteString(host.Hostname + "\n")
}

func MakeNewHostProcessor(hostFile string) *HostListProcessor {
	return &HostListProcessor{hostFile: hostFile}
}
