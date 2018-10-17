package main

import (
	"errors"
	"io/ioutil"
	"regexp"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// MapReduceJS Ensemble de fonctions JS pour mongodb
type MapReduceJS struct {
	Routine  string
	Scope    string
	Map      string
	Reduce   string
	Finalize string
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
		if match, _ := regexp.MatchString("^map.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Map = mr.Map + string(fp)
		}
		if match, _ := regexp.MatchString("^reduce.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Reduce = mr.Reduce + string(fp)
		}
		if match, _ := regexp.MatchString("^finalize.*js", f.Name()); match {
			fp, err := ioutil.ReadFile("js/" + routine + "/" + scope + "/" + f.Name())
			if err != nil {
				return errors.New("Lecture impossible: js/" + routine + "/" + scope + "/" + f.Name())
			}
			mr.Finalize = mr.Finalize + string(fp)
		}
	}
	return nil
}

func dataPrediction(c *gin.Context) {
	var prediction []Prediction
	var etablissement []ValueEtablissement
	var siret []string

	db, _ := c.Keys["DB"].(*mgo.Database)
	db.C("prediction").Find(nil).Sort("-prob").Limit(50).All(&prediction)
	for _, r := range prediction {
		siret = append(siret, r.Siret)
	}

	query := bson.M{"_id": bson.M{"$in": siret}}
	db.C("Etablissement").Find(query).All(&etablissement)
	c.JSON(200, bson.M{"prediction": prediction, "etablissement": etablissement})
}

func reduce(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)

	dateDebut, _ := time.Parse("2006-01-02", "2014-01-01")
	dateFin, _ := time.Parse("2006-01-02", "2018-10-01")
	dateFinEffectif, _ := time.Parse("2006-01-02", "2018-06-01")

	// Détermination scope traitement
	algo := c.Params.ByName("algo")
	batch := c.Params.ByName("batch")
	siret := c.Params.ByName("siret")

	db.C("Features").RemoveAll(bson.M{"_id.batch": batch, "_id.algo": algo})

	var queryEtablissement interface{}
	var queryEntreprise interface{}
	var output interface{}
	var result interface{}

	if siret == "" {
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
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}

	scope := bson.M{"date_debut": dateDebut,
		"date_fin":               dateFin,
		"date_fin_effectif":      dateFinEffectif,
		"serie_periode":          genereSeriePeriode(dateDebut, dateFin),
		"serie_periode_annuelle": genereSeriePeriodeAnnuelle(dateDebut, dateFin),
		"actual_batch":           batch,
	}

	jobEtablissement := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      bson.M{"replace": "MRWorkspace"},
		Scope:    scope,
	}

	_, err := db.C("Etablissement").Find(queryEtablissement).MapReduce(jobEtablissement, nil)
	if err != nil {
		c.JSON(500, err)
		return
	}

	jobEntreprise := &mgo.MapReduce{
		Map:      string(MREntreprise.Map),
		Reduce:   string(MREntreprise.Reduce),
		Finalize: string(MREntreprise.Finalize),
		Out:      bson.M{"merge": "MRWorkspace"},
		Scope:    scope,
	}

	_, err = db.C("Entreprise").Find(queryEntreprise).MapReduce(jobEntreprise, nil)
	if err != nil {
		c.JSON(500, err)
		return
	}

	jobUnion := &mgo.MapReduce{
		Map:      string(MRUnion.Map),
		Reduce:   string(MRUnion.Reduce),
		Finalize: string(MRUnion.Finalize),
		Out:      output,
		Scope:    scope,
	}

	if output == nil {
		_, err = db.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, &result)
	} else {
		_, err = db.C("MRWorkspace").Find(queryEntreprise).MapReduce(jobUnion, nil)
	}

	if err != nil {
		c.JSON(500, err)
		return
	}

	err = db.C("MRWorkspace").DropCollection()
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(200, result)

}

func compactEtablissement(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batches := getBatchesID(db)

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var etablissement []interface{}

	// Si le parametre siret est absent, on traite l'ensemble de la collection
	siret := c.Params.ByName("siret")
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
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      output,
		Scope: bson.M{"batches": batches,
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
			"deleteOld": []string{"effectif", "apdemande", "apconso", "altares"},
		},
	}

	err := errors.New("")
	if output == nil {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, &etablissement)
	} else {
		_, err = db.C("Etablissement").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		c.JSON(500, "Echec du traitement MR, message serveur: "+err.Error())
	} else {
		c.JSON(200, etablissement)
	}

}

func compactEntreprise(c *gin.Context) {
	db, _ := c.Keys["DB"].(*mgo.Database)
	batches := getBatchesID(db)

	// Détermination scope traitement
	var query interface{}
	var output interface{}
	var etablissement []interface{}

	// Si le parametre siren est absent, on traite l'ensemble de la collection
	siren := c.Params.ByName("siren")
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
		c.JSON(500, "Problème d'accès aux fichiers MapReduce")
		return
	}

	// Traitement MR
	job := &mgo.MapReduce{
		Map:      string(MREntreprise.Map),
		Reduce:   string(MREntreprise.Reduce),
		Finalize: string(MREntreprise.Finalize),
		Out:      output,
		Scope: bson.M{"batches": batches,
			"types": []string{
				"bdf",
				"diane",
			},
			"deleteOld": []string{"bdf"},
		},
	}

	var err error

	if output == nil {
		_, err = db.C("Entreprise").Find(query).MapReduce(job, &etablissement)
	} else {
		_, err = db.C("Entreprise").Find(query).MapReduce(job, nil)
	}

	if err != nil {
		c.JSON(500, "Echec du traitement MR, message serveur: "+err.Error())
	} else {
		c.JSON(200, etablissement)
	}

}

func dropBatch(c *gin.Context) {
	db := c.Keys["DB"].(*mgo.Database)

	batches := getBatchesID(db)
	sort.Slice(batches, func(i, j int) bool {
		return batches[i] > batches[j]
	})
	currentBatch := batches[0]

	MREtablissement := MapReduceJS{}
	errEt := MREtablissement.load("dropBatch", "etablissement")
	MREntreprise := MapReduceJS{}
	errEn := MREntreprise.load("dropBatch", "entreprise")

	if errEt != nil || errEn != nil {
		c.JSON(500, "Probleme d'accès aux ressources MapReduce")
		return
	}

	jobEt := &mgo.MapReduce{
		Map:      string(MREtablissement.Map),
		Reduce:   string(MREtablissement.Reduce),
		Finalize: string(MREtablissement.Finalize),
		Out:      bson.M{"replace": "Etablissement"},
		Scope:    bson.M{"currentBatch": currentBatch},
	}

	jobEn := &mgo.MapReduce{
		Map:      string(MREntreprise.Map),
		Reduce:   string(MREntreprise.Reduce),
		Finalize: string(MREntreprise.Finalize),
		Out:      bson.M{"replace": "Entreprise"},
		Scope:    bson.M{"currentBatch": currentBatch},
	}

	_, errEt = db.C("Etablissement").Find(nil).MapReduce(jobEt, nil)
	_, errEn = db.C("Entreprise").Find(nil).MapReduce(jobEn, nil)

	c.JSON(200, currentBatch)
}
