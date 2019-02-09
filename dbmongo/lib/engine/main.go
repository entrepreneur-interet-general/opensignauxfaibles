package engine

import (
	"strconv"
	"time"

	"github.com/chrnin/gournal"
	"github.com/globalsign/mgo/bson"
)

// Db connecteur exportable
var Db DB

// Priority test
type Priority string

// Code test
type Code string

// Event test
type Event struct {
	ID       bson.ObjectId `json:"id" bson:"_id"`
	Date     time.Time     `json:"date" bson:"date"`
	Comment  interface{}   `json:"event" bson:"event"`
	Priority Priority      `json:"priority" bson:"priority"`
	Code     Code          `json:"code" bson:"code"`
	Channel  chan Event    `json:"-"`
}

// GetBSON retourne l'objet Event sous une forme sérialisable
func (event Event) GetBSON() (interface{}, error) {
	var tmp struct {
		ID       bson.ObjectId `json:"id" bson:"_id"`
		Date     time.Time     `json:"date" bson:"date"`
		Comment  interface{}   `json:"event" bson:"event"`
		Priority Priority      `json:"priority" bson:"priority"`
		Code     Code          `json:"code" bson:"code"`
	}
	tmp.ID = event.ID
	tmp.Date = event.Date
	tmp.Comment = event.Comment
	tmp.Priority = event.Priority
	tmp.Code = event.Code
	return tmp, nil
}

// Debug test
var Debug = Priority("debug")

// Info test
var Info = Priority("info")

// Warning test
var Warning = Priority("warning")

// Critical test
var Critical = Priority("critical")

var unknownCode = Code("unknown")

// DiscardEvents supprime les évènements
func DiscardEvents(events chan Event) {
	go func() {
		for range events {
			// fmt.Println(e)
		}
	}()
}

// Debug produit un évènement de niveau Debug
func (event Event) Debug(comment interface{}) {
	event.ID = bson.NewObjectId()
	event.Date = time.Now()
	event.Comment = comment
	event.Priority = Debug
	if event.Code == "" {
		event.Code = unknownCode
	}
	event.Channel <- event
}

// Info produit un évènement de niveau Info
func (event Event) Info(comment interface{}) {
	event.ID = bson.NewObjectId()
	event.Date = time.Now()
	event.Comment = comment
	event.Priority = Info
	if event.Code == "" {
		event.Code = unknownCode
	}
	event.Channel <- event
}

// Warning produit un évènement de niveau Warning
func (event Event) Warning(comment interface{}) {
	event.ID = bson.NewObjectId()
	event.Date = time.Now()
	event.Comment = comment
	event.Priority = Warning
	if event.Code == "" {
		event.Code = unknownCode
	}
	event.Channel <- event
}

// Critical produit un évènement de niveau Critical
func (event Event) Critical(comment interface{}) {
	event.ID = bson.NewObjectId()
	event.Date = time.Now()
	event.Comment = comment
	event.Priority = Critical
	if event.Code == "" {
		event.Code = unknownCode
	}
	event.Channel <- event
}

func reportAbstract(tracker gournal.Tracker) interface{} {
	return tracker.Context["path"] + ": intégration terminée, " +
		strconv.Itoa(tracker.Count) + " éléments traités " +
		strconv.Itoa(tracker.CountErrors()) + " rejets."
}

func reportErrors(tracker gournal.Tracker) interface{} {
	return bson.M{
		"report":      tracker.Context["path"] + ": ligne " + strconv.Itoa(tracker.Count) + " ignorée",
		"errorReport": tracker.CurrentErrors(),
	}
}

func reportInvalidData(tracker gournal.Tracker) interface{} {
	if errs, ok := tracker.Errors[tracker.Count]; ok {
		return tracker.Context["path"] + ": cycle " + strconv.Itoa(tracker.Count) + " ignoré: " + errs[len(errs)-1].Error()
	}
	return nil
}

func reportFatalError(tracker gournal.Tracker) interface{} {
	if errs, ok := tracker.Errors[tracker.Count]; ok {
		return "Erreur fatale, abandon: " + errs[len(errs)-1].Error()
	}
	return nil
}

// TrackerReports contient les fonctions de reporting du moteur
var TrackerReports = map[string]gournal.ReportFunction{
	"abstract":    reportAbstract,
	"errors":      reportErrors,
	"invalidLine": reportInvalidData,
	"fatalError":  reportFatalError,
}
