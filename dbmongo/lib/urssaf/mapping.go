package urssaf

import (
  "bufio"
  "dbmongo/lib/engine"
  "encoding/csv"
  "errors"
  "io"
  "os"
  //"strconv"
  "time"
  "sort"
  "github.com/spf13/viper"
)

type SiretDate struct{
  Siret string
  Date time.Time
}
type Comptes map[string][]SiretDate


func (c Comptes) GetSiret(compte string, date time.Time) (string, error) {
  date_keys := make([]SiretDate, len(c))
  i := 0
  for _, k := range c[compte] {
    date_keys[i] = k
    i++
  }

  found := false
  i = 0
  siret := ""
  for !found && i < len(date_keys) {
    if date.Before(date_keys[i].Date){
      found = true
      siret = c[compte][i].Siret
    }

    i++
  }
  if siret == "" {
    return siret, errors.New("Pas de siret associé au compte " + compte + " à cette période")
  }
  return siret, nil
}

func getCompteSiretMapping(batch *engine.AdminBatch) (Comptes, error) {

  compteSiretMapping := make(map[string][]SiretDate)

  path := batch.Files["admin_urssaf"]
  basePath := viper.GetString("APP_DATA")

  for _, p := range path {
    file, err := os.Open(basePath + p)
    if err != nil {
      return map[string][]SiretDate{}, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
    }

    reader := csv.NewReader(bufio.NewReader(file))
    reader.Comma = ';'

    // discard header row
    reader.Read()

    compteIndex := 0
    siretIndex := 3
    fermetureIndex := 5

    for {
      row, err := reader.Read()
      if err == io.EOF {
        break
      } else if err != nil {
        return map[string][]SiretDate{}, err
      }

      maxTime := "9990101"

      if row[fermetureIndex] == "" {row[fermetureIndex] = "0"} // date de fermeture manquante
      if (row[fermetureIndex] == "0") { row[fermetureIndex] = maxTime } // compte non fermé

      fermeture, err := urssafToDate(row[fermetureIndex])
      if  err != nil {
        return map[string][]SiretDate{}, err // fermeture n'a pas pu être lue ou convertie en date
      }

      compte := row[compteIndex]
      siret := row[siretIndex]
      if len(siret) == 14 {
        //siret valide
        compteSiretMapping[compte] = append(compteSiretMapping[compte], SiretDate{siret, fermeture})
        sort.Slice(compteSiretMapping[compte],
        func(i, j int) bool {return(
          compteSiretMapping[compte][i].Date.Before(compteSiretMapping[compte][j].Date))})
        }
      }
      file.Close()

    }
    return compteSiretMapping, nil
  }
