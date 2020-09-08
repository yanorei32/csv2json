package main

import (
	"flag"
	"os"
	"log"
	"io"
	"encoding/json"
	"encoding/csv"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding/japanese"
)

func main() {
	flag.Parse()

	csvf, err := os.Open(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}
	defer csvf.Close()

	jsonf, err := os.Create(flag.Arg(1))
	if err != nil {
		log.Fatal(err)
	}
	defer jsonf.Close()

	reader := csv.NewReader(
		transform.NewReader(
			csvf,
			japanese.ShiftJIS.NewDecoder(),
		),
	)

	reader.LazyQuotes = true

	header, err := reader.Read()
	if err != nil {
		log.Fatal(err)
	}

	data := make([]map[string]string, 0, 0)
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			log.Fatal(err)
		}

		v := make(map[string]string, len(header))
		for i, col := range header {
			v[col] = record[i]
		}

		data = append(data, v)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := jsonf.Write(bytes); err != nil {
		log.Fatal(err)
	}
}

