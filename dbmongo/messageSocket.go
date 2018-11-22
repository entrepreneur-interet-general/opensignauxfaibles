package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

type journalEvent struct {
	ID       bson.ObjectId   `json:"id" bson:"_id"`
	Date     time.Time       `json:"date" bson:"date"`
	Comment  string          `json:"event" bson:"event"`
	Priority journalPriority `json:"priority" bson:"priority"`
	Code     journalCode     `json:"code" bson:"code"`
}

type socketMessage struct {
	JournalEvent journalEvent  `json:"journalEvent" bson:"journalEvent"`
	Batches      []AdminBatch  `json:"batches,omitempty" bson:"batches,omitempty"`
	Types        Types         `json:"types,omitempty" bson:"types,omitempty"`
	Features     []string      `json:"features,omitempty" bson:"features,omitempty"`
	Files        []fileSummary `json:"files,omitempty" bson:"files,omitempty"`
}

type journalCode string
type journalPriority string

type messageChannel chan socketMessage

var messageClientChannels = []messageChannel{}
var mainMessageChannel = messageDispatch()
var addClientChannel = make(chan messageChannel)

var debug = journalPriority("debug")
var info = journalPriority("info")
var warning = journalPriority("warning")
var critical = journalPriority("critical")

// messageSocketAddClientChannel surveille l'ajout de nouveaux clients pour propager le channel d'events
func messageSocketAddClient() {
	for clientChannel := range addClientChannel {
		messageClientChannels = append(messageClientChannels, clientChannel)
	}
}

// journal dispatch un event vers les clients et l'enregistre dans la bdd
func messageDispatch() chan socketMessage {
	channel := make(messageChannel)
	go func() {
		for event := range channel {
			db.DBStatus.C("Journal").Insert(event.JournalEvent)
			for _, clientChannel := range messageClientChannels {
				clientChannel <- event
			}
		}
	}()
	return channel
}

func log(priority journalPriority, code journalCode, comment string) journalEvent {
	e := journalEvent{
		ID:       bson.NewObjectId(),
		Date:     time.Now(),
		Comment:  comment,
		Priority: priority,
		Code:     code,
	}
	mainMessageChannel <- socketMessage{
		JournalEvent: e,
	}
	return e
}

func getLogs() ([]journalEvent, error) {
	var logs []journalEvent
	err := db.DB.C("Journal").Find(nil).Sort("-date").Limit(250).All(&logs)
	return logs, err
}

func getLogsHandler(c *gin.Context) {
	logs, err := getLogs()
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, logs)
	}
}
