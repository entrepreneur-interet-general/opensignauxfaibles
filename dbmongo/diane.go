package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/cnf/structhash"
	"github.com/spf13/viper"
)

// Diane Information financières
type Diane struct {
	Annee                                *int     `json:"exercice_diane" bson:"exercice_diane"`
	Marquee                              *float64 `json:"marquee" bson:"marquee"`
	NomEntreprise                        string   `json:"nom_entreprise" bson:"nom_entreprise"`
	NumeroSiren                          string   `json:"numero_siren" bson:"numero_siren"`
	StatutJuridique                      string   `json:"statut_juridique" bson:"statut_juridique"`
	ProcedureCollective                  bool     `json:"procedure_collective" bson:"procedure_collective"`
	CA                                   *float64 `json:"CA" bson:"CA"`
	ValeurAjoutee                        *float64 `json:"valeur_ajoutee" bson:"valeur_ajoutee"`
	ResultatNetConsolide                 *float64 `json:"resultat_net_consolide" bson:"resultat_net_consolide"`
	CapaciteAutofinancAvantRepartition   *float64 `json:"capacite_autofinanc_avant_repartition" bson:"capacite_autofinanc_avant_repartition"`
	CapitalSocialOuIndividuel            *float64 `json:"capital_social_ou_individuel" bson:"capital_social_ou_individuel"`
	CapitauxPropresDuGroupe              *float64 `json:"capitaux_propres_du_groupe" bson:"capitaux_propres_du_groupe"`
	FondsDeRoulNetGlobal                 *float64 `json:"fonds_de_roul_net_global" bson:"fonds_de_roul_net_global"`
	EndettementPourcent                  *float64 `json:"endettement_pourcent" bson:"endettement_pourcent"`
	LiquiditeReduite                     *float64 `json:"liquidite_reduite" bson:"liquidite_reduite"`
	RentabiliteNettePourcent             *float64 `json:"rentabilite_nette_pourcent" bson:"rentabilite_nette_pourcent"`
	RendDesCapitauxPropresNetsPourcent   *float64 `json:"rend_des_capitaux_propres_nets_pourcent" bson:"rend_des_capitaux_propres_nets_pourcent"`
	RendDesRessDurablesNettesPourcent    *float64 `json:"rend_des_ress_durables_nettes_pourcent" bson:"rend_des_ress_durables_nettes_pourcent"`
	EffectifConsolide                    *float64 `json:"effectif_consolide" bson:"effectif_consolide"`
	TotalActifImmob                      *float64 `json:"total_actif_immob" bson:"total_actif_immob"`
	TotalImmobFin                        *float64 `json:"total_immob_fin" bson:"total_immob_fin"`
	TotalImmobCorp                       *float64 `json:"total_immob_corp" bson:"total_immob_corp"`
	TotalImmobIncorp                     *float64 `json:"total_immob_incorp" bson:"total_immob_incorp"`
	Stocks                               *float64 `json:"stocks" bson:"stocks"`
	ProduitsIntermedEtFinis              *float64 `json:"produits_intermed_et_finis" bson:"produits_intermed_et_finis"`
	Marchandises                         *float64 `json:"marchandises" bson:"marchandises"`
	EnCoursDeProdDeBiens                 *float64 `json:"en_cours_de_prod_de_biens" bson:"en_cours_de_prod_de_biens"`
	MatieresPremApprov                   *float64 `json:"matieres_prem_approv" bson:"matieres_prem_approv"`
	CreancesExpl                         *float64 `json:"creances_expl" bson:"creances_expl"`
	ClientsEtCptesRatt                   *float64 `json:"clients_et_cptes_ratt" bson:"clients_et_cptes_ratt"`
	AvEtAcSurCommandes                   *float64 `json:"av_et_ac_sur_commandes" bson:"av_et_ac_sur_commandes"`
	Disponibilites                       *float64 `json:"disponibilites" bson:"disponibilites"`
	TotalActifCircChConstAv              *float64 `json:"total_actif_circ_ch_const_av" bson:"total_actif_circ_ch_const_av"`
	TotalActif                           *float64 `json:"total_actif" bson:"total_actif"`
	CapitauxPropresGroupe                *float64 `json:"capitaux_propres_groupe" bson:"capitaux_propres_groupe"`
	ResultatConsolidePartDuGroupe        *float64 `json:"resultat_consolide_part_du_groupe" bson:"resultat_consolide_part_du_groupe"`
	TotalDettesFin                       *float64 `json:"total_dettes_fin" bson:"total_dettes_fin"`
	TotalDetteExplEtDivers               *float64 `json:"total_dette_expl_et_divers" bson:"total_dette_expl_et_divers"`
	DettesFournEtCptesRatt               *float64 `json:"dettes_fourn_et_cptes_ratt" bson:"dettes_fourn_et_cptes_ratt"`
	DettesFiscalesEtSociales             *float64 `json:"dettes_fiscales_et_sociales" bson:"dettes_fiscales_et_sociales"`
	DettesSurImmobCptesRatt              *float64 `json:"dettes_sur_immob_cptes_ratt" bson:"dettes_sur_immob_cptes_ratt"`
	TotalDuPassif                        *float64 `json:"total_du_passif" bson:"total_du_passif"`
	ChiffreAffairesNet                   *float64 `json:"chiffre_affaires_net" bson:"chiffre_affaires_net"`
	ChiffreAffairesNetEnFrance           *float64 `json:"chiffre_affaires_net_en_france" bson:"chiffre_affaires_net_en_france"`
	ChiffreAffairesNetLieAuxExportations *float64 `json:"chiffre_affaires_net_lie_aux_exportations" bson:"chiffre_affaires_net_lie_aux_exportations"`
	SalairesEtTraitements                *float64 `json:"salaires_et_traitements" bson:"salaires_et_traitements"`
	ChargesSociales                      *float64 `json:"charges_sociales" bson:"charges_sociales"`
	TotalDesChargesExpl                  *float64 `json:"total_des_charges_expl" bson:"total_des_charges_expl"`
	ResultatExpl                         *float64 `json:"resultat_expl" bson:"resultat_expl"`
	TotalDesProduitsFin                  *float64 `json:"total_des_produits_fin" bson:"total_des_produits_fin"`
	TotalDesChargesFin                   *float64 `json:"total_des_charges_fin" bson:"total_des_charges_fin"`
	ResultatCourantAvantImpots           *float64 `json:"resultat_courant_avant_impots" bson:"resultat_courant_avant_impots"`
	ResultatFinancier                    *float64 `json:"resultat_financier" bson:"resultat_financier"`
	ResultatExceptionnel                 *float64 `json:"resultat_exceptionnel" bson:"resultat_exceptionnel"`
	TotalDesCharges                      *float64 `json:"total_des_charges" bson:"total_des_charges"`
	TotalDesProduits                     *float64 `json:"total_des_produits" bson:"total_des_produits"`
	FraisDeRetD                          *float64 `json:"frais_de_RetD" bson:"frais_de_RetD"`
	ConcesBrevEtDroitsSim                *float64 `json:"conces_brev_et_droits_sim" bson:"conces_brev_et_droits_sim"`
	NotePreface                          *float64 `json:"note_preface" bson:"note_preface"`
	NombreEtabSecondaire                 *int     `json:"nombre_etab_secondaire" bson:"nombre_etab_secondaire"`
}

