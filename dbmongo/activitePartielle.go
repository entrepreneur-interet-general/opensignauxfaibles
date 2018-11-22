package main

import (
	"fmt"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/viper"

	"github.com/cnf/structhash"
	"github.com/tealeg/xlsx"
)

// APDemande Demande d'activité partielle
type APDemande struct {
	ID                 string    `json:"id_demande" bson:"id_demande"`
	Siret              string    `json:"-" bson:"-"`
	EffectifEntreprise *int      `json:"effectif_entreprise" bson:"effectif_entreprise"`
	Effectif           *int      `json:"effectif" bson:"effectif"`
	DateStatut         time.Time `json:"date_statut" bson:"date_statut"`
	TxPC               *float64  `json:"tx_pc" bson:"tx_pc"`
	TxPCUnedicDares    *float64  `json:"tx_pc_unedic_dares" bson:"tx_pc_unedic_dares"`
	TxPCEtatDares      *float64  `json:"tx_pc_etat_dares" bson:"tx_pc_etat_dares"`
	Periode            Periode   `json:"periode" bson:"periode"`
	HTA                *float64  `json:"hta" bson:"hta"`
	MTA                *float64  `json:"mta" bson:"mta"`
	EffectifAutorise   *int      `json:"effectif_autorise" bson:"effectif_autorise"`
	ProdHTAEffectif    *float64  `json:"prod_hta_effectif" bson:"prod_hta_effectif"`
	MotifRecoursSE     *int      `json:"motif_recours_se" bson:"motif_recours_se"`
	Perimetre          *int      `json:"perimetre" bson:"perimetre"`
	RecoursAnterieur   *int      `json:"recours_anterieur" bson:"recours_anterieur"`
	AvisCE             *int      `json:"avis_ce" bson:"avis_ce"`
	HeureConsommee     *float64  `json:"heure_consommee" bson:"heure_consommee"`
	MontantConsomme    *float64  `json:"montant_consommee" bson:"montant_consommee"`
	EffectifConsomme   *int      `json:"effectif_consomme" bson:"effectif_consomme"`
}

// APConso Consommation d'activité partielle
type APConso struct {
	ID             string    `json:"id_conso" bson:"id_conso"`
	Siret          string    `json:"-" bson:"-"`
	HeureConsommee *float64  `json:"heure_consomme" bson:"heure_consomme"`
	Montant        *float64  `json:"montant" bson:"montant"`
	Effectif       *int      `json:"effectif" bson:"effectif"`
	Periode        time.Time `json:"periode" bson:"periode"`
}

