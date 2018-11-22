package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/cnf/structhash"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

// BDF Information Banque de France
type BDF struct {
	Siren               string    `json:"siren" bson:"siren"`
	Annee               *int       `json:"annee_bdf" bson:"annee_bdf"`
	ArreteBilan         time.Time `json:"arrete_bilan_bdf" bson:"arrete_bilan_bdf"`
	RaisonSociale       string    `json:"raison_sociale" bson:"raison_sociale"`
	Secteur             string    `json:"secteur" bson:"secteur"`
	PoidsFrng           *float64  `json:"poids_frng" bson:"poids_frng"`
	TauxMarge           *float64  `json:"taux_marge" bson:"taux_marge"`
	DelaiFournisseur    *float64  `json:"delai_fournisseur" bson:"delai_fournisseur"`
	DetteFiscale        *float64  `json:"dette_fiscale" bson:"dette_fiscale"`
	FinancierCourtTerme *float64  `json:"financier_court_terme" bson:"financier_court_terme"`
	FraisFinancier      *float64  `json:"frais_financier" bson:"frais_financier"`
}

func parseBDF(path []string) chan *BDF {
	outputChannel := make(chan *BDF)

	go func() {
		for _, file := range path {
			var errorLines []int
			n := 0
			e := 0

			xlFile, err := xlsx.OpenFile(viper.GetString("APP_DATA") + file)
			if err != nil {
				log(critical, "importBDF", "Erreur à l'ouverture du fichier "+file+": "+err.Error())
			} else {
				log(info, "importBDF", "Ouverture du fichier "+file)

				for _, sheet := range xlFile.Sheets {
					for _, row := range sheet.Rows[1:] {
						n++
						var errors [11]error
						bdf := BDF{}
						bdf.Siren = strings.Replace(row.Cells[0].Value, " ", "", -1)
						bdf.Annee, errors[0] = parsePInt(row.Cells[1].Value)
						bdf.ArreteBilan, errors[1] = excelToTime(row.Cells[2].Value)
						bdf.RaisonSociale = row.Cells[3].Value
						bdf.Secteur = row.Cells[4].Value
						bdf.PoidsFrng, errors[5] = parsePFloat(row.Cells[5].Value)
						bdf.TauxMarge, errors[6] = parsePFloat(row.Cells[6].Value)
						bdf.DelaiFournisseur, errors[7] = parsePFloat(row.Cells[7].Value)
						bdf.DetteFiscale, errors[8] = parsePFloat(row.Cells[8].Value)
						bdf.FinancierCourtTerme, errors[9] = parsePFloat(row.Cells[9].Value)
						bdf.FraisFinancier, errors[10] = parsePFloat(row.Cells[10].Value)

						if allErrors(errors[:], nil) {
							outputChannel <- &bdf
						} else {
							e++
							errorLines = append(errorLines, n)
						}
					}
				}
			}
			log(debug, "importBDF", "Import du fichier Banque de France "+file+" terminé. "+fmt.Sprint(n)+" lignes traitée(s), "+fmt.Sprint(e)+" rejet(s)")
			if len(errorLines) > 0 {
				log(warning, "importBDF", "Erreurs de conversion constatées aux lignes suivantes: "+fmt.Sprintf("%v", errorLines))
			}
		}
		close(outputChannel)
	}()
	return outputChannel
}

func importBDF(batch *AdminBatch) error {
	log(info, "importBDF", "Import du batch "+batch.ID.Key+": Banque de France")

	for bdf := range parseBDF(batch.Files["bdf"]) {
		hash := fmt.Sprintf("%x", structhash.Md5(bdf, 1))

		value := ValueEntreprise{
			Value: Entreprise{
				Siren: bdf.Siren,
				Batch: map[string]Batch{
					batch.ID.Key: Batch{
						BDF: map[string]*BDF{
							hash: bdf,
						}}}}}
		db.ChanEntreprise <- &value
	}
	db.ChanEntreprise <- &ValueEntreprise{}
	log(info, "importBDF", "Fin de l'import du batch "+batch.ID.Key+": Banque de France")

	return nil
}
