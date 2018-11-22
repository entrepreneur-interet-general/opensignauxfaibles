package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/cnf/structhash"
)

// Debit Débit – fichier Urssaf
type Debit struct {
	NumeroCompte                 string    `json:"numero_compte" bson:"numero_compte"`
	NumeroEcartNegatif           string    `json:"numero_ecart_negatif" bson:"numero_ecart_negatif"`
	DateTraitement               time.Time `json:"date_traitement" bson:"date_traitement"`
	PartOuvriere                 float64   `json:"part_ouvriere" bson:"part_ouvriere"`
	PartPatronale                float64   `json:"part_patronale" bson:"part_patronale"`
	NumeroHistoriqueEcartNegatif int       `json:"numero_historique" bson:"numero_historique"`
	EtatCompte                   int       `json:"etat_compte" bson:"etat_compte"`
	CodeProcedureCollective      string    `json:"code_procedure_collective" bson:"code_procedure_collective"`
	Periode                      Periode   `json:"periode" bson:"periode"`
	CodeOperationEcartNegatif    string    `json:"code_operation_ecart_negatif" bson:"code_operation_ecart_negatif"`
	CodeMotifEcartNegatif        string    `json:"code_motif_ecart_negatif" bson:"code_motif_ecart_negatif"`
	DebitSuivant                 string    `json:"debit_suivant,omitempty" bson:"debit_suivant,omitempty"`
	MontantMajorations			 float64   `json:"montant_majorations" bson:"montant_majorations"`
}

func parseDebit(paths []string) chan *Debit {
	outputChannel := make(chan *Debit)

	go func() {
		for _, path := range paths {
			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			fields, err := reader.Read()

			dateTraitementIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Dt_trt_ecn" })
			partOuvriereIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PO" })
			partPatronaleIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Mt_PP" })
			numeroHistoriqueEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Hist_Ecn" })
			periodeIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Periode" })
			etatCompteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Etat_cpte" })

			numeroCompteIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "num_cpte" })
			numeroEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Num_Ecn" })
			codeProcedureCollectiveIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_pro_col" })
			codeOperationEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Cd_op_ecn" })
			codeMotifEcartNegatifIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Motif_ecn" })
			montantMajorationsIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] ==  "Montant majorations de retard en centimes"})
			for {
				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					// log.Fatal(error)
				}

				debit := Debit{
					NumeroCompte:              row[numeroCompteIndex],
					NumeroEcartNegatif:        row[numeroEcartNegatifIndex],
					CodeProcedureCollective:   row[codeProcedureCollectiveIndex],
					CodeOperationEcartNegatif: row[codeOperationEcartNegatifIndex],
					CodeMotifEcartNegatif:     row[codeMotifEcartNegatifIndex],
				}
				debit.DateTraitement, err = UrssafToDate(row[dateTraitementIndex])
				debit.PartOuvriere, err = strconv.ParseFloat(row[partOuvriereIndex], 64)
				debit.PartOuvriere = debit.PartOuvriere / 100
				debit.PartPatronale, err = strconv.ParseFloat(row[partPatronaleIndex], 64)
				debit.PartPatronale = debit.PartPatronale / 100
				debit.NumeroHistoriqueEcartNegatif, err = strconv.Atoi(row[numeroHistoriqueEcartNegatifIndex])
				debit.EtatCompte, err = strconv.Atoi(row[etatCompteIndex])
				debit.Periode, err = UrssafToPeriod(row[periodeIndex])
				debit.MontantMajorations, err =  strconv.ParseFloat(row[montantMajorationsIndex], 64)
				debit.MontantMajorations = debit.MontantMajorations / 100

				outputChannel <- &debit
			}
			file.Close()
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importDebit(batch *AdminBatch) error {
	mapping, _ := getCompteSiretMapping(batch)

	for debit := range parseDebit(batch.Files["debit"]) {
		if siret, ok := mapping[debit.NumeroCompte]; ok {
			hash := fmt.Sprintf("%x", structhash.Md5(debit, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							Debit: map[string]*Debit{
								hash: debit,
							}}}}}
			db.ChanEtablissement <- &value
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	return nil
}
