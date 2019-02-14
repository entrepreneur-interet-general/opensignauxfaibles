package altares

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/testtools"
	"testing"

	"github.com/spf13/viper"
)

func Test_ParserEventsNoFile(t *testing.T) {
	viper.SetDefault("APP_DATA", ".")

	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"altares": []string{"nofile"},
		},
	}
	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardOutputChannel(outputChannel)
	events := testtools.EventChannelToSlice(eventChannel)

	if len(events) == 1 {
		t.Log("No file should return 1 event: Ok")
	} else {
		t.Fatal("No file should return 1 event: Fail")
	}

	if events[0].Code == "altaresParser" {
		t.Log("Code should be altaresParser: Ok")
	} else {
		t.Fatal("Code should be altaresParser: Fail")
	}

	if events[0].Priority == "critical" {
		t.Log("Priority should be critical: Ok")
	} else {
		t.Fatal("Priority should be critical: Fail")
	}
}

func Test_ParserDatasNoFile(t *testing.T) {
	viper.SetDefault("APP_DATA", ".")

	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"altares": []string{"nofile"},
		},
	}
	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardEventChannel(eventChannel)
	datas := testtools.OutputChannelToSlice(outputChannel)

	if len(datas) == 0 {
		t.Log("No file should return no data: Ok")
	} else {
		t.Fatal("No file should return no data: Fail")
	}
}

func Test_ParserEventsBadHeaders(t *testing.T) {
	viper.SetDefault("APP_DATA", ".")

	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"altares": []string{"/testData/testWithBadHeaders.csvfile"},
		},
	}
	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardOutputChannel(outputChannel)
	events := testtools.EventChannelToSlice(eventChannel)
	if len(events) == 2 {
		t.Log("Bad headers file should return 2 events: Ok")
	} else {
		t.Fatal("Bad headers file should return 2 events: Fail")
	}

	if events[0].Priority == "info" && events[1].Priority == "critical" {
		t.Log("Priority are info then critical: Ok")
	} else {
		t.Fatal("Priority are info then critical: Fail")
	}
}

func Test_ParserDatasBadHeaders(t *testing.T) {
	viper.SetDefault("APP_DATA", ".")

	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"altares": []string{"/testData/testWithBadHeaders.csvfile"},
		},
	}
	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardEventChannel(eventChannel)
	datas := testtools.OutputChannelToSlice(outputChannel)

	if len(datas) == 0 {
		t.Log("Bad headers file file should return no data: Ok")
	} else {
		t.Fatal("Bad headers should return no data: Fail")
	}
}

func Test_ParserEventsWithGoodFile(t *testing.T) {
	viper.SetDefault("APP_DATA", ".")

	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"altares": []string{"/testData/testWithoutErrors.csvfile"},
		},
	}
	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardOutputChannel(outputChannel)
	events := testtools.EventChannelToSlice(eventChannel)

	if len(events) == 2 {
		t.Log("Good file should return 2 events: Ok")
	} else {
		t.Fatal("Good file should return 2 events: Fail")
	}

	fail := false
	for _, e := range events {
		if e.Code != "altaresParser" {
			fail = true
		}
	}

	if !fail {
		t.Log("Code should be altaresParser: Ok")
	} else {
		t.Fatal("Code should be altaresParser: Fail")
	}

	if events[0].Priority == "info" && events[1].Priority == "info" {
		t.Log("Priority should be info: Ok")
	} else {
		t.Fatal("Priority should be info: Fail")
	}
}

func Test_ParserEventsWithBadFile(t *testing.T) {
	viper.SetDefault("APP_DATA", ".")

	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"altares": []string{"/testData/testWithErrors.csvfile"},
		},
	}
	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardOutputChannel(outputChannel)
	events := testtools.EventChannelToSlice(eventChannel)

	if len(events) == 3 {
		t.Log("Bad file should return 3 events: Ok")
	} else {
		t.Fatal("Bad file should return 3 events: Fail")
	}

	fail := false
	for _, e := range events {
		if e.Code != "altaresParser" {
			fail = true
		}
	}

	if !fail {
		t.Log("Code should be altaresParser: Ok")
	} else {
		t.Fatal("Code should be altaresParser: Fail")
	}

	if events[0].Priority == "info" && events[2].Priority == "info" && events[1].Priority == "debug" {
		t.Log("Priority should be info: Ok")
	} else {
		t.Fatal("Priority should be info: Fail")
	}
}

func Test_AltaresType(t *testing.T) {
	a := Altares{
		Siret: "test",
	}
	if a.Key() == "test" && a.Type() == "altares" && a.Scope() == "etablissement" {
		t.Log("Tuple interface is correct: Ok")
	} else {
		t.Fatal("Tuple interface is correct: Fail")
	}
}
