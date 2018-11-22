package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

func readAndRandomComptes(fileName string, outputFileName string) (map[string]string, error) {
	// source
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	date := "1000101"
	mapping := make(map[string]string)
	// destination
	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	// ligne de titre
	row, err := reader.Read()
	outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return nil, err
	}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		siret := row[3]
		compte := row[0]
		newSiret := ""
		newCompte := ""
		for {
			newSiret = randStringBytesRmndr(len(siret))
			if _, ok := mapping[newSiret]; !ok && newSiret != siret {
				break
			} else {
				fmt.Println("bounce !")
			}
		}
		for {
			newCompte = randStringBytesRmndr(len(compte))
			if _, ok := mapping[newCompte]; !ok && newCompte != compte {
				break
			} else {
				fmt.Println("bounce !")
			}
		}
		mapping[compte] = newCompte
		mapping[siret] = newSiret

		row[0] = newCompte
		row[2] = newSiret[0:9]
		row[3] = newSiret
		row[4] = date
		row[5] = date

		outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
		_, err = outputFile.WriteString(outputRow)
		if err != nil {
			return nil, err
		}
	}
	return mapping, nil
}
