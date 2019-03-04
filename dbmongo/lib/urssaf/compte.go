package urssaf

import (
  "dbmongo/lib/engine"
  "dbmongo/lib/misc"
  "time"
)

// Compte tuple fichier ursaff
type Compte struct {
  key               string
  Siret             string    `json:"siret" bson:"siret"`
  NumeroCompte      string    `json:"numero_compte" bson:"numero_compte"`
  Periode           time.Time `json:"periode" bson:"periode"`
}

// Key _id de l'objet
func (compte Compte) Key() string {
  return compte.key
}

// Scope de l'objet
func (compte Compte) Scope() string {
  return "etablissemnt"
}

// Type de l'objet
func (compte Compte) Type() string {
  return "compte"
}

func parseCompte(batch engine.AdminBatch, mapping Comptes) (chan engine.Tuple, chan engine.Event) {
  outputChannel := make(chan engine.Tuple)
  eventChannel := make(chan engine.Event)

  //	event := engine.Event{
  //	Code:    "compteParser",
  //	Channel: eventChannel,
  //}

  go func(){
    periode_init, _ := time.Parse("2006-01-02", "2014-01-01")
    periodes := misc.GenereSeriePeriode(periode_init, time.Now()) //[]time.Time
    for c := range(mapping){
      for _, p := range(periodes) {
        compte := Compte{}
        compte.NumeroCompte = c
        compte.Periode = p
        compte.Siret, _ = mapping.GetSiret(c, p)
        outputChannel <- compte
      }
    }
    close(outputChannel)
    close(eventChannel)
  }()
  return outputChannel, eventChannel
}
