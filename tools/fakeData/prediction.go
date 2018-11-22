package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/spf13/viper"
)

// Prediction ligne de prédiction.
type Prediction struct {
	ID struct {
		Siret string `json:"siret" bson:"siret"`
		Batch string `json:"batch" bson:"batch"`
		Algo  string `json:"algo" bson:"algo"`
	} `json:"id" bson:"_id"`

	Prob          float64 `json:"prob" bson:"prob"`
	Diff          float64 `json:"diff" bson:"diff"`
	RaiSoc        string  `json:"raison_sociale" bson:"raison_sociale"`
	Departement   string  `json:"departement" bson:"departement"`
	Region        string  `json:"region" bson:"region"`
	EtatProCol    string  `json:"procol" bson:"procol"`
	DefaultUrssaf bool    `json:"default_urssaf" bson:"default_urssaf"`
	Connu         bool    `json:"connu" bson:"connu"`
	Niveau1       string  `json:"naf1" bson:"naf1"`
	Effectif      int     `json:"effectif" bson:"effectif"`
	CCSF          bool    `json:"ccsf" bson:"ccsf"`
}

func readSireneRaisonSociale(sirenePath string) map[string]string {
	// source
	file, err := os.Open(sirenePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ','

	reader.Read()

	raisoc := make(map[string]string)
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		siret := row[0] + row[1]
		raisoc[siret] = row[2]
	}

	return raisoc
}

func readAndRandomPrediction(fileName string, outputFileName string, mapping map[string]string) error {
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

	sireneFile := outputFileNamePrefixer(viper.GetString("prefixOutput"), viper.GetString("sirene"))
	raisoc := readSireneRaisonSociale(sireneFile)
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
		row[3] = raisoc[mapping[row[0]]]
		row[0] = mapping[row[0]]
		row[4] = "21"
		row[5] = "Bourgogne-Franche-Comté"
		outputRow := "\"" + strings.Join(row, "\";\"") + "\"\n"
		_, err = outputFile.WriteString(outputRow)
		if err != nil {
			return err
		}
	}

	return nil
}

func importPrediction(mapping map[string]string) {

	csvFile, _ := os.Open("./prediction_1809_utf.csv")

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.Read()
	var prediction []interface{}

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println("niktoo")
			log.Fatal(error)
		}

		proba, _ := strconv.ParseFloat(line[1], 64)
		diff, _ := strconv.ParseFloat(line[2], 64)
		effectif, _ := strconv.Atoi(line[11])
		var ccsf bool
		if line[12] == "" {
			ccsf = false
		} else {
			ccsf = true
		}

		var defaultUrssaf bool
		if line[7] == "TRUE" {
			defaultUrssaf = true
		} else {
			defaultUrssaf = false
		}

		var connu bool
		if line[8] == "0" {
			connu = false
		} else {
			connu = true
		}

		p := Prediction{
			ID: struct {
				Siret string `json:"siret" bson:"siret"`
				Batch string `json:"batch" bson:"batch"`
				Algo  string `json:"algo" bson:"algo"`
			}{
				Siret: line[0],
				Batch: "1802",
				Algo:  "algo1",
			},
			Prob:          proba,
			Diff:          diff,
			RaiSoc:        line[3],
			Departement:   line[4],
			Region:        line[5],
			EtatProCol:    line[6],
			DefaultUrssaf: defaultUrssaf,
			Connu:         connu,
			Niveau1:       line[10],
			Effectif:      effectif,
			CCSF:          ccsf,
		}

		prediction = append(prediction, p)
	}

	mongodb, err := mgo.Dial("")
	if err != nil {
		fmt.Println("pouet pouet")
		fmt.Println(err)
		return
	}
	db := mongodb.DB("signauxfaibles")

	db.C("Prediction").Insert(prediction...)
}
