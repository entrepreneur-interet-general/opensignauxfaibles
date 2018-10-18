package main

import (
	"bufio"
	"encoding/csv"
//	"fmt"
	"io"
	"os"

	"github.com/spf13/viper"
)

// NAF libellés et liens N5/N1
type NAF struct {
	N1    map[string]string `json:"n1" bson:"n1"`
	N2    map[string]string `json:"n2" bson:"n2"`
	N3    map[string]string `json:"n3" bson:"n3"`
	N4    map[string]string `json:"n4" bson:"n4"`
	N5    map[string]string `json:"n5" bson:"n5"`
	N5to1 map[string]string `json:"n5to1" bson:"n5to1"`
}

func loadNAF() (NAF, error) {
	naf := NAF{}
	naf.N1 = make(map[string]string)
	naf.N2 = make(map[string]string)
	naf.N3 = make(map[string]string)
	naf.N4 = make(map[string]string)
	naf.N5 = make(map[string]string)
	naf.N5to1 = make(map[string]string)

	NAF1 := viper.GetString("NAF_L1")
	NAF2 := viper.GetString("NAF_L2")
	NAF3 := viper.GetString("NAF_L3")
	NAF4 := viper.GetString("NAF_L4")
	NAF5 := viper.GetString("NAF_L5")
	NAF5to1 := viper.GetString("NAF_5TO1")

	//NAF1
	NAF1File, NAF1err := os.Open(NAF1)

	if NAF1err != nil {
		return NAF{}, NAF1err
	}

	NAF1reader := csv.NewReader(bufio.NewReader(NAF1File))
	NAF1reader.Comma = ';'
	NAF1reader.Read()
	for {
		row, error := NAF1reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			// log.Fatal(error)
		}
		naf.N1[row[0]] = row[1]
	}

	//NAF2
	NAF2File, NAF2err := os.Open(NAF2)

	if NAF2err != nil {
		return NAF{}, NAF2err
	}

	NAF2reader := csv.NewReader(bufio.NewReader(NAF2File))
	NAF2reader.Comma = ';'
	NAF2reader.Read()
	for {
		row, error := NAF2reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
		}
		naf.N2[row[0]] = row[1]
	}

	//NAF3

	NAF3File, NAF3err := os.Open(NAF3)

	if NAF3err != nil {
		return NAF{}, NAF3err
	}

	NAF3reader := csv.NewReader(bufio.NewReader(NAF3File))
	NAF3reader.Comma = ';'
	NAF3reader.Read()
	for {
		row, error := NAF3reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
		}
		naf.N3[row[0]] = row[1]
	}
	//NAF4

	NAF4File, NAF4err := os.Open(NAF4)

	if NAF4err != nil {
		return NAF{}, NAF4err
	}

	NAF4reader := csv.NewReader(bufio.NewReader(NAF4File))
	NAF4reader.Comma = ';'
	NAF4reader.Read()
	for {
		row, error := NAF4reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
		}
		naf.N4[row[0]] = row[1]
	}

	//NAF5

	NAF5File, NAF5err := os.Open(NAF5)
	if NAF5err != nil {
		return NAF{}, NAF5err
	}

	NAF5reader := csv.NewReader(bufio.NewReader(NAF5File))
	NAF5reader.Comma = ';'
	NAF5reader.Read()
	for {
		row, error := NAF5reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			// log.Fatal(error)
		}

		naf.N5[row[0]] = row[1]
	}

	//NAF1TO5
	NAF5to1File, NAF5to1err := os.Open(NAF5to1)
	if NAF5to1err != nil {
		return NAF{}, NAF5to1err
	}

	NAF5to1reader := csv.NewReader(bufio.NewReader(NAF5to1File))
	NAF5to1reader.Comma = ';'
	NAF5to1reader.Read()
	for {
		row, error := NAF5to1reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			// log.Fatal(error)
		}
		naf.N5to1[row[0]] = row[1]
	}

	return naf, nil
}
