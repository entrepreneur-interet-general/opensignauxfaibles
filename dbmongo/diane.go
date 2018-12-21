package main

import (
  "encoding/csv"
  "fmt"
  "io"
  "os/exec"
  "strconv"
  "time"

  "github.com/cnf/structhash"
  "github.com/spf13/viper"
)

// Diane Information financières
type Diane struct {
  Annee                                *int      `json:"exercice_diane" bson:"exercice_diane"`
  NomEntreprise                        string    `json:"nom_entreprise" bson:"nom_entreprise"`
  NumeroSiren                          string    `json:"numero_siren" bson:"numero_siren"`
  StatutJuridique                      string    `json:"statut_juridique" bson:"statut_juridique"`
  ProcedureCollective                  bool      `json:"procedure_collective" bson:"procedure_collective"`
  EffectifConsolide                    *int      `json:"effectif_consolide" bson:"effectif_consolide"`
  DetteFiscaleEtSociale                *float64  `json:"dette_fiscale_et_sociale" bson:"dette_fiscale_et_sociale"`
  FraisDeRetD                          *float64  `json:"frais_de_RetD" bson:"frais_de_RetD"`
  ConcesBrevEtDroitsSim                *float64  `json:"conces_brev_et_droits_sim" bson:"conces_brev_et_droits_sim"`
  NotePreface                          *float64  `json:"note_preface" bson:"note_preface"`
  NombreEtabSecondaire                 *int      `json:"nombre_etab_secondaire" bson:"nombre_etab_secondaire"`
  NombreFiliale                        *int      `json:"nombre_filiale" bson:"nombre_filiale"`
  TailleCompoGroupe                    *int      `json:"taille_compo_groupe" bson:"taille_compo_groupe"`
  ArreteBilan                          time.Time `json:"arrete_bilan_diane" bson:"arrete_bilan_diane"`
  NombreMois                           *int      `json:"nombre_mois" bson:"nombre_mois"`
  ConcoursBancaireCourant              *float64  `json:"concours_bancaire_courant" bson:"concours_bancaire_courant"`
  EquilibreFinancier                   *float64  `json:"equilibre_financier" bson:"equilibre_financier"`
  IndependanceFinanciere               *float64  `json:"independance_financiere" bson:"independance_financiere"`
  Endettement                          *float64  `json:"endettement" bson:"endettement"`
  AutonomieFinanciere                  *float64  `json:"autonomie_financiere" bson:"autonomie_financiere"`
  DegreImmoCorporelle                  *float64  `json:"degre_immo_corporelle" bson:"degre_immo_corporelle"`
  FinancementActifCirculant            *float64  `json:"financement_actif_circulant" bson:"financement_actif_circulant"`
  LiquiditeGenerale                    *float64  `json:"liquidite_generale" bson:"liquidite_generale"`
  LiquiditeReduite                     *float64  `json:"liquidite_reduite" bson:"liquidite_reduite"`
  RotationStocks                       *float64  `json:"rotation_stocks" bson:"rotation_stocks"`
  CreditClient                         *float64  `json:"credit_client" bson:"credit_client"`
  CreditFournisseur                    *float64  `json:"credit_fournisseur" bson:"credit_fournisseur"`
  CAparEffectif                        *float64  `json:"ca_par_effectif" bson:"ca_apar_effectif"`
  TauxInteretFinancier                 *float64  `json:"taux_interet_financier" bson:"taux_interet_financier"`
  TauxInteretSurCA                     *float64  `json:"taux_interet_sur_ca" bson:"taux_interet_sur_ca"`
  EndettementGlobal                    *float64  `json:"endettement_global" bson:"endettement_global"`
  TauxEndettement                      *float64  `json:"taux_endettement" bson:"taux_endettement"`
  CapaciteRemboursement                *float64  `json:"capacite_remboursement" bson:"capacite_remboursement"`
  CapaciteAutofinancement              *float64  `json:"capacite_autofinancement" bson:"capacite_autofinancement"`
  CouvertureCA_FDR                     *float64  `json:"couverture_ca_fdr" bson:"couverture_ca_fdr"`
  CouvertureCA_BesoinFDR               *float64  `json:"couverture_ca_besoin_fdr" bson:"couverture_ca_besoin_fdr"`
  PoidsBFRExploitation                 *float64  `json:"poids_bfr_exploitation" bson:"poids_bfr_exploitation"`
  Exportation                          *float64  `json:"exportation" bson:"exportation"`
  EfficaciteEconomique                 *float64  `json:"efficacite_economique" bson:"efficacite_economique"`
  ProductivitePotentielProduction      *float64  `json:"productivite_potentiel_production" bson:"productivite_potentiel_production"`
  ProductiviteCapitalFinancier         *float64  `json:"productivite_capital_financier" bson:"productivite_capital_financier"`
  ProductiviteCapitalInvesti           *float64  `json:"productivite_capital_investi" bson:"productivite_capital_investi"`
  TauxDInvestissementProductif         *float64  `json:"taux_d_investissement_productif" bson:"taux_d_investissement_productif"`
  RentabiliteEconomique                *float64  `json:"rentabilite_economique" bson:"rentabilite_economique"`
  Performance                          *float64  `json:"performance" bson:"performance"`
  RendementBrutFondsPropres            *float64  `json:"rendement_brut_fonds_propres" bson:"rendement_brut_fonds_propres"`
  RentabiliteNette                     *float64  `json:"rentabilite_nette" bson:"rentabilite_nette"`
  RendementCapitauxPropres             *float64  `json:"rendement_capitaux_propres" bson:"rendement_capitaux_propres"`
  RendementRessourcesDurables          *float64  `json:"rendement_ressources_durables" bson:"rendement_ressources_durables"`
  TauxMargeCommerciale                 *float64  `json:"taux_marge_commerciale" bson:"taux_marge_commerciale"`
  TauxValeurAjoutee                    *float64  `json:"taux_valeur_ajoutee" bson:"taux_valeur_ajoutee"`
  PartSalaries                         *float64  `json:"part_salaries" bson:"part_salaries"`
  PartEtat                             *float64  `json:"part_etat" bson:"part_etat"`
  PartPreteur                          *float64  `json:"part_preteur" bson:"part_preteur"`
  PartAutofinancement                  *float64  `json:"part_autofinancement" bson:"part_autofinancement"`
  CA                                   *float64  `json:"ca" bson:"ca"`
  CAExportation                        *float64  `json:"ca_exportation" bson:"ca_exportation"`
  AchatMarchandises                    *float64  `json:"achat_marchandises" bson:"achat_marchandises"`
  AchatMatieresPremieres               *float64  `json:"achat_matieres_premieres" bson:"achat_matieres_premieres"`
  Production                           *float64  `json:"production" bson:"production"`
  MargeCommerciale                     *float64  `json:"marge_commerciale" bson:"marge_commerciale"`
  Consommation                         *float64  `json:"consommation" bson:"consommation"`
  AutresAchatsChargesExternes          *float64  `json:"autres_achats_charges_externes" bson:"autres_achats_charges_externes"`
  ValeurAjoutee                        *float64  `json:"valeur_ajoutee" bson:"valeur_ajoutee"`
  ChargePersonnel                      *float64  `json:"charge_personnel" bson:"charge_personnel"`
  ImpotsTaxes                          *float64  `json:"impots_taxes" bson:"impots_taxes"`
  SubventionsDExploitation             *float64  `json:"subventions_d_exploitation" bson:"subventions_d_exploitation"`
  ExcedentBrutDExploitation            *float64  `json:"excedent_brut_d_exploitation" bson:"excedent_brut_d_exploitation"`
  AutresProduitsChargesReprises        *float64  `json:"autres_produits_charges_reprises" bson:"autres_produits_charges_reprises"`
  DotationAmortissement                *float64  `json:"dotation_amortissement" bson:"dotation_amortissement"`
  ResultatExpl                         *float64  `json:"resultat_expl" bson:"resultat_expl"`
  OperationsCommun                     *float64  `json:"operations_commun" bson:"operations_commun"`
  ProduitsFinanciers                   *float64  `json:"produits_financiers" bson:"produits_financiers"`
  ChargesFinancieres                   *float64  `json:"charges_financieres" bson:"charges_financieres"`
  Interets                             *float64  `json:"interets" bson:"interets"`
  ResultatAvantImpot                   *float64  `json:"resultat_avant_impot" bson:"resultat_avant_impot"`
  ProduitExceptionnel                  *float64  `json:"produit_exceptionnel" bson:"produit_exceptionnel"`
  ChargeExceptionnelle                 *float64  `json:"charge_exceptionnelle" bson:"charge_exceptionnelle"`
  ParticipationSalaries                *float64  `json:"participation_salaries" bson:"participation_salaries"`
  ImpotBenefice                        *float64  `json:"impot_benefice" bson:"impot_benefice"`
  BeneficeOuPerte                      *float64  `json:"benefice_ou_perte" bson:"benefice_ou_perte"`
}

