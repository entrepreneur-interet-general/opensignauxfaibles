package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

// MapReduceJS Ensemble de fonctions JS pour mongodb
type MapReduceJS struct {
	Routine  string
	Scope    string
	Map      string
	Reduce   string
	Finalize string
}

func loadMR(typeJob string, target string) (*mgo.MapReduce, error) {
	mr := &mgo.MapReduce{}

	file, err := ioutil.ReadDir("js/" + typeJob + "/" + target)
	sort.Slice(file, func(i, j int) bool {
		return file[i].Name() < file[j].Name()
	})

	if err != nil {
		return nil, errors.New("Chemin introuvable")
	}

	mr.Map = ""
	mr.Reduce = ""
	mr.Finalize = ""

	for _, f := range file {
		if match, _ := regexp.MatchString("^map.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + typeJob + "/" + target + "/" + f.Name())
			if err != nil {
				return nil, errors.New("Lecture impossible: js/" + typeJob + "/" + target + "/" + f.Name())
			}
			mr.Map = mr.Map + string(fp)
		}
		if match, _ := regexp.MatchString("^reduce.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + typeJob + "/" + target + "/" + f.Name())
			if err != nil {
				return nil, errors.New("Lecture impossible: js/" + typeJob + "/" + target + "/" + f.Name())
			}
			mr.Reduce = mr.Reduce + string(fp)
		}
		if match, _ := regexp.MatchString("^finalize.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + typeJob + "/" + target + "/" + f.Name())
			if err != nil {
				return nil, errors.New("Lecture impossible: js/" + typeJob + "/" + target + "/" + f.Name())
			}
			mr.Finalize = mr.Finalize + string(fp)
		}
	}
	return mr, nil

}

func (mr *MapReduceJS) load(routine string, scope string) error {
	file, err := ioutil.ReadDir("js/" + routine + "/" + scope)
	sort.Slice(file, func(i, j int) bool {
		return file[i].Name() < file[j].Name()
	})

	if err != nil {
		return errors.New("Chemin introuvable")
	}

	mr.Routine = routine
	mr.Scope = scope
	mr.Map = ""
	mr.Reduce = ""
	mr.Finalize = ""

	for _, f := range file {
		if match, _ := regexp.MatchString("^map.*js$", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Map = mr.Map + string(fp)
		}
		if match, _ := regexp.MatchString("^reduce.*js$", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Reduce = mr.Reduce + string(fp)
		}
		if match, _ := regexp.MatchString("^finalize.*js$", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Finalize = mr.Finalize + string(fp)
		}
	}
	return nil
}

func reduceHandler(c *gin.Context) {
	algo := c.Params.ByName("algo")
	batchKey := c.Params.ByName("batch")
	siret := c.Params.ByName("siret")

	batch, _ := getBatch(batchKey)
	result, err := reduce(batch, algo, siret)

	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, result)
	}
}

func reduce(batch AdminBatch, algo string, siret string) (interface{}, error) {
	var queryEtablissement interface{}
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	dateDebut := batch.Params.DateDebut
	dateFin := batch.Params.DateFin
	dateFinEffectif := batch.Params.DateFinEffectif

	if siret == "" {
		db.DB.C("Features").RemoveAll(bson.M{"_id.batch": batch.ID.Key, "_id.algo": algo})
		queryEtablissement = bson.M{"value.index." + algo: true}
		queryEntreprise = nil
		output = bson.M{"merge": "Features"}
	} else {
		queryEtablissement = bson.M{"value.siret": siret}
		queryEntreprise = bson.M{"value.siren": siret[0:9]}
		output = nil
	}

	MREtablissement := MapReduceJS{}
	MREntreprise := MapReduceJS{}
	MRUnion := MapReduceJS{}
	errEt := MREtablissement.load(algo, "etablissement")
	errEn := MREntreprise.load(algo, "entreprise")
	errUn := MRUnion.load(algo, "union")

	if errEt != nil || errEn != nil || errUn != nil {
		return nil, fmt.Errorf("Problème d'accès aux fichiers MapReduce")
	}

	naf, errNAF := loadNAF()
	if errNAF != nil {
		return nil, fmt.Errorf("Problème d'accès aux fichiers naf")
	}

	scope := bson.M{
		"date_debut":             dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
	}

	jobEtablissement := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      bson.M{"replace": "MRWorkspace"},
		Scope:    scope,
	}

	_, err := db.DB.C("Etablissement").Find(queryEtablissement).MapReduce(jobEtablissement, nil)
	if err != nil {
		return nil, fmt.Errorf("Erreur du traitement MapReduce Etablissement: " + err.Error())
	}

	jobEntreprise := &mgo.MapReduce{
		Map:      string(MREntreprise.Map),
		Reduce:   string(MREntreprise.Reduce),
		Finalize: string(MREntreprise.Finalize),
		Out:      bson.M{"merge": "MRWorkspace"},
		Scope:    scope,
	}

	_, err = db.DB.C("Entreprise").Find(queryEntreprise).MapReduce(jobEntreprise, nil)
	if err != nil {
		return nil, fmt.Errorf("Erreur du traitement Entreprise: " + err.Error())
	}

	jobUnion := &mgo.MapReduce{
		Map:      string(MRUnion.Map),
		Reduce:   string(MRUnion.Reduce),
		Finalize: string(MRUnion.Finalize),
		Out:      output,
		Scope:    scope,
	}

	if output == nil {
		_, err = db.DB.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, &result)
	} else {
		_, err = db.DB.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, nil)
	}

	if err != nil {
		return result, fmt.Errorf("Erreur du traitement MapReduce Union")
	}

	return result, nil

}

