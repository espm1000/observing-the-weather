package report

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

type Dogs struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func CsvHandler(filename string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := parseCsv(filename)
		if err != nil {
			log.Printf("error parsing %s: %v", filename, err)
			return
		}
		fmt.Fprint(w, string(data))
	}
}

func parseCsv(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println("error closing file: ", err)
		}
	}()

	data, err := os.ReadFile(file.Name())
	if err != nil {
		log.Printf("error reading file %s: %v", file.Name(), err)
		return nil, err
	}

	return data, nil
}
