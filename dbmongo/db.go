package main

import (
	"errors"
	"regexp"
	"time"

	"github.com/spf13/viper"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

// DB contient des méthodes pour accéder à la base de données
type DB struct {
	DB                *mgo.Database
	DBStatus          *mgo.Database
	Status            Status
	ChanEntreprise    chan *ValueEntreprise
	ChanEtablissement chan *ValueEtablissement
}

// DB Initialisation de la connexion MongoDB
func initDB() DB {
	loadConfig()

	dbDial := viper.GetString("DB_DIAL")
	dbDatabase := viper.GetString("DB")

	// définition de 3 connexions pour isoler les curseurs
	mongostatus, err := mgo.Dial(dbDial)
	if err != nil {
		// log.Panic(err)
	}

	mongodb, err := mgo.Dial(dbDial)
	if err != nil {
		// log.Panic(err)
	}

	mongostatus.SetSocketTimeout(3600 * time.Second)
	mongodb.SetSocketTimeout(3600 * time.Second)
	dbstatus := mongostatus.DB(dbDatabase)
	db := mongodb.DB(dbDatabase)

	firstBatchID := viper.GetString("FIRST_BATCH")
	if !isBatchID(firstBatchID) {
		panic("Paramètre FIRST_BATCH incorrect, vérifiez la configuration.")
	}

	// firstBatch, err := getBatch(db, firstBatchID)
	var firstBatch AdminBatch
	err = db.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": firstBatchID}).One(&firstBatch)

	if firstBatch.ID.Type == "" {
		firstBatch = AdminBatch{
			ID: AdminID{
				Key:  firstBatchID,
				Type: "batch",
			},
		}

		err := db.C("Admin").Insert(firstBatch)

		if err != nil {
			panic("Impossible de créer le premier batch: " + err.Error())
		}
	}

	status := Status{
		DB: dbstatus,
		ID: AdminID{
			Key:  "status",
			Type: "status",
		},
	}
	status.write()

	// pousse les fonctions partagées JS
	err = declareServerFunctions(db)
	if err != nil {
		// log.Panic(err)
	}

	chanEntreprise := insertEntreprise(db)
	chanEtablissement := insertEtablissement(db)

	// envoie un struct vide pour purger les channels au cas où il reste les objets non insérés
	go func() {
		for range time.Tick(1 * time.Second) {
			chanEntreprise <- &ValueEntreprise{}
			chanEtablissement <- &ValueEtablissement{}
		}
	}()

	return DB{
		DB:                db,
		DBStatus:          dbstatus,
		ChanEntreprise:    chanEntreprise,
		ChanEtablissement: chanEtablissement,
		Status:            status,
	}
}

func logErrors(db *mgo.Database, err error) {
	db.C("Journal").Insert(err)
}

func insertEntreprise(db *mgo.Database) chan *ValueEntreprise {
	source := make(chan *ValueEntreprise, 1000)

	go func(chan *ValueEntreprise) {
		buffer := make(map[string]*ValueEntreprise)
		objects := make([]interface{}, 0)
		i := 0

		for value := range source {
			if (value.Value.Batch == nil) || i >= 100 {
				for _, v := range buffer {
					objects = append(objects, *v)
				}
				if len(objects) > 0 {
					go func(o []interface{}) { db.C("Entreprise").Insert(o...) }(objects)
				}
				buffer = make(map[string]*ValueEntreprise)
				objects = make([]interface{}, 0)
				i = 0
			}
			if match, _ := regexp.MatchString("^[0-9]{9}$", value.Value.Siren); value.Value.Batch != nil && match {
				if knownValue, ok := buffer[value.Value.Siren]; ok {
					newValue, _ := (*knownValue).merge(*value)
					buffer[value.Value.Siren] = &newValue
				} else {
					value.ID = bson.NewObjectId()
					buffer[value.Value.Siren] = value
					i++
				}
			}
		}
	}(source)

	return source
}

func insertEtablissement(db *mgo.Database) chan *ValueEtablissement {
	source := make(chan *ValueEtablissement, 1000)

	go func(chan *ValueEtablissement) {
		buffer := make(map[string]*ValueEtablissement)
		objects := make([]interface{}, 0)
		i := 0

		for value := range source {
			if (value.Value.Batch == nil) || i >= 100 {
				for _, v := range buffer {
					objects = append(objects, *v)
				}
				if len(objects) > 0 {
					go func(o []interface{}) { db.C("Etablissement").Insert(o...) }(objects)
				}
				buffer = make(map[string]*ValueEtablissement)
				objects = make([]interface{}, 0)
				i = 0
			}
			if match, _ := regexp.MatchString("^[0-9]{14}$", value.Value.Siret); value.Value.Batch != nil && match {
				if knownValue, ok := buffer[value.Value.Siret]; ok {
					newValue, _ := (*knownValue).merge(*value)
					buffer[value.Value.Siret] = &newValue
				} else {
					value.ID = bson.NewObjectId()
					buffer[value.Value.Siret] = value
					i++
				}
			}
		}
	}(source)

	return source
}

// ServerJSFunc Function à injecter dans l'instance MongoDB
type ServerJSFunc struct {
	ID    string          `json:"id" bson:"_id"`
	Value bson.JavaScript `json:"value" bson:"value"`
}

// Add Méthode pour upsérer une fonction serveur
func (f ServerJSFunc) Add(db *mgo.Database) error {
	_, err := db.C("system.js").Upsert(bson.M{"_id": f.ID}, f)
	return err
}

// Drop Méthode pour supprimer une fonction serveur
func (f ServerJSFunc) Drop(db *mgo.Database) error {
	return db.C("system.js").Remove(bson.M{"_id": f.ID})
}

func declareDatabaseCopy(db *mgo.Database, from string, to string) {
	f := ServerJSFunc{
		ID:    "copyDatabase",
		Value: bson.JavaScript{Code: "function () {db.copyDatabase('" + from + "', '" + to + "')}"},
	}
	f.Add(db)
}

func removeDatabaseCopy(db *mgo.Database) {
	f := ServerJSFunc{
		ID: "copyDatabase",
	}
	f.Drop(db)
}

func declareServerFunctions(db *mgo.Database) error {

	f := ServerJSFunc{
		ID:    "generatePeriodSerie",
		Value: bson.JavaScript{Code: "function (date_debut, date_fin) {var date_next = new Date(date_debut.getTime());var serie = [];while (date_next.getTime() < date_fin.getTime()) {serie.push(new Date(date_next.getTime()));date_next.setUTCMonth(date_next.getUTCMonth() + 1);}return serie;}"},
	}
	err := f.Add(db)
	if err != nil {
		return err
	}

	f = ServerJSFunc{
		ID:    "compareDebit",
		Value: bson.JavaScript{Code: `function(a,b) {if (a.numero_historique < b.numero_historique) return -1;if (a.numero_historique > b.numero_historique) return 1;return 0;}`},
	}
	err = f.Add(db)
	if err != nil {
		return err
	}

	f = ServerJSFunc{
		ID:    "isRJLJ",
		Value: bson.JavaScript{Code: `function(code) {codes = ['PCL010501','PCL010502','PCL030105','PCL05010102','PCL05010203','PCL05010402','PCL05010302','PCL05010502','PCL05010702','PCL05010802','PCL05010901','PCL05011003','PCL05011101','PCL05011203','PCL05011303','PCL05011403','PCL05011503','PCL05011603','PCL05011902','PCL05012003','PCL0108','PCL0109','PCL030107','PCL030108','PCL030307','PCL030308','PCL05010103','PCL05010104','PCL05010204','PCL05010205','PCL05010303','PCL05010304','PCL05010403','PCL05010404','PCL05010503','PCL05010504','PCL05010703','PCL05010803','PCL05011004','PCL05011005','PCL05011102','PCL05011103','PCL05011204','PCL05011205','PCL05011304','PCL05011305','PCL05011404','PCL05011405','PCL05011504','PCL05011505','PCL05011604','PCL05011605','PCL05011903','PCL05011904','PCL05012004','PCL05012005','PCL040802'];return codes.includes(code);}`},
	}
	err = f.Add(db)
	if err != nil {
		return err
	}

	altaresCodes := `function(code) {var codeLiquidation = ['PCL0108', 'PCL010801','PCL010802','PCL030107','PCL030307','PCL030311','PCL05010103','PCL05010204','PCL05010303','PCL05010403','PCL05010503','PCL05010703','PCL05011004','PCL05011102','PCL05011204','PCL05011206','PCL05011304','PCL05011404','PCL05011504','PCL05011604','PCL05011903','PCL05012004','PCL050204','PCL0109','PCL010901','PCL030108','PCL030308','PCL05010104','PCL05010205','PCL05010304','PCL05010404','PCL05010504','PCL05010803','PCL05011005','PCL05011103','PCL05011205','PCL05011207','PCL05011305','PCL05011405','PCL05011505','PCL05011904','PCL05011605','PCL05012005'];
		var codePlanSauvegarde = ['PCL010601','PCL0106','PCL010602','PCL030103','PCL030303','PCL03030301','PCL05010101','PCL05010202','PCL05010301','PCL05010401','PCL05010501','PCL05010506','PCL05010701','PCL05010705','PCL05010801','PCL05010805','PCL05011002','PCL05011202','PCL05011302','PCL05011402','PCL05011502','PCL05011602','PCL05011901','PCL0114','PCL030110','PCL030310'];
		var codeRedressement = ['PCL0105','PCL010501','PCL010502','PCL010503','PCL030105','PCL030305','PCL05010102','PCL05010203','PCL05010302','PCL05010402','PCL05010502','PCL05010702','PCL05010706','PCL05010802','PCL05010806','PCL05010901','PCL05011003','PCL05011101','PCL05011203','PCL05011303','PCL05011403','PCL05011503','PCL05011603','PCL05011902','PCL05012003'];
		var codeInBonis = ['PCL05','PCL0501','PCL050101','PCL050102','PCL050103','PCL050104','PCL050105','PCL050106','PCL050107','PCL050108','PCL050109','PCL050110','PCL050111','PCL050112','PCL050113','PCL050114','PCL050115','PCL050116','PCL050119','PCL050120','PCL050121','PCL0503','PCL050301','PCL050302','PCL0508','PCL010504','PCL010803','PCL010902','PCL050901','PCL050902','PCL050903','PCL050904','PCL0504','PCL050303','PCL050401','PCL050402','PCL050403','PCL050404','PCL050405','PCL050406'];
		var codeContinuation = ['PCL0202'];
		var codeSauvegarde = ['PCL0203','PCL020301','PCL0205','PCL040408'];
		var codeCession = ['PCL0204','PCL020401','PCL020402','PCL020403'];
		var res = null;
		if (codeLiquidation.includes(code)) 
			res = 'liquidation';
		else if (codePlanSauvegarde.includes(code))
			res = 'plan_sauvegarde';
		else if (codeRedressement.includes(code))
			res = 'plan_redressement';
		else if (codeInBonis.includes(code))
			res = 'in_bonis';
		else if (codeContinuation.includes(code))
			res = 'continuation';
		else if (codeSauvegarde.includes(code))
			res = 'sauvegarde';
		else if (codeCession.includes(code))
			res = 'cession';

		return res;
	}	`

	f = ServerJSFunc{
		ID:    "altaresToHuman",
		Value: bson.JavaScript{Code: altaresCodes},
	}
	err = f.Add(db)
	if err != nil {
		return err
	}

	f = ServerJSFunc{
		ID:    "DateAddMonth",
		Value: bson.JavaScript{Code: `function(date, nbMonth) {var result = new Date(date.getTime());result.setUTCMonth(result.getUTCMonth() + nbMonth);return result;}`},
	}
	err = f.Add(db)
	if err != nil {
		return err
	}

	return nil
}

// Status statut de la base de données
type Status struct {
	ID     AdminID       `json:"id" bson:"_id"`
	Status *string       `json:"status" bson:"status"`
	Epoch  int           `json:"epoch" bson:"epoch"`
	DB     *mgo.Database `json:"-" bson:"-"`
}

func (status *Status) read() error {
	var tempStatus Status
	err := status.DB.C("Admin").Find(bson.M{"_id.key": "status", "_id.type": "status"}).One(&tempStatus)
	status.ID = tempStatus.ID
	status.Status = tempStatus.Status
	return err
}

func (status *Status) write() error {
	status.Epoch++
	_, err := status.DB.C("Admin").Upsert(bson.M{"_id": bson.M{"key": "status", "type": "status"}}, status)
	return err
}

func getDBStatus(c *gin.Context) {
	c.JSON(200, db.Status.Status)
}

func (status *Status) setDBStatus(message *string) error {
	if status.Status != nil && message != nil {
		return errors.New("Ne peut remplacer une activité en cours")
	}
	status.Status = message
	return status.write()
}

func epoch(c *gin.Context) {
	c.JSON(200, db.Status.Epoch)
}
