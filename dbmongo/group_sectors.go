package main

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)


func group_sectors_aggregation(c *gin.Context) {

	var pipeline1, pipeline2 []bson.M
  // variables := [...]string{
  //   "ca",
  // }

// pipeline 1 creates the intermediate AggWorkspace collection

pipeline1 = append(pipeline1, bson.M{
  "$unwind": bson.M{
    "path" : "$value",
    "preserveNullAndEmptyArrays" : false,
  },
})


pipeline1 = append(pipeline1, bson.M{
  "$group": bson.M{
    "_id": bson.M{ "$substr": []interface{}{"$value.code_ape", 0, 3}},
    "mean_ca": bson.M{"$avg": "$value.ca"},
    "std_ca": bson.M{"$stdDevSamp": "$value.ca"},
    "siret": bson.M{"$addToSet": "$value.siret"},
  },
})


pipeline1 = append(pipeline1, bson.M{
  "$match": bson.M{
   "mean_ca": bson.M{"$ne": nil},
   "std_ca": bson.M{"$nin": []interface{}{0, nil}},
 },
})

pipeline1 = append(pipeline1, bson.M{
  "$unwind": bson.M{
    "path" : "$siret",
    "preserveNullAndEmptyArrays" : false,
 },
})

pipeline1 = append(pipeline1, bson.M{
  "$project" : bson.M{
    "_id" : "$siret",
    "mean_ca" : 1.0,
    "std_ca" : 1.0,
  },
})

pipeline1 = append(pipeline1, bson.M{
  "$out" : "AggWorkspace",
})

// Pipeline 2: l'aggrégation raccroche les données sectorielles aux sirets correspondants lorsque possible
pipeline2 = append(pipeline2, bson.M{
  "$unwind" : bson.M{
    "path" : "$value",
    "preserveNullAndEmptyArrays" : false,
  },
})

pipeline2 = append(pipeline2, bson.M{
  "$lookup" : bson.M{
    "from" : "AggWorkspace",
    "localField" : "value.siret",
    "foreignField" : "_id",
    "as" : "sectorial_data",
  },
})

pipeline2 = append(pipeline2, bson.M{
    "$unwind" : bson.M{
      "path" : "$sectorial_data",
      "preserveNullAndEmptyArrays" : false,
  },
})
pipeline2 = append(pipeline2, bson.M{
  "$addFields" : bson.M{
    "value.ca_norm_ape_3" : bson.M{
      "$ifNull" : []interface{}{
        bson.M{
          "$divide" : []interface{}{
            bson.M{
              "$subtract" : []interface{}{
                "$value.ca",
                "$sectorial_data.mean_ca",
              },
            },
            "$sectorial_data.std_ca",
          },
        },
        "$$REMOVE",
      },
    },
  },
})
pipeline2 = append(pipeline2, bson.M{
  "$project" : bson.M{
    "sectorial_data" : 0.0,
  },
})


pipeline2 = append(pipeline2, bson.M{
  "$out" : "AggWorkspace2",
})


var result []interface{}

err := db.DB.C("Features").Pipe(pipeline1).All(&result)
if err != nil {
  log(critical, "group_sectors", "Erreur lors des comparaisons par groupes sectoriels: "+err.Error())
  c.JSON(500, err)
  return
}

err2 := db.DB.C("Features").Pipe(pipeline2).All(&result)

if err2 != nil {
  log(critical, "group_sectors", "Erreur lors des comparaisons par groupes sectoriels: "+err2.Error())
  c.JSON(500, err2)
  return
}

}
