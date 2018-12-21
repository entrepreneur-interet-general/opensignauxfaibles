package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/cnf/structhash"
	"github.com/spf13/viper"
)

// Sirene informations sur les entreprises
type Sirene struct {
	Siren              string    `json,omitempty:"siren" bson,omitempty:"siren"`
	Nic                string    `json,omitempty:"nic" bson,omitempty:"nic"`
	NicSiege           string    `json,omitempty:"nic_siege" bson,omitempty:"nic_siege"`
	RaisonSociale      string    `json,omitempty:"raison_sociale" bson,omitempty:"raison_sociale"`
	NumVoie            string    `json,omitempty:"numero_voie" bson,omitempty:"numero_voie"`
	IndRep             string    `json,omitempty:"indrep" bson,omitempty:"indrep"`
	TypeVoie           string    `json,omitempty:"type_voie" bson,omitempty:"type_voie"`
	CodePostal         string    `json,omitempty:"code_postal" bson,omitempty:"code_postal"`
	Cedex              string    `json,omitempty:"cedex" bson,omitempty:"cedex"`
	Region             string    `json,omitempty:"region" bson,omitempty:"region"`
	Departement        string    `json,omitempty:"departement" bson,omitempty:"departement"`
	Commune            string    `json,omitempty:"commune" bson,omitempty:"commune"`
	APE                string    `json,omitempty:"ape" bson,omitempty:"ape"`
	NatureActivite     string    `json,omitempty:"nature_activite" bson,omitempty:"nature_activite"`
	ActiviteSaisoniere string    `json,omitempty:"activite_saisoniere" bson,omitempty:"activite_sai"`
	ModaliteActivite   string    `json,omitempty:"modalite_activite" bson,omitempty:"modalite_activite"`
	Productif          string    `json,omitempty:"productif" bson,omitempty:"productif"`
	NatureJuridique    string    `json,omitempty:"nature_juridique" bson,omitempty:"nature_juridique"`
	Categorie          string    `json,omitempty:"categorie" bson,omitempty:"categorie"`
	Creation           time.Time `json,omitempty:"date_creation" bson,omitempty:"date_creation"`
	IndiceMonoactivite *int      `json,omitempty:"indice_monoactivite" bson,omitempty:"indice_monoactivite"`
	TrancheCA          *int      `json,omitempty:"tranche_ca" bson,omitempty:"tranche_ca"`
	Sigle              string    `json,omitempty:"sigle" bson,omitempty:"sigle"`
	Longitude          *float64  `json,omitempty:"longitude" bson:"longitude"`
	Lattitude          *float64  `json,omitempty:"lattitude" bson:"lattitude"`
	Adresse            [7]string `json:"adresse" bson:"adresse"`
}

func parseSirene(paths []string, mapping map[string]bool) chan *Sirene {
	outputChannel := make(chan *Sirene)
	//ignoreSirene := stringSlice(viper.GetStringSlice("SIRENE_IGNORE"))

	go func() {
		// var errorLines []int
		n := 0
		// e := 0
		for _, path := range paths {
			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
        log(critical, "importSirene", "Erreur à l'ouverture du fichier: "+path+", erreur: "+err.Error())
        close(outputChannel)
      }

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ','

			for {
				row, err := reader.Read()
				if err == io.EOF {
					break
				} else if err != nil {
					log(critical, "importSirene", "Erreur lors de la lecture du fichier "+path+". Abandon.")
					close(outputChannel)
				}


//				if _, ok := mapping[row[0]]; ok {
          n++

          sirene := Sirene{}
					sirene.Siren = row[0]


					sirene.Nic = row[1]
					sirene.NicSiege = row[65]
					sirene.RaisonSociale = row[2]
					sirene.NumVoie = row[16]
					sirene.IndRep = row[17]
					sirene.TypeVoie = row[19]
					sirene.CodePostal = row[20]
					sirene.Cedex = row[21]
					sirene.Region = row[23]
					sirene.Departement = row[24]
					sirene.Commune = row[28]
					sirene.APE = row[42]
					sirene.NatureActivite = row[52]
					sirene.ActiviteSaisoniere = row[55]
					sirene.ModaliteActivite = row[56]
					sirene.Productif = row[57]
					sirene.NatureJuridique = row[71]
					sirene.Categorie = row[82]

          if i, err := time.Parse("20060102", row[50]); err == nil {
            sirene.Creation = i
          }
          if i, err := strconv.Atoi(row[85]); err == nil {
            sirene.IndiceMonoactivite = &i
          }
          if i, err := strconv.Atoi(row[89]); err == nil {
            sirene.TrancheCA = &i
          }

          sirene.Sigle = row[61]

          if i , err := strconv.ParseFloat(row[100], 64); err == nil {
            sirene.Longitude = &i
          }
          if i, err := strconv.ParseFloat(row[101], 64); err == nil {
            sirene.Lattitude = &i
          }
          sirene.Adresse = [7]string{row[2], row[3], row[4], row[5], row[6], row[7], row[8]}

          //     for  i:=0; i < len(errors);i++ {
          //       if (errors[i] != nil){
          //         log(debug, "importSirene", "Erreur ? "+fmt.Sprint(i)+" - "+errors[i].Error())
          //       }
          //     }

          //			if allErrors(errors[:], nil) {
          //				outputChannel <- &sirene
          //			} else {
          //				e++
          //				errorLines = append(errorLines, n)
          //			}
          outputChannel <- &sirene
          // }

      }
      file.Close()
      log(debug, "importSirene", "Import du fichier "+path+" terminé. "+fmt.Sprint(n)+" lignes traitée(s)")//, "+fmt.Sprint(e)+" rejet(s)")
      // if len(errorLines) > 0 {
      //   log(warning, "importSirene", "Erreurs de conversion constatées aux lignes suivantes: "+fmt.Sprintf("%v", errorLines))
      // }
    }
    close(outputChannel)
  }()

  return outputChannel
}

// hash := fmt.Sprintf("%x", structhash.Md5(sirene, 1))

func importSirene(batch *AdminBatch) error {

  mapping, _ := getSirensFromMapping(batch)

  for sirene := range parseSirene(batch.Files["sirene"], mapping) {
    hash := fmt.Sprintf("%x", structhash.Md5(sirene, 1))

    value := ValueEtablissement{
      Value: Etablissement{
        Siret: sirene.Siren + sirene.Nic,
        Batch: map[string]Batch{
          batch.ID.Key: Batch{
            Sirene: map[string]*Sirene{
              hash: sirene,
            },
          },
        },
      },
    }
    db.ChanEtablissement <- &value
  }

  db.ChanEtablissement <- &ValueEtablissement{}
  return nil
}
