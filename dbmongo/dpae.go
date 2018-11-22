package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/cnf/structhash"
	"github.com/spf13/viper"
)

// DPAE Déclaration préalabre à l'embauche
type DPAE struct {
	Siret    string    `json:"-" bson:"-"`
	Date     time.Time `json:"date" bson:"date"`
	CDI      float64   `json:"cdi" bson:"cdi"`
	CDDLong  float64   `json:"cdd_long" bson:"cdd_long"`
	CDDCourt float64   `json:"cdd_court" bson:"cdd_court"`
}

func parseDPAE(path string) chan *DPAE {
	outputChannel := make(chan *DPAE)

	file, err := os.Open(viper.GetString("APP_DATA") + path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.Read()

	go func() {
		for {
			row, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				// log.Fatal(error)
			}

			date, err := time.Parse("20060102", row[1]+row[2]+"01")

			dpae := DPAE{
				Siret: row[0],
				Date:  date,
			}
			dpae.CDI, _ = strconv.ParseFloat(row[3], 64)
			dpae.CDDLong, _ = strconv.ParseFloat(row[4], 64)
			dpae.CDDCourt, _ = strconv.ParseFloat(row[5], 64)

			if err == nil {
				outputChannel <- &dpae

			}
		}
		file.Close()
		close(outputChannel)
	}()
	return outputChannel
}

func importDPAE(batch *AdminBatch) error {
	for _, dpaeFile := range batch.Files["dpae"] {
		for dpae := range parseDPAE(dpaeFile) {
			hash := fmt.Sprintf("%x", structhash.Md5(dpae, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: dpae.Siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							DPAE: map[string]*DPAE{
								hash: dpae,
							}}}}}
			db.ChanEtablissement <- &value
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	return nil
}
