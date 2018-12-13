package main

import (
  "bufio"
  "encoding/csv"
  "fmt"
  "io"
  "os"
  "time"
  "strings"
  "strconv"
  "regexp"

  "github.com/cnf/structhash"
  "github.com/spf13/viper"
)

// Procol Procédures collectives, extraction URSSAF
type Procol struct {
  DateEffet     time.Time `json:"date_effet" bson:"date_effet"`
  ActionProcol  string    `json:"action_procol" bson:"action_procol"`
  StadeProcol   string    `json:"stade_procol" bson:"stade_procol"`
  Siret         string    `json:"-" bson:"-"`
}

func parseProcol(path string) chan *Procol {
  outputChannel := make(chan *Procol)

  file, err := os.Open(path)
  if err != nil {
    log(critical, "importProcol", "Erreur à l'ouverture du fichier: "+path+", erreur: "+err.Error())
    close(outputChannel)
  }

  reader := csv.NewReader(bufio.NewReader(file))
  reader.Comma = ';'
  reader.LazyQuotes = true
  fields, err := reader.Read() //headline

  dateEffetIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "dt_effet" })
  actionStadeIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "lib_actx_stdx" })
  siretIndex := sliceIndex(len(fields), func(i int) bool { return fields[i] == "Siret" })

  go func() {
    var errorLines []int
    n := 0
    e := 0

    for {

      row, err := reader.Read()

      if err == io.EOF {
        break
      } else if err != nil {
        log(critical, "importProcol", "Erreur de lecture pendant l'import du fichier "+path+". Abandon.")
        close(outputChannel)
      }
      if _, err := strconv.Atoi(row[siretIndex]); err == nil && len(row[siretIndex]) == 14 {
        // Uniquement les lignes avec des sirets valides sont comptées
        n++


        var errors [2]error
        procol := Procol{}

        date_formatee := row[dateEffetIndex]
        date_formatee = date_formatee[:3]+strings.ToLower(date_formatee[4:5])+date_formatee[6:]
        procol.DateEffet, errors[1] = time.Parse("02Jan2006", row[dateEffetIndex])
        procol.Siret = row[siretIndex]
        splitted := strings.Split(strings.ToLower(row[actionStadeIndex]), "_")

        for i, v := range splitted {
          r, _ := regexp.Compile("liquidation|redressement|sauvegarde")
          if match := r.MatchString(v); match {
            procol.ActionProcol = v
            procol.StadeProcol = strings.Join(append(splitted[:i], splitted[i+1:]...), "_")
            break
          }
        }

        if allErrors(errors[:], nil) && procol.Siret != "" {
          outputChannel <- &procol
        } else {
          e++
          errorLines = append(errorLines, n)
        }
      }
    }

    log(debug, "importProcol", "Import du fichier "+path+" terminé. "+fmt.Sprint(n)+" lignes traitée(s), "+fmt.Sprint(e)+" rejet(s)")
    if len(errorLines) > 0 {
       log(warning, "importProcol", "Erreurs de conversion constatées aux lignes suivantes: "+fmt.Sprintf("%v", errorLines))
    }
    file.Close()
    close(outputChannel)
  }()
  return outputChannel
}

func importProcol(batch *AdminBatch) error {
  for _, fileName := range batch.Files["procol"] {
    for procol := range parseProcol(viper.GetString("APP_DATA") + fileName) {
      hash := fmt.Sprintf("%x", structhash.Md5(procol, 1))

      value := ValueEtablissement{
        Value: Etablissement{
          Siret: procol.Siret,
          Batch: map[string]Batch{
            batch.ID.Key: Batch{
              Procol: map[string]*Procol{
                hash: procol,
              },
            },
          },
        },
      }
      db.ChanEtablissement <- &value
    }
  }
  db.ChanEtablissement <- &ValueEtablissement{}
  return nil
}
