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

func readAndRandomCotisations(fileName string, outputFileName string, mapping map[string]string) error {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	coef := make(map[string]float64)

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
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		mer, _ := strconv.ParseFloat(row[2], 64)
		encDirect, _ := strconv.ParseFloat(row[3], 64)
		cotisDue, _ := strconv.ParseFloat(row[5], 64)

		if c, ok := coef[row[0]]; ok {
			row[2] = strconv.Itoa(int(mer * c))
			row[3] = strconv.Itoa(int(encDirect * c))
			row[5] = strconv.Itoa(int(cotisDue * c))
		} else {
			coef[row[0]] = rand.Float64() * rand.Float64() / 150
			row[2] = strconv.Itoa(int(mer * coef[row[0]]))
			row[3] = strconv.Itoa(int(encDirect * coef[row[0]]))
			row[5] = strconv.Itoa(int(cotisDue * coef[row[0]]))
		}

		row[0] = mapping[row[0]]

		outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
		_, err = outputFile.WriteString(outputRow)
		if err != nil {
			return err
		}
	}

	return nil
}
