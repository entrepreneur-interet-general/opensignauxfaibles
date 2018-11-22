package main

import (
	"testing"
)

func Test_nextBatchID(t *testing.T) {
	batchID, err := nextBatchID("abcd")

	if err == nil {
		t.Error("Erreur attendue absente. '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test valeur erronée ok")
	}
	batchID, err = nextBatchID("1801")
	if err != nil || batchID != "1802" {
		t.Error("'1802' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test valeur courante ok")
	}
	batchID, err = nextBatchID("1812")
	if err != nil || batchID != "1901" {
		t.Error("'1901' attendu, '" + batchID + "' obtenu. err = " + err.Error())
	} else {
		t.Log("Test passage nouvelle année ok")
	}
}

func Test_isBatchID(t *testing.T) {
	if isBatchID("1801") {
		t.Log("1801 est un ID de batch")
	} else {
		t.Error("1801 devrait être un ID de batch")
	}

	if isBatchID("") {
		t.Error("'' ne devrait pas être considéré comme un ID de batch")
	} else {
		t.Log("'' est bien rejeté")
	}

	if isBatchID("190193039") {
		t.Error("'190193039' ne devrait pas être considéré comme un ID de batch")
	} else {
		t.Log("'190193039' est bien rejeté: ")
	}

	if isBatchID("abcd") {
		t.Error("'abcd' ne devrait pas être considéré comme un ID de batch")
	} else {
		t.Log("'abcd' est bien rejeté: ")
	}
}
