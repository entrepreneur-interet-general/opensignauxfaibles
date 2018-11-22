package main

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
)

func readAndRandomSirene(fileName string, outputFileName string, mapping map[string]string) error {
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
	reader.Comma = ','

	title, err := reader.Read()
	outputRow := strings.Join(title, ",") + "\n"
	outputLength := len(title)

	_, err = outputFile.WriteString(outputRow)
	if err != nil {
		return err
	}

	var wordBase = make(map[string]struct{})

	for i := 0; i < 100000; i++ {
		row, err := reader.Read()
		if err == io.EOF {
			return nil
		} else if err != nil {
			panic(err)
		}
		words := strings.Split(row[2], " ")
		for _, word := range words {
			wordBase[word] = struct{}{}
		}
	}

	var wordList []string
	for key := range wordBase {
		wordList = append(wordList, key)
	}

	wordLength := float64(len(wordList))
	length := int(rand.Float64()*4) + 1

	for _, v := range mapping {
		if len(v) == 14 {
			output := make([]string, outputLength)
			// siren
			output[0] = v[0:9]
			// nic
			output[1] = v[9:14]
			// nic siege
			output[65] = v[9:14]

			output[16] = "21"
			output[19] = "RUE"
			output[20] = "21000"
			output[23] = "Bourgogne Franche-Comté"
			output[24] = "Côte d'or"
			output[28] = "Dijon"
			output[42] = "2120Z"
			output[71] = "SA"
			output[81] = "2018-01-01"
			output[51] = "2018-01-01"
			output[100] = "5.051709"
			output[101] = "47.315471"
			output[2] = "21 BOULEVARD VOLTAIRE"
			output[4] = "21000 DIJON"

			var sentence []string
			for i := 0; i < length; i++ {
				sentence = append(sentence, wordList[int(rand.Float64()*wordLength)])
			}
			// raison sociale
			output[2] = strings.Join(sentence, " ")

			outputRow := "\"" + strings.Join(output, "\",\"") + "\"\n"
			_, err = outputFile.WriteString(outputRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
