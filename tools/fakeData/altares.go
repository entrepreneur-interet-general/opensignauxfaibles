package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strings"
)

func readAndRandomAltares(fileName string, outputFileName string, mapping map[string]string) error {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','
	reader.LazyQuotes = true
	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// ligne de titre
	row, err := reader.Read()
	outputRow := strings.Join(row, ",") + "\n"
	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return err
	}

	idx := index(row)
	lidx := len(idx)

	for {
		row, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		row[0] = mapping[row[0]]
		output := make([]string, lidx)

		output[idx["Date d'effet"]] = row[idx["Date d'effet"]]
		output[idx["Date parution"]] = row[idx["Date parution"]]
		output[idx["Code du journal"]] = row[idx["Code du journal"]]
		output[idx["Code de la nature de l'événement"]] = row[idx["Code de la nature de l'événement"]]
		newSiret, ok := mapping[row[idx["Siret"]]]
		output[idx["Siret"]] = newSiret
		if ok {
			outputRow := strings.Join(output, ",") + "\n"

			_, err = outputFile.WriteString(outputRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func index(row []string) map[string]int {
	index := make(map[string]int)
	for i, v := range row {
		index[v] = i
	}
	return index
}
