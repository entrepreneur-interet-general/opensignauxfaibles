package apartconso

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"time"

	"github.com/chrnin/gournal"
	"github.com/spf13/viper"
	"github.com/tealeg/xlsx"
)

// APConso Consommation d'activité partielle
type APConso struct {
	key            string
	scope          string
	datatype       string
	ID             string    `json:"id_conso" bson:"id_conso"`
	Siret          string    `json:"-" bson:"-"`
	HeureConsommee *float64  `json:"heure_consomme" bson:"heure_consomme"`
	Montant        *float64  `json:"montant" bson:"montant"`
	Effectif       *int      `json:"effectif" bson:"effectif"`
	Periode        time.Time `json:"periode" bson:"periode"`
}

// Key id de l'objet
func (apconso APConso) Key() string {
	return apconso.key
}

// Type de données
func (apconso APConso) Type() string {
	return apconso.datatype
}

// Scope de l'objet
func (apconso APConso) Scope() string {
	return apconso.scope
}

// Parser produit des lignes de consommation d'activité partielle
func Parser(batch engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)
	event := engine.Event{
		Code:    "parserApconso",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["apconso"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			xlFile, err := xlsx.OpenFile(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture du fichier, abandon:" + err.Error())
				continue
			}

			event.Info(path + ": ouverture")

			for _, sheet := range xlFile.Sheets {
				fields := sheet.Rows[0]
				idxID := misc.SliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "ID_DA" })
				idxSiret := misc.SliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "ETAB_SIRET" })
				idxPeriode := misc.SliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "MOIS" })
				idxHeureConsommee := misc.SliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "HEURES" })
				idxMontants := misc.SliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "MONTANTS" })
				idxEffectifs := misc.SliceIndex(35, func(i int) bool { return fields.Cells[i].Value == "EFFECTIFS" })
				if misc.SliceMin(idxID, idxSiret, idxPeriode, idxHeureConsommee, idxMontants, idxEffectifs) == -1 {
					event.Critical(path + ": entête non conforme, fichier ignoré")
					continue
				}

				for _, row := range sheet.Rows[1:] {
					if len(row.Cells) > 0 {
						apconso := APConso{}

						apconso.ID = row.Cells[idxID].Value
						apconso.Siret = row.Cells[idxSiret].Value
						apconso.key = apconso.Siret
						apconso.datatype = "apconso"
						apconso.scope = "etablissement"
						apconso.Periode, err = misc.ExcelToTime(row.Cells[idxPeriode].Value)
						tracker.Error(err)
						apconso.HeureConsommee, err = misc.ParsePFloat(row.Cells[idxHeureConsommee].Value)
						tracker.Error(err)
						apconso.Montant, err = misc.ParsePFloat(row.Cells[idxMontants].Value)
						tracker.Error(err)
						apconso.Effectif, err = misc.ParsePInt(row.Cells[idxEffectifs].Value)
						tracker.Error(err)

						if !tracker.ErrorInCycle() && apconso.Siret != "" {
							outputChannel <- apconso
						} else {
							event.Debug(tracker.Report("errors"))
						}
						tracker.Next()
					}
				}
			}
			event.Info(tracker.Report("abstract"))
		}
		close(outputChannel)
		close(eventChannel)
	}()

	return outputChannel, eventChannel
}
