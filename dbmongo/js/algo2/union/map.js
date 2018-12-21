function map() { 
  let value_array = (this.value.entreprise || this.value[this._id])
  for (var i=0; i< value_array.length; i++){
    if (value_array[i]){

      var m = {}
      if ("siret" in value_array[i]){
        m[value_array[i]["siret"]] = value_array[i] 
      } else {
        m["entreprise"] = value_array[i]
      }
      emit(
        {"siren": this.value.siren,
          "batch": actual_batch,
          "algo": "algo2",
          "periode": value_array[i].periode},
        m)
    }
  }
} 
