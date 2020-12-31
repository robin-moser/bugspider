package processor

import (
	"encoding/csv"
	"os"
)

func appendToFile(slice []string, filepath string) error {

	f, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	csvwriter := csv.NewWriter(f)
	err = csvwriter.Write(slice)
	if err != nil {
		return err
	}

	csvwriter.Flush()
	return nil

}
