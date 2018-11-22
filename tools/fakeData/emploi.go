package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func readAndRandomEmploi(fileName string, outputFileName string, mapping map[string]string) error {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	// ligne de titre
	row, err := reader.Read()
	outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return err
	}
	for {
		c := rand.Float64() * 2.0
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		if _, ok := mapping[row[95]]; ok {

			for i := 0; i < 93; i++ {
				if row[i] != "" {
					val, _ := strconv.ParseFloat(row[i], 64)
					row[i] = strconv.Itoa(int(val * c))
				}
			}
			row[93] = ""
			row[94] = ""
			row[96] = ""
			row[95] = mapping[row[95]]
			row[97] = mapping[row[97]]
			row[98] = ""
			outputRow := strings.Join(row, ";") + "\n"
			_, err = outputFile.WriteString(outputRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
