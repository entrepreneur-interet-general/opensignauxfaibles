package main

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
)

// ValueEntreprise permet de stocker une entreprise dans un objet Bson
type ValueEntreprise struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Value Entreprise    `json:"value" bson:"value"`
}

// ValueEtablissement structure pour un établissement
type ValueEtablissement struct {
	ID    bson.ObjectId `json:"id" bson:"_id"`
	Value Etablissement `json:"value" bson:"value"`
}

// Etablissement objet établissement (/entreprise/)
type Etablissement struct {
	Siret       string           `json:"siret" bson:"siret"`
	Region      string           `json:"region,omitempty" bson:"region,omitempty"`
	Key         string           `json:"-" bson:"-"`
	AncienSiret []string         `json:"ancien_siret,omitempty" bson:"ancien_siret,omitempty"`
	Batch       map[string]Batch `json:"batch,omitempty" bson:"batch,omitempty"`
}

// Entreprise object Entreprise
type Entreprise struct {
	Siren string           `json:"siren" bson:"siren"`
	Key   string           `json:"-" bson:"-"`
	Batch map[string]Batch `json:"batch,omitempty" bson:"batch,omitempty"`
}

// Batch lot de data
type Batch struct {
	Compact    map[string]bool        `json:"compact,omitempty" bson:"compact,omitempty"`
	Effectif   map[string]*Effectif   `json:"effectif,omitempty" bson:"effectif,omitempty"`
	Delai      map[string]*Delai      `json:"delai,omitempty" bson:"delai,omitempty"`
	Debit      map[string]*Debit      `json:"debit,omitempty" bson:"debit,omitempty"`
	CCSF       map[string]*CCSF       `json:"ccsf,omitempty" bson:"ccsf,omitempty"`
	Cotisation map[string]*Cotisation `json:"cotisation,omitempty" bson:"cotisation,omitempty"`
	Altares    map[string]*Altares    `json:"altares,omitempty" bson:"altares,omitempty"`
  Procol     map[string]*Procol     `json:"procol,omitempty" bson:"procol,omitempty"`
	APDemande  map[string]*APDemande  `json:"apdemande,omitempty" bson:"apdemande,omitempty"`
	APConso    map[string]*APConso    `json:"apconso,omitempty" bson:"apconso,omitempty"`
	Sirene     map[string]*Sirene     `json:"sirene,omitempty" bson:"sirene,omitempty"`
	BDF        map[string]*BDF        `json:"bdf,omitempty" bson:"bdf,omitempty"`
	Diane      map[string]*Diane      `json:"diane,omitempty" bson:"diane,omitempty"`
	DPAE       map[string]*DPAE       `json:"dpae,omitempty" bson:"dpae,omitempty"`
}

func (batch1 Batch) merge(batch2 Batch) {
	if batch1.Altares == nil {
		batch1.Altares = make(map[string]*Altares)
	}
	for hash, altares := range batch2.Altares {
		batch1.Altares[hash] = altares
	}

	if batch1.Procol == nil {
		batch1.Procol = make(map[string]*Procol)
	}
	for hash, procol := range batch2.Procol {
		batch1.Procol[hash] = procol
	}

	if batch1.Cotisation == nil {
		batch1.Cotisation = make(map[string]*Cotisation)
	}
	for hash, cotisation := range batch2.Cotisation {
		batch1.Cotisation[hash] = cotisation
	}

	if batch1.Debit == nil {
		batch1.Debit = make(map[string]*Debit)
	}
	for hash, debit := range batch2.Debit {
		batch1.Debit[hash] = debit
	}

	if batch1.Delai == nil {
		batch1.Delai = make(map[string]*Delai)
	}
	for hash, delai := range batch2.Delai {
		batch1.Delai[hash] = delai
	}

	if batch1.Effectif == nil {
		batch1.Effectif = make(map[string]*Effectif)
	}
	for hash, effectif := range batch2.Effectif {
		batch1.Effectif[hash] = effectif
	}

	if batch1.APDemande == nil {
		batch1.APDemande = make(map[string]*APDemande)
	}
	for hash, apdemande := range batch2.APDemande {
		batch1.APDemande[hash] = apdemande
	}

	if batch1.APConso == nil {
		batch1.APConso = make(map[string]*APConso)
	}
	for hash, apconso := range batch2.APConso {
		batch1.APConso[hash] = apconso
	}

	if batch1.Sirene == nil {
		batch1.Sirene = make(map[string]*Sirene)
	}
	for hash, sirene := range batch2.Sirene {
		batch1.Sirene[hash] = sirene
	}

	if batch1.BDF == nil {
		batch1.BDF = make(map[string]*BDF)
	}
	for hash, bdf := range batch2.BDF {
		batch1.BDF[hash] = bdf
	}

	if batch1.Diane == nil {
		batch1.Diane = make(map[string]*Diane)
	}
	for hash, diane := range batch2.Diane {
		batch1.Diane[hash] = diane
	}

	if batch1.DPAE == nil {
		batch1.DPAE = make(map[string]*DPAE)
	}
	for hash, dpae := range batch2.DPAE {
		batch1.DPAE[hash] = dpae
	}
}

func (value1 ValueEtablissement) merge(value2 ValueEtablissement) (ValueEtablissement, error) {
	if value1.Value.Siret != value2.Value.Siret {
		return ValueEtablissement{},
			errors.New("Valeurs non missibles: sirets '" +
				value1.Value.Siret + "' et '" +
				value2.Value.Siret + "'")
	}
	for idBatch := range value2.Value.Batch {
		if value1.Value.Batch == nil {
			value1.Value.Batch = make(map[string]Batch)
		}
		value1.Value.Batch[idBatch].merge(value2.Value.Batch[idBatch])
	}
	return value1, nil
}

func (value1 ValueEntreprise) merge(value2 ValueEntreprise) (ValueEntreprise, error) {
	if value1.Value.Siren != value2.Value.Siren {
		return ValueEntreprise{},
			errors.New("Valeurs non missibles: sirens '" +
				value1.Value.Siren + "' et '" +
				value2.Value.Siren + "'")
	}
	for idBatch := range value2.Value.Batch {
		if value1.Value.Batch == nil {
			value1.Value.Batch = make(map[string]Batch)
		}
		value1.Value.Batch[idBatch].merge(value2.Value.Batch[idBatch])
	}
	return value1, nil
}

// Periode Période de temps avec un début et une fin
type Periode struct {
	Start time.Time `json:"start" bson:"start"`
	End   time.Time `json:"end" bson:"end"`
}

// Prediction structure des prédiction #FASTDIRTY #FIXME
type Prediction struct {
	ID struct {
		Siret string `json:"siret" bson:"siret"`
		Batch string `json:"batch" bson:"batch"`
		Algo  string `json:"algo" bson:"algo"`
	} `json:"id" bson:"_id"`

	Prob          float64 `json:"prob" bson:"prob"`
	Diff          float64 `json:"diff" bson:"diff"`
	RaiSoc        string  `json:"raison_sociale" bson:"raison_sociale"`
	Departement   string  `json:"departement" bson:"departement"`
	Region        string  `json:"region" bson:"region"`
	EtatProCol    string  `json:"procol" bson:"procol"`
	DefaultUrssaf bool    `json:"default_urssaf" bson:"default_urssaf"`
	Connu         bool    `json:"connu" bson:"connu"`
	Niveau1       string  `json:"naf1" bson:"naf1"`
	Effectif      int     `json:"effectif" bson:"effectif"`
	CCSF          bool    `json:"ccsf" bson:"ccsf"`
}
