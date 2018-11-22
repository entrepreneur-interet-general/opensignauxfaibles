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

// Effectif Urssaf
type Effectif struct {
	Siret        string    `json:"-" bson:"-"`
	NumeroCompte string    `json:"numero_compte" bson:"numero_compte"`
	Periode      time.Time `json:"periode" bson:"periode"`
	Effectif     int       `json:"effectif" bson:"effectif"`
}

// ParseEffectifPeriod Transforme un tableau de périodes telles qu'écrites dans l'entête du tableau d'effectif urssaf en date de début
func parseEffectifPeriod(effectifPeriods []string) ([]time.Time, error) {
	periods := []time.Time{}
	for _, period := range effectifPeriods {
		urssaf := period[3:9]
		date, _ := UrssafToPeriod(urssaf)
		periods = append(periods, date.Start)
	}

	return periods, nil
}

func parseEffectif(paths []string) chan map[string]*Effectif {
	outputChannel := make(chan map[string]*Effectif)

	go func() {
		for _, path := range paths {

			file, err := os.Open(viper.GetString("APP_DATA") + path)

			if err != nil {
				fmt.Println("Error", err)

			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			fields, err := reader.Read()

			if err != nil {
				fmt.Println("Error", err)
			}

			siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "SIRET" })
			compteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "compte" })

			boundaryIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "ape_ins" })
			periods, err := parseEffectifPeriod(fields[0:boundaryIndex])

			if err != nil {
				// log.Panic("Aborting: could not read a period:", err)
			}

			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					// log.Fatal(error)
				}

				i := 0
				effectif := make(map[string]*Effectif)

				for i < boundaryIndex {
					e, _ := strconv.Atoi(row[i])
					if e > 0 {
						eff := Effectif{
							Siret:        row[siretIndex],
							NumeroCompte: row[compteIndex],
							Periode:      periods[i],
							Effectif:     e}
						hash := fmt.Sprintf("%x", structhash.Md5(eff, 1))
						effectif[hash] = &eff
					}
					i++
				}

				if len(row[siretIndex]) == 14 {
					outputChannel <- effectif
				}
			}
			file.Close()
		}
		close(outputChannel)
	}()
	return outputChannel
}

func importEffectif(batch *AdminBatch) error {
	for effectif := range parseEffectif(batch.Files["effectif"]) {

		randomKey := ""
		for randomKey = range effectif {
			break
		}

		if randomKey != "" {
			siret := effectif[randomKey].Siret

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							Effectif: effectif,
						}}}}
			db.ChanEtablissement <- &value
		}
	}

	db.ChanEtablissement <- &ValueEtablissement{}
	return nil
}