func compactEtablissementHandler(c *gin.Context) {
	siret := c.Params.ByName("siret")

	err := compactEtablissement(siret)

	if err != nil {
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
	}
}

func compactEtablissement(siret string) error {
	batches, _ := getBatches()

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var etablissement []interface{}
	var completeTypes = make(map[string][]string)
	var batchesID []string

	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}

	// Si le parametre siret est absent, on traite l'ensemble de la collection
	if siret == "" {
		query = nil
		output = bson.M{"replace": "Etablissement"}
		etablissement = nil
	} else {
		query = bson.M{"value.siret": siret}
		output = nil
	}

	// Ressources JS
	MREtablissement := MapReduceJS{}
	errEt := MREtablissement.load("compact", "etablissement")

	if errEt != nil {
		return errEt
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      output,
		Scope: bson.M{"batches": batchesID,
			"types": []string{
				"altares",
				"apconso",
				"apdemande",
				"ccsf",
				"cotisation",
				"debit",
				"delai",
				"effectif",
				"sirene",
				"dpae",
			},
			"completeTypes": completeTypes,
		},
	}

	err := errors.New("")
	if output == nil {
		_, err = db.DB.C("Etablissement").Find(query).MapReduce(job, &etablissement)
	} else {
		_, err = db.DB.C("Etablissement").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		return err
	}
	return nil
}

func getFeatures(c *gin.Context) {
	var data []interface{}
	db.DB.C("Features").Find(nil).All(&data)
	c.JSON(200, data)
}

func compactEntrepriseHandler(c *gin.Context) {
	siren := c.Params.ByName("siren")
	err := compactEntreprise(siren)

	if err != nil {
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}
}
func compactEntreprise(siren string) error {
	batches, _ := getBatches()

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var etablissement []interface{}
	var completeTypes = make(map[string][]string)
	var batchesID []string

	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}

	if siren == "" {
		query = nil
		output = bson.M{"replace": "Entreprise"}
		etablissement = nil
	} else {
		query = bson.M{"value.siren": siren}
		output = nil
	}

	// Ressources JS
	MREntreprise := MapReduceJS{}
	errEn := MREntreprise.load("compact", "entreprise")

	if errEn != nil {
		return errEn
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(MREntreprise.Map),
		Reduce:   string(MREntreprise.Reduce),
		Finalize: string(MREntreprise.Finalize),
		Out:      output,
		Scope: bson.M{
			"batches": batchesID,
			"types": []string{
				"bdf",
				"diane",
			},
			"completeTypes": completeTypes,
		},
	}

	var err error

	if output == nil {
		_, err = db.DB.C("Entreprise").Find(query).MapReduce(job, &etablissement)
	} else {
		_, err = db.DB.C("Entreprise").Find(query).MapReduce(job, nil)
	}

	return err
}

func getNAF(c *gin.Context) {
	naf, err := loadNAF()
	if err != nil {
		c.JSON(500, err)
	}
	c.JSON(200, naf)
}

// NAF libellés et liens N5/N1
type NAF struct {
	N1    map[string]string `json:"n1" bson:"n1"`
	N5    map[string]string `json:"n5" bson:"n5"`
	N5to1 map[string]string `json:"n5to1" bson:"n5to1"`
}

func loadNAF() (NAF, error) {
	naf := NAF{}
	naf.N1 = make(map[string]string)
	naf.N5 = make(map[string]string)
	naf.N5to1 = make(map[string]string)

	NAF1 := viper.GetString("NAF_L1")
	NAF5 := viper.GetString("NAF_L5")
	NAF5to1 := viper.GetString("NAF_5TO1")

	NAF1File, NAF1err := os.Open(NAF1)
	if NAF1err != nil {
		fmt.Println(NAF1err)
		return NAF{}, NAF1err
	}

	NAF1reader := csv.NewReader(bufio.NewReader(NAF1File))
	NAF1reader.Comma = ';'
	NAF1reader.Read()
	for {
		row, err := NAF1reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		naf.N1[row[0]] = row[1]
	}

	NAF5to1File, NAF5to1err := os.Open(NAF5to1)
	if NAF5to1err != nil {
		fmt.Println(NAF5to1err)
		return NAF{}, NAF5to1err
	}

	NAF5to1reader := csv.NewReader(bufio.NewReader(NAF5to1File))
	NAF5to1reader.Comma = ';'
	NAF5to1reader.Read()
	for {
		row, err := NAF5to1reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}
		naf.N5to1[row[0]] = row[1]
	}

	NAF5File, NAF5err := os.Open(NAF5)
	if NAF5err != nil {
		fmt.Println(NAF5err)
		return NAF{}, NAF5err
	}

	NAF5reader := csv.NewReader(bufio.NewReader(NAF5File))
	NAF5reader.Comma = ';'
	NAF5reader.Read()
	for {
		row, err := NAF5reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
		}

		naf.N5[row[0]] = row[1]

	}
	return naf, nil
}
