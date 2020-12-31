package processor

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/robin-moser/bugspider/request"
)

const gitConfigDir string = "output/configs"
const csvGitLogFile string = "output/opengit.csv"
const csvGitErrorFile string = "output/errors.csv"

func ProcessOpenGit(currentHost *Host) (bool, int, error) {

	success, status, err := openGitRequest(currentHost)

	var hostArray []string

	hostArray = append(hostArray, currentHost.Hostname)
	hostArray = append(hostArray, strconv.FormatBool(success))
	hostArray = append(hostArray, strconv.Itoa(status))
	hostArray = append(hostArray, time.Now().Format("2006-01-02 15:04:05"))

	appendToFile(hostArray, csvGitLogFile)

	if err != nil {
		hostArray = append(hostArray, fmt.Sprintf("%v", err))
		appendToFile(hostArray, csvGitErrorFile)
	}

	return success, status, err

}

func openGitRequest(currentHost *Host) (bool, int, error) {

	configPath := path.Join(".", gitConfigDir)
	hostFile := path.Join(configPath, currentHost.Hostname)

	url := "https://" + currentHost.Hostname + "/.git/config"
	body, status, err := request.GetResponseBody(url, true)
	if err != nil {
		return false, 0, err
	}

	if strings.Contains(string(body), "[core]") {

		// create the specified dir structure if not existant
		_, err := os.Stat(configPath)
		if os.IsNotExist(err) {
			log.Println("Creating directory ", configPath)
			os.MkdirAll(configPath, os.FileMode(0755))
		}

		file, err := os.OpenFile(hostFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			return false, status, err
		}
		defer file.Close()

		_, err = file.Write(body)
		if err != nil {
			return false, status, err
		}

		return true, status, nil

	}

	return false, status, nil

}
