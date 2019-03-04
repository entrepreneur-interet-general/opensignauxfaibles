function compte (compte, periodes) {
  let output_compte = {}

  //  var offset_compte = 3

  Object.keys(compte).forEach(hash =>{
    var one_compte = compte[hash]
    var periode = one_compte.periode.getTime()

    if (periode in periodes){
      output_compte[periode] =  {}
      output_compte[periode].compte_urssaf =  one_compte.compte
    }
  })

  return output_compte
}
