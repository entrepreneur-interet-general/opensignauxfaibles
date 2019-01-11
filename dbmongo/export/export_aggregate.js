db = db.getSiblingDB(db_name)
db.getCollection("Features").aggregate(
  [
    { 
      "$match" : {
        "_id.batch" : batch, 
        "_id.algo" : algo
      }
    }, 
    { 
      "$unwind" : {
        "path" : "$value"
      }
    }, 
    {
      $replaceRoot: { newRoot: "$value" } 	
    },
    {
      "$match": {
        "effectif" : {$gte: 10}
      }
    },
    { 
      "$project" : {
        "_id" : 0.0
      }
    }, 
    { 
      "$out" : "to_export"
    }
  ], 
  { 
    "allowDiskUse" : true
  }
);