func parseDiane(paths []string) chan *Diane {
	outputChannel := make(chan *Diane)

	go func() {
		for _, path := range paths {
			n := 0

			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				fmt.Println("Error", err)
				log(critical, "importDiane", "Erreur à l'ouverture du fichier "+path+", abandon.")
			} else {
				log(debug, "importDiane", "Ouverture du fichier "+path+".")
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.FieldsPerRecord = 69

			for {
				row, error := reader.Read()

				if error == io.EOF {
					log(debug, "importDiane", "Traitement du fichier "+path+" terminé: "+fmt.Sprint(n)+" lignes traitées")
					break
				} else if error != nil {
					log(critical, "importDiane", "Erreur lors de la lecture du fichier "+path+": "+err.Error()+". Abandon !")
					break
				}
				n++
				diane := Diane{}
				if i, err := strconv.Atoi(row[0]); err == nil {
					diane.Annee = &i
				}
				if i, err := strconv.ParseFloat(row[1], 64); err == nil {
					diane.Marquee = &i
				}

				diane.NomEntreprise = row[2]
				diane.NumeroSiren = row[3]
				diane.StatutJuridique = row[4]
				diane.ProcedureCollective = (row[5] == "Oui")
				if i, err := strconv.ParseFloat(row[6], 64); err == nil {
					diane.CA = &i
				}
				if i, err := strconv.ParseFloat(row[7], 64); err == nil {
					diane.ValeurAjoutee = &i
				}
				if i, err := strconv.ParseFloat(row[8], 64); err == nil {
					diane.ResultatNetConsolide = &i
				}
				if i, err := strconv.ParseFloat(row[9], 64); err == nil {
					diane.CapaciteAutofinancAvantRepartition = &i
				}
				if i, err := strconv.ParseFloat(row[10], 64); err == nil {
					diane.CapitalSocialOuIndividuel = &i
				}
				if i, err := strconv.ParseFloat(row[11], 64); err == nil {
					diane.CapitauxPropresDuGroupe = &i
				}
				if i, err := strconv.ParseFloat(row[12], 64); err == nil {
					diane.FondsDeRoulNetGlobal = &i
				}
				if i, err := strconv.ParseFloat(row[13], 64); err == nil {
					diane.EndettementPourcent = &i
				}
				if i, err := strconv.ParseFloat(row[14], 64); err == nil {
					diane.LiquiditeReduite = &i
				}
				if i, err := strconv.ParseFloat(row[15], 64); err == nil {
					diane.RentabiliteNettePourcent = &i
				}
				if i, err := strconv.ParseFloat(row[16], 64); err == nil {
					diane.RendDesCapitauxPropresNetsPourcent = &i
				}
				if i, err := strconv.ParseFloat(row[17], 64); err == nil {
					diane.RendDesRessDurablesNettesPourcent = &i
				}
				if i, err := strconv.ParseFloat(row[18], 64); err == nil {
					diane.EffectifConsolide = &i
				}
				if i, err := strconv.ParseFloat(row[19], 64); err == nil {
					diane.TotalActifImmob = &i
				}
				if i, err := strconv.ParseFloat(row[20], 64); err == nil {
					diane.TotalImmobFin = &i
				}
				if i, err := strconv.ParseFloat(row[21], 64); err == nil {
					diane.TotalImmobCorp = &i
				}
				if i, err := strconv.ParseFloat(row[22], 64); err == nil {
					diane.TotalImmobIncorp = &i
				}
				if i, err := strconv.ParseFloat(row[23], 64); err == nil {
					diane.Stocks = &i
				}
				if i, err := strconv.ParseFloat(row[24], 64); err == nil {
					diane.ProduitsIntermedEtFinis = &i
				}
				if i, err := strconv.ParseFloat(row[25], 64); err == nil {
					diane.Marchandises = &i
				}
				if i, err := strconv.ParseFloat(row[26], 64); err == nil {
					diane.EnCoursDeProdDeBiens = &i
				}
				if i, err := strconv.ParseFloat(row[27], 64); err == nil {
					diane.MatieresPremApprov = &i
				}
				if i, err := strconv.ParseFloat(row[28], 64); err == nil {
					diane.CreancesExpl = &i
				}
				if i, err := strconv.ParseFloat(row[29], 64); err == nil {
					diane.ClientsEtCptesRatt = &i
				}
				if i, err := strconv.ParseFloat(row[30], 64); err == nil {
					diane.AvEtAcSurCommandes = &i
				}
				if i, err := strconv.ParseFloat(row[31], 64); err == nil {
					diane.Disponibilites = &i
				}
				if i, err := strconv.ParseFloat(row[32], 64); err == nil {
					diane.TotalActifCircChConstAv = &i
				}
				if i, err := strconv.ParseFloat(row[33], 64); err == nil {
					diane.TotalActif = &i
				}
				if i, err := strconv.ParseFloat(row[35], 64); err == nil {
					diane.CapitauxPropresGroupe = &i
				}
				if i, err := strconv.ParseFloat(row[36], 64); err == nil {
					diane.ResultatConsolidePartDuGroupe = &i
				}
				if i, err := strconv.ParseFloat(row[37], 64); err == nil {
					diane.TotalDettesFin = &i
				}
				if i, err := strconv.ParseFloat(row[38], 64); err == nil {
					diane.TotalDetteExplEtDivers = &i
				}
				if i, err := strconv.ParseFloat(row[39], 64); err == nil {
					diane.DettesFournEtCptesRatt = &i
				}
				if i, err := strconv.ParseFloat(row[40], 64); err == nil {
					diane.DettesFiscalesEtSociales = &i
				}
				if i, err := strconv.ParseFloat(row[41], 64); err == nil {
					diane.DettesSurImmobCptesRatt = &i
				}
				if i, err := strconv.ParseFloat(row[42], 64); err == nil {
					diane.TotalDuPassif = &i
				}
				if i, err := strconv.ParseFloat(row[43], 64); err == nil {
					diane.ChiffreAffairesNet = &i
				}
				if i, err := strconv.ParseFloat(row[44], 64); err == nil {
					diane.ChiffreAffairesNetEnFrance = &i
				}
				if i, err := strconv.ParseFloat(row[45], 64); err == nil {
					diane.ChiffreAffairesNetLieAuxExportations = &i
				}
				if i, err := strconv.ParseFloat(row[46], 64); err == nil {
					diane.SalairesEtTraitements = &i
				}
				if i, err := strconv.ParseFloat(row[47], 64); err == nil {
					diane.ChargesSociales = &i
				}
				if i, err := strconv.ParseFloat(row[48], 64); err == nil {
					diane.TotalDesChargesExpl = &i
				}
				if i, err := strconv.ParseFloat(row[49], 64); err == nil {
					diane.ResultatExpl = &i
				}
				if i, err := strconv.ParseFloat(row[50], 64); err == nil {
					diane.TotalDesProduitsFin = &i
				}
				if i, err := strconv.ParseFloat(row[51], 64); err == nil {
					diane.TotalDesChargesFin = &i
				}
				if i, err := strconv.ParseFloat(row[52], 64); err == nil {
					diane.ResultatCourantAvantImpots = &i
				}
				if i, err := strconv.ParseFloat(row[53], 64); err == nil {
					diane.ResultatFinancier = &i
				}
				if i, err := strconv.ParseFloat(row[54], 64); err == nil {
					diane.ResultatExceptionnel = &i
				}
				if i, err := strconv.ParseFloat(row[56], 64); err == nil {
					diane.TotalDesCharges = &i
				}
				if i, err := strconv.ParseFloat(row[57], 64); err == nil {
					diane.TotalDesProduits = &i
				}
				if i, err := strconv.ParseFloat(row[58], 64); err == nil {
					diane.FraisDeRetD = &i
				}
				if i, err := strconv.ParseFloat(row[59], 64); err == nil {
					diane.ConcesBrevEtDroitsSim = &i
				}
				if i, err := strconv.ParseFloat(row[60], 64); err == nil {
					diane.NotePreface = &i
				}
				if i, err := strconv.Atoi(row[61]); err == nil {
					diane.NombreEtabSecondaire = &i
				}
				outputChannel <- &diane

			}
			file.Close()
		}
		close(outputChannel)
	}()

	return outputChannel
}

func importDiane(batch *AdminBatch) error {
	log(info, "importDiane", "Import batch "+batch.ID.Key+": début import Diane")
	for diane := range parseDiane(batch.Files["diane"]) {
		hash := fmt.Sprintf("%x", structhash.Md5(diane, 1))

		value := ValueEntreprise{
			Value: Entreprise{
				Siren: diane.NumeroSiren,
				Batch: map[string]Batch{
					batch.ID.Key: Batch{
						Diane: map[string]*Diane{
							hash: diane,
						}}}}}
		db.ChanEntreprise <- &value
	}
	db.ChanEntreprise <- &ValueEntreprise{}
	log(info, "importDiane", "Import batch "+batch.ID.Key+": fin import Diane")
	return nil
}
