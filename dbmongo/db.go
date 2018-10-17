package main

import (
  "log"
  "time"
  "github.com/spf13/viper"

  "github.com/gin-gonic/gin"
  "github.com/globalsign/mgo"
  "github.com/globalsign/mgo/bson"
)

// DB Initialisation de la connexion MongoDB
func DB() gin.HandlerFunc {

  dbDial := viper.GetString("DB_DIAL")
  dbDatabase := viper.GetString("DB")

  mongodb, err := mgo.Dial(dbDial)
  mongodb.SetSocketTimeout(3600 * time.Second)
  db := mongodb.DB(dbDatabase)

  // pousse les fonctions partagées JS
  declareServerFunctions(db)

  if err != nil {
    log.Panic(err)
  }

  chanEntreprise := insertEntreprise(db)
  chanEtablissement := insertEtablissement(db)

  go func() {
    for range time.Tick(30 * time.Second) {
      chanEntreprise <- &ValueEntreprise{}
      chanEtablissement <- &ValueEtablissement{}
    }
  }()

  return func(c *gin.Context) {
    c.Set("ChanEntreprise", chanEntreprise)
    c.Set("ChanEtablissement", chanEtablissement)
    c.Set("DBSESSION", mongodb)
    c.Set("DB", db)
    c.Next()
  }
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
			if value.Value.Batch != nil {
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
			if value.Value.Batch != nil {
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
func (f ServerJSFunc) Add(db *mgo.Database) {
  db.C("system.js").Upsert(bson.M{"_id": f.ID}, f)
}

// Drop Méthode pour supprimer une fonction serveur
func (f ServerJSFunc) Drop(db *mgo.Database) {
  db.C("system.js").Remove(bson.M{"_id": f.ID})
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

func declareServerFunctions(db *mgo.Database) {

  f := ServerJSFunc{
    ID:    "generatePeriodSerie",
    Value: bson.JavaScript{Code: "function (date_debut, date_fin) {var date_next = new Date(date_debut.getTime());var serie = [];while (date_next.getTime() < date_fin.getTime()) {serie.push(new Date(date_next.getTime()));date_next.setUTCMonth(date_next.getUTCMonth() + 1);}return serie;}"},
  }
  f.Add(db)
  f = ServerJSFunc{
    ID:    "compareDebit",
    Value: bson.JavaScript{Code: `function(a,b) {if (a.numero_historique < b.numero_historique) return -1;if (a.numero_historique > b.numero_historique) return 1;return 0;}`},
  }
  f.Add(db)

  f = ServerJSFunc{
    ID:    "isRJLJ",
    Value: bson.JavaScript{Code: `function(code) {codes = ['PCL010501','PCL010502','PCL030105','PCL05010102','PCL05010203','PCL05010402','PCL05010302','PCL05010502','PCL05010702','PCL05010802','PCL05010901','PCL05011003','PCL05011101','PCL05011203','PCL05011303','PCL05011403','PCL05011503','PCL05011603','PCL05011902','PCL05012003','PCL0108','PCL0109','PCL030107','PCL030108','PCL030307','PCL030308','PCL05010103','PCL05010104','PCL05010204','PCL05010205','PCL05010303','PCL05010304','PCL05010403','PCL05010404','PCL05010503','PCL05010504','PCL05010703','PCL05010803','PCL05011004','PCL05011005','PCL05011102','PCL05011103','PCL05011204','PCL05011205','PCL05011304','PCL05011305','PCL05011404','PCL05011405','PCL05011504','PCL05011505','PCL05011604','PCL05011605','PCL05011903','PCL05011904','PCL05012004','PCL05012005','PCL040802'];return codes.includes(code);}`},
  }
  f.Add(db)

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
f.Add(db)

f = ServerJSFunc{
  ID:    "DateAddMonth",
  Value: bson.JavaScript{Code: `function(date, nbMonth) {var result = new Date(date.getTime());result.setUTCMonth(result.getUTCMonth() + nbMonth);return result;}`},
}
f.Add(db)
}
