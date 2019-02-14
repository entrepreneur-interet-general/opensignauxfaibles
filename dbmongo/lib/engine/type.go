package engine

// GetTypes retourne la liste des types déclarés
func GetTypes() []Type {
	return []Type{
		{"admin_urssaf", "Siret/Compte URSSAF", "Liste comptes"},
		{"apconso", "Consommation Activité Partielle", "conso"},
		{"bdf", "Ratios Banque de France", "bdf"},
		{"cotisation", "Cotisations URSSAF", "cotisation"},
		{"delai", "Délais URSSAF", "delais|Délais"},
		{"dpae", "Déclaration Préalable à l'embauche", "DPAE"},
		{"interim", "Base Interim", "interim"},
		{"altares", "Base Altarès", "ALTARES"},
		{"procol", "Procédures collectives", "procol"},
		{"apdemande", "Demande Activité Partielle", "dde"},
		{"ccsf", "Stock CCSF à date", "ccsf"},
		{"debit", "Débits URSSAF", "debit"},
		{"dmmo", "Déclaration Mouvement de Main d'Œuvre", "dmmo"},
		{"effectif", "Emplois URSSAF", "Emploi"},
		{"sirene", "Base GéoSirene", "sirene"},
		{"diane", "Diane", "diane"},
	}
}

// Type description des types de fichiers pris en charge
type Type struct {
	Type    string `json:"type" bson:"type"`
	Libelle string `json:"text" bson:"text"`
	Filter  string `json:"filter" bson:"filter"`
}

// Parser fonction de traitement de données en entrée
type Parser func(AdminBatch) (chan Tuple, chan Event)

// Browseable est le type qui permet d'envoyer les objets vers le frontend
// Voir la fonction Browse
type Browseable struct {
	ID struct {
		Key   string   `json:"key" bson:"key"`
		Scope []string `json:"scope" bson:"scope"`
	}
	Value map[string]interface{} `json:"value" bson:"value"`
}
