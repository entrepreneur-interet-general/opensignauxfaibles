package urssaf

import (
	"dbmongo/lib/engine"
	"testing"
  "time"
)


func Test_parseCompte(t *testing.T) {
	batch := engine.AdminBatch{
		Files: engine.BatchFiles{
			"admin_urssaf": []string{"testData/admin_urssaf.csv"},
		},
	}

	dataChannel, eventChannel := Parser(batch)
	engine.DiscardEvents(eventChannel)

  tuple := <-dataChannel
	data := tuple.(Compte)
  periode_init, _ := time.Parse("2006-01-02", "2014-01-01")


	if data.Siret == "12345678900001" &&
  data.NumeroCompte == "123456789123456789" &&
  data.Periode ==  periode_init {
		t.Log("Test parseCompte: lecture ok")
	} else {
		t.Error("Erreur parseCompte: données différentes du fichier de test")
	}

  for range(dataChannel){
  }

}
