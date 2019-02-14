package apartconso

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/testtools"
	"testing"
)

func Test_parseAPConso(t *testing.T) {
	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"apconso": []string{"testData/apconso.excelsheet"},
		},
	}

	outputChannel, eventChannel := Parser(batch)
	go testtools.DiscardEventChannel(eventChannel)
	tuples := testtools.OutputChannelToSlice(outputChannel)
	data := tuples[0].(APConso)

	if len(tuples) == 1 {
		t.Log("Le fichier contient une seule ligne: Ok")
	} else {
		t.Fatal("Le fichier contient une seule ligne: Fail")
	}

	if *data.Effectif == 30 &&
		*data.HeureConsommee == 100 &&
		data.ID == "123456789" &&
		data.Siret == "12345678901234" &&
		*data.Montant == 200 {
		t.Log("Les bonne données sont lues: ok")
	} else {
		t.Error("Les bonnes données sont lues: fail")
	}

	outputChannel, eventChannel = Parser(batch)
	go testtools.DiscardOutputChannel(outputChannel)
	events := testtools.EventChannelToSlice(eventChannel)

	if len(events) == 2 {
		t.Log("L'intégration génère 2 évènements: Ok")
	} else {
		t.Fatal("L'intégration génère 2 évènements: Fail")
	}

	if events[0].Code == "info" && events[1].Code == "info" {
		t.Log("Les codes sont fixés correctement: Ok")
	} else {
		t.Fatal("Les codes sont fixés correctement: Fail")
	}
	// outputChannel, eventChannel = Parser(batch)
	// testtools.DiscardEventChannel(eventChannel)
	// data :=
	// if tuple != nil {
	// 	t.Error("Erreur parseAPConso: fichier absent: le channel devrait être vide")
	// } else {
	// 	t.Log("Test parseAPConso: fichier absent ok")
	// }
}
