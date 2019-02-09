package altares

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/misc"
	"encoding/csv"
	"io"
	"os"
	"time"

	"github.com/chrnin/gournal"
	"github.com/spf13/viper"
)

// Altares Extrait du récapitulatif altarès
type Altares struct {
	DateEffet     time.Time  `json:"date_effet" bson:"date_effet"`
	DateParution  *time.Time `json:"date_parution" bson:"date_parution"`
	CodeJournal   string     `json:"code_journal" bson:"code_journal"`
	CodeEvenement string     `json:"code_evenement" bson:"code_evenement"`
	Siret         string     `json:"-" bson:"-"`
}

// Key id de l'objet
func (altares Altares) Key() string {
	return altares.Siret
}

// Type de données
func (altares Altares) Type() string {
	return "altares"
}

// Scope de l'objet
func (altares Altares) Scope() string {
	return "etablissement"
}

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

// Parser du type Altares
func Parser(batch engine.AdminBatch) (chan engine.Tuple, chan engine.Event) {
	outputChannel := make(chan engine.Tuple)
	eventChannel := make(chan engine.Event)

	event := engine.Event{
		Code:    "altaresParser",
		Channel: eventChannel,
	}

	go func() {
		for _, path := range batch.Files["altares"] {
			tracker := gournal.NewTracker(
				map[string]string{"path": path},
				engine.TrackerReports)

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				event.Critical(path + ": erreur à l'ouverture: " + err.Error())
				continue
			}
			event.Info(path + ": ouverture")
			reader := csv.NewReader(file)
			reader.Comma = ','
			reader.LazyQuotes = true
			// ligne de titre
			fields, err := reader.Read()

			dateEffetIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Date d'effet" })
			dateParutionIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Date parution" })
			codeJournalIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Code du journal" })
			codeEvenementIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Code de la nature de l'événement" })
			siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Siret" })

			if misc.SliceMin(dateEffetIndex, dateParutionIndex, codeJournalIndex, codeEvenementIndex, siretIndex) == -1 {
				event.Critical(path + ": entête non conforme, fichier ignoré")
				file.Close()
				continue
			}

			for {
				row, err := reader.Read()
				if err == io.EOF {
					file.Close()
					break
				}

				dateEffet, err := time.Parse("2006-01-02", row[dateEffetIndex])
				tracker.Error(err)

				// si erreur, pointeur nil.
				var dateParution *time.Time
				if d, err := time.Parse("2006-01-02", row[dateParutionIndex]); err != nil {
					dateParution = &d
				}

				altares := Altares{
					Siret:         row[siretIndex],
					DateEffet:     dateEffet,
					DateParution:  dateParution,
					CodeJournal:   row[codeJournalIndex],
					CodeEvenement: row[codeEvenementIndex],
				}
				if !tracker.ErrorInCycle() {
					outputChannel <- altares
				} else {
					event.Debug(tracker.Report("errors"))
				}
				tracker.Next()
			}
			event.Info(tracker.Report("abstract"))
		}
		close(eventChannel)
		close(outputChannel)
	}()
	return outputChannel, eventChannel
}
