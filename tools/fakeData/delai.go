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

// "Num�ro de compte externe";"Le n� de structure est l'identifiant d'un dossier contentieux";"Date de cr�ation";"Date d'�ch�ance-2";"Dur�e d�lai";"D�nomination premi�re ligne";"Indic 6M";"ann�e (  Date de cr�ation  )";"Montant global de l'�ch�ancier";"Num�ro de structure";"Code externe du stade";"Code externe de l'action"
// "267000001600134585";2014099226;2014-08-25;2014-10-27;63;"PIZZA SAONE LA ROYALE";"INF";2014;7529;2014099226;"APPROB";"SUR PO"

func readAndRandomDelais(fileName string, outputFileName string, mapping map[string]string) error {
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
		var output [12]string

		var c float64

		if c, ok := coef[row[0]]; !ok {
			c = rand.Float64() * rand.Float64()
			coef[row[1]] = c
		}

		if k, ok := mapping[row[0]]; ok {
			output[0] = k
			output[2] = "2017-01-01"
			output[3] = "2017-06-01"
			output[4] = "181"
			output[6] = row[6]
			output[7] = "2017"

			output[10] = row[10]
			output[11] = row[11]

			montant, err := strconv.ParseFloat(row[8], 64)
			output[8] = strconv.Itoa(int(montant * c))
			row[0] = mapping[row[0]]

			outputRow := "\"" + strings.Join(output[:], "\";\"") + "\"\n"
			_, err = outputFile.WriteString(outputRow)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