func parseDiane(paths []string) chan *Diane {
  outputChannel := make(chan *Diane)

  go func() {
    n := 0
    var full_paths []string
    for _, path := range paths {
      full_paths = append(full_paths, viper.GetString("APP_DATA")+path)
    }

    if (len(full_paths) != 0){
      full_paths = append([]string{viper.GetString("APP_DATA") + "/../script/convert_diane.sh"}, full_paths...)

      cmd := exec.Command("/bin/bash", full_paths...)
      stdout, err := cmd.StdoutPipe()

      if err != nil {
        fmt.Println("Error", err)
        log(critical, "importDiane", "Erreur à l'ouverture des fichiers Dianes, abandon.")
      } else {
        log(debug, "importDiane", "Ouverture du fichier Diane.")
      }
      reader := csv.NewReader(stdout)
      reader.Comma = ';'
      reader.LazyQuotes = true
      cmd.Start()

      for {
        row, err := reader.Read()
        if err == io.EOF {
          log(debug, "importDiane", "Traitement des fichiers diane terminé: "+fmt.Sprint(n)+" lignes traitées")
          break
        } else if err != nil {
          log(critical, "importDiane", "Erreur lors de la lecture du fichier diane: "+err.Error()+". Abandon !")
          break
        }
        n++
        diane := Diane{}

        if i, err := strconv.Atoi(row[0]); err == nil {
          diane.Annee = &i
        }
        if 2 < len(row){
          diane.NomEntreprise = row[2]
        } else {
          log(critical, "importDiane", "Ligne invalide. Abandon !")
          break
        }
        if 3 < len(row){
          diane.NumeroSiren = row[3]
        } else {
          log(critical, "importDiane", "Ligne invalide. Abandon !")
          break
        }
        if 4 < len(row){
          diane.StatutJuridique = row[4]
        } else {
          log(critical, "importDiane", "Ligne invalide. Abandon !")
          break
        }
        if 5 < len(row){
          diane.ProcedureCollective = (row[5] == "Oui")
        } else {
          log(critical, "importDiane", "Ligne invalide. Abandon !")
          break
        }

        if i, err := strconv.Atoi(row[6]); err == nil {
          diane.EffectifConsolide = &i
        }
        if i, err := strconv.ParseFloat(row[7], 64); err == nil {
          diane.DetteFiscaleEtSociale = &i
        }
        if i, err := strconv.ParseFloat(row[8], 64); err == nil {
          diane.FraisDeRetD = &i
        }
        if i, err := strconv.ParseFloat(row[9], 64); err == nil {
          diane.ConcesBrevEtDroitsSim = &i
        }
        if i, err := strconv.ParseFloat(row[10], 64); err == nil {
          diane.NotePreface = &i
        }
        if i, err := strconv.Atoi(row[11]); err == nil {
          diane.NombreEtabSecondaire = &i
        }
        if i, err := strconv.Atoi(row[12]); err == nil {
          diane.NombreFiliale = &i
        }
        if i, err := strconv.Atoi(row[13]); err == nil {
          diane.TailleCompoGroupe = &i
        }
        if i, err := time.Parse("02/01/2006", row[15]); err == nil {
          diane.ArreteBilan = i
        }
        if i, err := strconv.Atoi(row[16]); err == nil {
          diane.NombreMois = &i
        }
        if i, err := strconv.ParseFloat(row[17], 64); err == nil {
          diane.ConcoursBancaireCourant = &i
        }
        if i, err := strconv.ParseFloat(row[18], 64); err == nil {
          diane.EquilibreFinancier = &i
        }
        if i, err := strconv.ParseFloat(row[19], 64); err == nil {
          diane.IndependanceFinanciere = &i
        }
        if i, err := strconv.ParseFloat(row[20], 64); err == nil {
          diane.Endettement = &i
        }
        if i, err := strconv.ParseFloat(row[21], 64); err == nil {
          diane.AutonomieFinanciere = &i
        }
        if i, err := strconv.ParseFloat(row[22], 64); err == nil {
          diane.DegreImmoCorporelle = &i
        }
        if i, err := strconv.ParseFloat(row[23], 64); err == nil {
          diane.FinancementActifCirculant = &i
        }
        if i, err := strconv.ParseFloat(row[24], 64); err == nil {
          diane.LiquiditeGenerale = &i
        }
        if i, err := strconv.ParseFloat(row[25], 64); err == nil {
          diane.LiquiditeReduite = &i
        }
        if i, err := strconv.ParseFloat(row[26], 64); err == nil {
          diane.RotationStocks = &i
        }
        if i, err := strconv.ParseFloat(row[27], 64); err == nil {
          diane.CreditClient = &i
        }
        if i, err := strconv.ParseFloat(row[28], 64); err == nil {
          diane.CreditFournisseur = &i
        }
        if i, err := strconv.ParseFloat(row[29], 64); err == nil {
          diane.CAparEffectif = &i
        }
        if i, err := strconv.ParseFloat(row[30], 64); err == nil {
          diane.TauxInteretFinancier = &i
        }
        if i, err := strconv.ParseFloat(row[31], 64); err == nil {
          diane.TauxInteretSurCA = &i
        }
        if i, err := strconv.ParseFloat(row[32], 64); err == nil {
          diane.EndettementGlobal = &i
        }
        if i, err := strconv.ParseFloat(row[33], 64); err == nil {
          diane.TauxEndettement = &i
        }
        if i, err := strconv.ParseFloat(row[34], 64); err == nil {
          diane.CapaciteRemboursement = &i
        }
        if i, err := strconv.ParseFloat(row[35], 64); err == nil {
          diane.CapaciteAutofinancement = &i
        }
        if i, err := strconv.ParseFloat(row[36], 64); err == nil {
          diane.CouvertureCA_FDR = &i
        }
        if i, err := strconv.ParseFloat(row[37], 64); err == nil {
          diane.CouvertureCA_BesoinFDR = &i
        }
        if i, err := strconv.ParseFloat(row[38], 64); err == nil {
          diane.PoidsBFRExploitation = &i
        }
        if i, err := strconv.ParseFloat(row[39], 64); err == nil {
          diane.Exportation = &i
        }
        if i, err := strconv.ParseFloat(row[40], 64); err == nil {
          diane.EfficaciteEconomique = &i
        }
        if i, err := strconv.ParseFloat(row[41], 64); err == nil {
          diane.ProductivitePotentielProduction = &i
        }
        if i, err := strconv.ParseFloat(row[42], 64); err == nil {
          diane.ProductiviteCapitalFinancier = &i
        }
        if i, err := strconv.ParseFloat(row[43], 64); err == nil {
          diane.ProductiviteCapitalInvesti = &i
        }
        if i, err := strconv.ParseFloat(row[44], 64); err == nil {
          diane.TauxDInvestissementProductif = &i
        }
        if i, err := strconv.ParseFloat(row[45], 64); err == nil {
          diane.RentabiliteEconomique = &i
        }
        if i, err := strconv.ParseFloat(row[46], 64); err == nil {
          diane.Performance = &i
        }
        if i, err := strconv.ParseFloat(row[47], 64); err == nil {
          diane.RendementBrutFondsPropres = &i
        }
        if i, err := strconv.ParseFloat(row[48], 64); err == nil {
          diane.RentabiliteNette = &i
        }
        if i, err := strconv.ParseFloat(row[49], 64); err == nil {
          diane.RendementCapitauxPropres = &i
        }
        if i, err := strconv.ParseFloat(row[50], 64); err == nil {
          diane.RendementRessourcesDurables = &i
        }
        if i, err := strconv.ParseFloat(row[51], 64); err == nil {
          diane.TauxMargeCommerciale = &i
        }
        if i, err := strconv.ParseFloat(row[52], 64); err == nil {
          diane.TauxValeurAjoutee = &i
        }
        if i, err := strconv.ParseFloat(row[53], 64); err == nil {
          diane.PartSalaries = &i
        }
        if i, err := strconv.ParseFloat(row[54], 64); err == nil {
          diane.PartEtat = &i
        }
        if i, err := strconv.ParseFloat(row[55], 64); err == nil {
          diane.PartPreteur = &i
        }
        if i, err := strconv.ParseFloat(row[56], 64); err == nil {
          diane.PartAutofinancement = &i
        }
        if i, err := strconv.ParseFloat(row[57], 64); err == nil {
          diane.CA = &i
        }
        if i, err := strconv.ParseFloat(row[58], 64); err == nil {
          diane.CAExportation = &i
        }
        if i, err := strconv.ParseFloat(row[59], 64); err == nil {
          diane.AchatMarchandises = &i
        }
        if i, err := strconv.ParseFloat(row[61], 64); err == nil {
          diane.AchatMatieresPremieres = &i
        }
        if i, err := strconv.ParseFloat(row[62], 64); err == nil {
          diane.Production = &i
        }
        if i, err := strconv.ParseFloat(row[63], 64); err == nil {
          diane.MargeCommerciale = &i
        }
        if i, err := strconv.ParseFloat(row[64], 64); err == nil {
          diane.Consommation = &i
        }
        if i, err := strconv.ParseFloat(row[65], 64); err == nil {
          diane.AutresAchatsChargesExternes = &i
        }
        if i, err := strconv.ParseFloat(row[66], 64); err == nil {
          diane.ValeurAjoutee = &i
        }
        if i, err := strconv.ParseFloat(row[67], 64); err == nil {
          diane.ChargePersonnel = &i
        }
        if i, err := strconv.ParseFloat(row[68], 64); err == nil {
          diane.ImpotsTaxes = &i
        }
        if i, err := strconv.ParseFloat(row[69], 64); err == nil {
          diane.SubventionsDExploitation = &i
        }
        if i, err := strconv.ParseFloat(row[70], 64); err == nil {
          diane.ExcedentBrutDExploitation = &i
        }
        if i, err := strconv.ParseFloat(row[71], 64); err == nil {
          diane.AutresProduitsChargesReprises = &i
        }
        if i, err := strconv.ParseFloat(row[72], 64); err == nil {
          diane.DotationAmortissement = &i
        }
        if i, err := strconv.ParseFloat(row[73], 64); err == nil {
          diane.ResultatExpl = &i
        }
        if i, err := strconv.ParseFloat(row[74], 64); err == nil {
          diane.OperationsCommun = &i
        }
        if i, err := strconv.ParseFloat(row[75], 64); err == nil {
          diane.ProduitsFinanciers = &i
        }
        if i, err := strconv.ParseFloat(row[76], 64); err == nil {
          diane.ChargesFinancieres = &i
        }
        if i, err := strconv.ParseFloat(row[77], 64); err == nil {
          diane.Interets = &i
        }
        if i, err := strconv.ParseFloat(row[78], 64); err == nil {
          diane.ResultatAvantImpot = &i
        }
        if i, err := strconv.ParseFloat(row[79], 64); err == nil {
          diane.ProduitExceptionnel = &i
        }
        if i, err := strconv.ParseFloat(row[80], 64); err == nil {
          diane.ChargeExceptionnelle = &i
        }
        if i, err := strconv.ParseFloat(row[81], 64); err == nil {
          diane.ParticipationSalaries = &i
        }
        if i, err := strconv.ParseFloat(row[82], 64); err == nil {
          diane.ImpotBenefice = &i
        }
        if i, err := strconv.ParseFloat(row[83], 64); err == nil {
          diane.BeneficeOuPerte = &i
        }

        outputChannel <- &diane

      }
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
