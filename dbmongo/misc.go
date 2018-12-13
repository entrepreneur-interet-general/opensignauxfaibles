package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/spf13/viper"
)

func sliceIndex(limit int, predicate func(i int) bool) int {
	for i := 0; i < limit; i++ {
		if predicate(i) {
			return i
		}
	}
	return -1
}

func sliceIndexArray(target []string, data []string) []int {
	l := len(target)
	var idx []int

	for a := 1; a <= l; a++ {
		idx = append(idx, -1)
	}

	for idxTarget, t := range target {
		for idxData, d := range data {
			if d == t {
				idx[idxTarget] = idxData
			}
		}
	}

	return idx
}

type stringSlice []string

func (arr stringSlice) contains(str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}

func excelToTime(excel string) (time.Time, error) {
	excelInt, err := strconv.ParseInt(excel, 10, 64)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}
	return time.Unix((excelInt-25569)*3600*24, 0), nil
}

func getCompteSiretMapping(batch *AdminBatch) (map[string]string, error) {
	compteSiretMapping := make(map[string]string)
  compteSiretLast := make(map[string]int)

  path := batch.Files["admin_urssaf"]
	basePath := viper.GetString("APP_DATA")

	for _, p := range path {
		file, err := os.Open(basePath + p)
		if err != nil {
			return map[string]string{}, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
		}

		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'

		// discard header row
		reader.Read()

		siretIndex := 3
		compteIndex := 0
    fermetureIndex := 5

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log(critical, "importCompteSiret", "Erreur à la lecture du fichier "+file.Name())
				return map[string]string{}, err
			}

      _, err1 := strconv.Atoi(row[siretIndex]);
      fermeture, err2 := strconv.Atoi(row[fermetureIndex]);
      if err2 != nil {
        if row[fermetureIndex] == "" {
          fermeture = 1
        } else {
          log(critical, "importCompteSiret", "Erreur (2) à la lecture du fichier "+file.Name()+err2.Error())
          return map[string]string{}, err2
        }
      }
      derniereFermetureLue, ok  := compteSiretLast[row[compteIndex]];
      if  err1 == nil &&
          len(row[siretIndex]) == 14 &&
          (!ok ||
             (derniereFermetureLue != 0 && derniereFermetureLue < fermeture) ||
             fermeture == 0) {

				compteSiretMapping[row[compteIndex]] = row[siretIndex]
        compteSiretLast[row[compteIndex]] = fermeture
			}
		}
		file.Close()
	}
	return compteSiretMapping, nil
}

func getSirensFromMapping(batch *AdminBatch) (map[string]bool, error) {
	SirensFromMapping := make(map[string]bool)
	path := batch.Files["admin_urssaf"]
	basePath := viper.GetString("APP_DATA")

	for _, p := range path {
		file, err := os.Open(basePath + p)
		if err != nil {
			return map[string]bool{}, errors.New("Erreur à l'ouverture du fichier, " + err.Error())
		}

		reader := csv.NewReader(bufio.NewReader(file))
		reader.Comma = ';'

		// discard header row
		reader.Read()

		siretIndex := 3

		for {
			row, err := reader.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log(critical, "importCompteSiret", "Erreur à la lecture du fichier "+file.Name())
				return map[string]bool{}, err
			}
			if _, err := strconv.Atoi(row[siretIndex]); err == nil && len(row[siretIndex]) == 14 {
				SirensFromMapping[row[siretIndex][0:9]] = true
			}
		}
		file.Close()
	}
	return SirensFromMapping, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func genereSeriePeriode(debut time.Time, fin time.Time) []time.Time {
	var serie []time.Time
	for fin.After(debut) {
		debut = debut.AddDate(0, 1, 0)
		serie = append(serie, debut)
	}
	return serie
}

func genereSeriePeriodeAnnuelle(debut time.Time, fin time.Time) []int {
	var serie []int
	for debut.Year() <= fin.Year() {
		serie = append(serie, debut.Year())
		debut = debut.AddDate(1, 0, 0)
	}
	return serie
}

func allErrors(slice []error, item interface{}) bool {
	for _, i := range slice {
		if i != item {
			return false
		}
	}
	return true
}

func parsePInt(s string) (*int, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.Atoi(s)
	return &i, err
}

func parsePFloat(s string) (*float64, error) {
	if s == "" {
		return nil, nil
	}
	i, err := strconv.ParseFloat(s, 64)
	return &i, err
}

// UrssafToPeriod convertit le format de période urssaf en type Period
func UrssafToPeriod(urssaf string) (Periode, error) {
	// format en 4 ou 6 caractère YYQM ou YYYYQM
	// si YY < 50 alors YYYY = 20YY sinon YYYY = 19YY
	// si QM == 62 alors période annuelle sur YYYY
	// si M == 0 alors période trimestrielle sur le trimestre Q de YYYY
	// si 0 < M < 4 alors mois M du trimestre Q

	period := Periode{}

	if len(urssaf) == 4 {
		if urssaf[0:2] < "50" {
			urssaf = "20" + urssaf
		} else {
			urssaf = "19" + urssaf
		}
	}

	if len(urssaf) != 6 {
		return period, errors.New("Valeur non autorisée")
	}

	year, err := strconv.Atoi(urssaf[0:4])
	if err != nil {
		return period, err
	}

	if urssaf[4:6] == "62" {
		period.Start = time.Date(year, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
		period.End = time.Date(year+1, time.Month(1), 1, 0, 0, 0, 0, time.UTC)
	} else {
		quarter, err := strconv.Atoi(urssaf[4:5])
		if err != nil {
			return period, err
		}
		monthOfQuarter, err := strconv.Atoi(urssaf[5:6])
		if err != nil {
			return period, err
		}
		if monthOfQuarter == 0 {
			period.Start = time.Date(year, time.Month((quarter-1)*3+1), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+4), 1, 0, 0, 0, 0, time.UTC)
		} else {
			period.Start = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter), 1, 0, 0, 0, 0, time.UTC)
			period.End = time.Date(year, time.Month((quarter-1)*3+monthOfQuarter+1), 1, 0, 0, 0, 0, time.UTC)
		}
	}
	return period, nil
}

// UrssafToDate Convertit le format de date urssaf en type Date
func UrssafToDate(urssaf string) (time.Time, error) {
	// Date au format YYYMMJJ
	// YYY = YYYY - 1900

	intUrsaff, err := strconv.Atoi(urssaf)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}
	strDate := strconv.Itoa(intUrsaff + 19000000)
	date, err := time.Parse("20060102", strDate)
	if err != nil {
		return time.Time{}, errors.New("Valeur non autorisée")
	}

	return date, nil
}
