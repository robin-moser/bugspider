package processor

import (
	"fmt"
	"log"
	"path"
	"time"

	"go.etcd.io/bbolt"
)

const Database string = "output/hosts.db"

var dbPath string = path.Join(".", Database)
var csvHostPath string = path.Dir(dbPath)

var db *bbolt.DB

func InitDB() {
	var err error
	db, err = bbolt.Open(dbPath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDB() {
	db.Close()
}

func ProcessDeduplication(currentHost *Host) (bool, error) {
	alreadyInDB, err := alreadyInDB(currentHost.Hostname)

	if err != nil {
		log.Fatal(err)
		return false, err
	}
	if !alreadyInDB {

		var hostArray []string

		hostArray = append(hostArray, currentHost.Hostname)
		hostArray = append(hostArray, currentHost.Source)
		hostArray = append(hostArray, currentHost.Date.Format("2006-01-02 15:04:05"))

		err := appendToDB(hostArray)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

func alreadyInDB(hostname string) (bool, error) {
	err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("Hosts"))
		if bucket == nil {
			return nil
		}
		val := bucket.Get([]byte(hostname))
		if val != nil {
			return fmt.Errorf("Hostname already exists")
		}
		return nil
	})
	if err != nil {
		return true, nil
	}
	return false, nil
}

func appendToDB(hostArray []string) error {
	err := db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("Hosts"))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(hostArray[0]), []byte(hostArray[1]+","+hostArray[2]))
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