func parseAPDemande(path string) chan *APDemande {
	outputChannel := make(chan *APDemande)

	go func() {
		var errorLines []int
		n := 0
		e := 0

		xlFile, err := xlsx.OpenFile(path)
		if err != nil {
			log(critical, "importAPDemande", "Erreur à l'ouverture du fichier: "+path+", erreur: "+err.Error())
			close(outputChannel)
		} else {
			log(debug, "importAPDemande", "Ouverture du fichier, feuille 0, "+path)
			sheet := xlFile.Sheets[0]
			f := make(map[string]int)
			for idx, cell := range sheet.Rows[0].Cells {
				f[cell.Value] = idx
			}

			fields := []string{
				"ID_DA",
				"ETAB_SIRET",
				"EFF_ENT",
				"EFF_ETAB",
				"DATE_STATUT",
				"DATE_DEB",
				"DATE_FIN",
				"HTA",
				"EFF_AUTO",
				"MOTIF_RECOURS_SE",
				"S_HEURE_CONSOM_TOT",
				"S_EFF_CONSOM_TOT",
			}
			minLength := 0
			for _, field := range fields {
				if i, err := f[field]; err {
					minLength = max(minLength, i)
				} else {
					log(critical, "importAPDemande", "Import du fichier "+path+". "+field+" non trouvé. Abandon.")
					close(outputChannel)
					return
				}
			}
			// FIXME: établir une meilleure méthode pour tester la validité de la ligne
			for _, row := range sheet.Rows[3:] {
				var errors [10]error
				if len(row.Cells) >= minLength {
					n++
					apdemande := APDemande{}
					apdemande.ID = row.Cells[f["ID_DA"]].Value
					apdemande.Siret = row.Cells[f["ETAB_SIRET"]].Value
					apdemande.EffectifEntreprise, errors[0] = parsePInt(row.Cells[f["EFF_ENT"]].Value)
					apdemande.Effectif, errors[1] = parsePInt(row.Cells[f["EFF_ETAB"]].Value)
					apdemande.DateStatut, errors[2] = excelToTime(row.Cells[f["DATE_STATUT"]].Value)
					apdemande.Periode = Periode{}
					apdemande.Periode.Start, errors[3] = excelToTime(row.Cells[f["DATE_DEB"]].Value)
					apdemande.Periode.End, errors[4] = excelToTime(row.Cells[f["DATE_FIN"]].Value)
					apdemande.HTA, errors[5] = parsePFloat(row.Cells[f["HTA"]].Value)
					apdemande.EffectifAutorise, errors[6] = parsePInt(row.Cells[f["EFF_AUTO"]].Value)
					apdemande.MotifRecoursSE, errors[7] = parsePInt(row.Cells[f["MOTIF_RECOURS_SE"]].Value)
					apdemande.HeureConsommee, errors[8] = parsePFloat(row.Cells[f["S_HEURE_CONSOM_TOT"]].Value)
					apdemande.EffectifConsomme, errors[9] = parsePInt(row.Cells[f["S_EFF_CONSOM_TOT"]].Value)

					if allErrors(errors[:], nil) && apdemande.Siret != "" {
						outputChannel <- &apdemande
					} else {
						e++
						errorLines = append(errorLines, n)
					}
				}

			}
		}
		log(debug, "importAPDemande", "Import du fichier "+path+" terminé. "+fmt.Sprint(n)+" lignes traitée(s), "+fmt.Sprint(e)+" rejet(s)")
		if len(errorLines) > 0 {
			log(warning, "importAPDemande", "Erreurs de conversion constatées aux lignes suivantes: "+fmt.Sprintf("%v", errorLines))
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importAPDemande(batch *AdminBatch) error {
	log(info, "importAPDemande", "Import du batch "+batch.ID.Key+": APDemande")
	for _, file := range batch.Files["apdemande"] {
		for apdemande := range parseAPDemande(viper.GetString("APP_DATA") + file) {
			hash := fmt.Sprintf("%x", structhash.Md5(apdemande, 1))
			value := ValueEtablissement{
				Value: Etablissement{
					Siret: apdemande.Siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							Compact: map[string]bool{
								"status": false,
							},
							APDemande: map[string]*APDemande{
								hash: apdemande,
							}}}}}
			if value.Value.Batch == nil {
				spew.Dump(value)
			}
			db.ChanEtablissement <- &value
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	log(info, "importAPDemande", "Import APDemande du batch "+batch.ID.Key+" terminé")
	return nil
}

func parseAPConso(path string) chan *APConso {
	outputChannel := make(chan *APConso)
	xlFile, err := xlsx.OpenFile(viper.GetString("APP_DATA") + path)
	if err != nil {
		log(critical, "importApconso", "Erreur à l'ouverture du fichier "+path+": "+err.Error())
	} else {
		log(debug, "importAPConso", "Ouverture du fichier "+path)
		go func() {
			var errorLines []int
			n := 0
			e := 0
			for idx, sheet := range xlFile.Sheets {
				log(debug, "importAPConso", "Import de la feuille "+fmt.Sprint(idx))
				fields := sheet.Rows[0]
				idxID := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "ID_DA" })
				idxSiret := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "ETAB_SIRET" })
				idxPeriode := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "MOIS" })
				idxHeureConsommee := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "HEURES" })
				idxMontants := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "MONTANTS" })
				idxEffectifs := sliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "EFFECTIFS" })

				for _, row := range sheet.Rows[1:] {
					if len(row.Cells) > 0 {
						n++
						var errors [4]error
						apconso := APConso{}
						apconso.ID = row.Cells[idxID].Value
						apconso.Siret = row.Cells[idxSiret].Value
						apconso.Periode, errors[0] = excelToTime(row.Cells[idxPeriode].Value)
						apconso.HeureConsommee, errors[1] = parsePFloat(row.Cells[idxHeureConsommee].Value)
						apconso.Montant, errors[2] = parsePFloat(row.Cells[idxMontants].Value)
						apconso.Effectif, errors[3] = parsePInt(row.Cells[idxEffectifs].Value)

						if allErrors(errors[:], nil) && apconso.Siret != "" {
							outputChannel <- &apconso
						} else {
							e++
							errorLines = append(errorLines, n)
						}
					}
				}
			}
			log(debug, "importAPConso", "Import du fichier délais "+path+" terminé. "+fmt.Sprint(n)+" lignes traitée(s), "+fmt.Sprint(e)+" rejet(s)")
			if len(errorLines) > 0 {
				log(warning, "importAPConso", "Erreurs de conversion constatées aux lignes suivantes: "+fmt.Sprintf("%v", errorLines))
			}
			close(outputChannel)
		}()
	}
	return outputChannel
}

func importAPConso(batch *AdminBatch) error {
	log(info, "importAPConso", "Import du batch "+batch.ID.Key+": APDemande")

	for _, file := range batch.Files["apconso"] {
		for apconso := range parseAPConso(file) {
			hash := fmt.Sprintf("%x", structhash.Md5(apconso, 1))
			value := ValueEtablissement{
				Value: Etablissement{
					Siret: apconso.Siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							APConso: map[string]*APConso{
								hash: apconso,
							}}}}}
			if value.Value.Batch == nil {
				spew.Dump(value)
			}
			db.ChanEtablissement <- &value
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	log(info, "importAPConso", "Fin de l'import du batch "+batch.ID.Key+": APConso")
	return nil
}
