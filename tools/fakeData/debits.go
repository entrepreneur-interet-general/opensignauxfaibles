package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func readAndRandomDebits(fileName string, outputFileName string, mapping map[string]string) error {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'

	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// ligne de titre
	row, err := reader.Read()
	outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return err
	}

	// map des coefficients générés
	coef := make(map[string]float64)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		partOuvriere, err := strconv.ParseFloat(row[4], 64)
		if err != nil {
			return err
		}
		partPatronale, _ := strconv.ParseFloat(row[5], 64)
		if err != nil {
			return err
		}
		partPenalite, _ := strconv.ParseFloat(row[15], 64)
		if err != nil {
			return err
		}

		if c, ok := coef[row[0]]; ok {
			row[4] = strconv.Itoa(int(partOuvriere * c))
			row[5] = strconv.Itoa(int(partPatronale * c))
			row[15] = strconv.Itoa(int(partPenalite * c))
		} else {
			coef[row[0]] = rand.Float64() * rand.Float64() / 150
			row[4] = strconv.Itoa(int(partOuvriere * coef[row[0]]))
			row[5] = strconv.Itoa(int(partPatronale * coef[row[0]]))
			row[15] = strconv.Itoa(int(partPenalite * coef[row[0]]))
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
