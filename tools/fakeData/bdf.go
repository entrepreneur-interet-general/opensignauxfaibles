package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"
)

func readAndRandomBDF(fileName string, outputFileName string, mapping map[string]string) error {
	file, err := xlsx.OpenFile(fileName)
	if err != nil {
		return err
	}
	// destination
	outputFile := xlsx.NewFile()
	outputFile.AddSheet("Sheet1")
	sirens := make(map[string]string)
	for k, v := range mapping {
		sirens[k[0:9]] = v[0:9]
	}
	for _, sheet := range file.Sheets {
		for _, row := range sheet.Rows[1:] {
			row.Cells[3].Value = ""
			siren := strings.Replace(row.Cells[0].Value, " ", "", -1)
			row.Cells[0].Value = sirens[siren][0:3] + " " + sirens[siren][3:6] + " " + sirens[siren][6:9]
			c5, _ := strconv.ParseFloat(row.Cells[5].Value, 64)
			c6, _ := strconv.ParseFloat(row.Cells[6].Value, 64)
			c7, _ := strconv.ParseFloat(row.Cells[7].Value, 64)
			c8, _ := strconv.ParseFloat(row.Cells[8].Value, 64)
			c9, _ := strconv.ParseFloat(row.Cells[9].Value, 64)
			c10, _ := strconv.ParseFloat(row.Cells[10].Value, 64)
			if row.Cells[5].Value != "" {
				row.Cells[5].Value = fmt.Sprintf("%f", c5*(rand.Float64()+0.50))
			}
			if row.Cells[6].Value != "" {
				row.Cells[6].Value = fmt.Sprintf("%f", c6*(rand.Float64()+0.50))
			}
			if row.Cells[7].Value != "" {
				row.Cells[7].Value = fmt.Sprintf("%f", c7*(rand.Float64()+0.50))
			}
			if row.Cells[8].Value != "" {
				row.Cells[8].Value = fmt.Sprintf("%f", c8*(rand.Float64()+0.50))
			}
			if row.Cells[9].Value != "" {
				row.Cells[9].Value = fmt.Sprintf("%f", c9*(rand.Float64()+0.50))
			}
			if row.Cells[10].Value != "" {
				row.Cells[10].Value = fmt.Sprintf("%f", c10*(rand.Float64()+0.50))
			}
		}
		outputFile.Sheet["Sheet1"] = sheet
	}
	err = file.Save(outputFileName)

	return nil
}
