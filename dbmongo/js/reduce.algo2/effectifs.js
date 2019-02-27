function effectifs (v, output_indexed) {

  let output_effectif = {}  

  let map_effectif = Object.keys(v.effectif).reduce(function (map_effectif, hash) {
    var effectif = v.effectif[hash]
    if (effectif == null) {
      return map_effectif
    }
    var effectifTime = effectif.periode.getTime()
    map_effectif[effectifTime] = (map_effectif[effectifTime] || 0) + effectif.effectif
    return map_effectif
  }, {})

  Object.keys(map_effectif).forEach(time =>{
    var time_d = new Date(parseInt(time))
    var time_offset = DateAddMonth(time_d, -offset_effectif -1)
    //if (time_offset.getTime() in output_indexed){
    output_effectif[time_offset.getTime()] = output_effectif[time_offset.getTime()] || {}
    output_effectif[time_offset.getTime()].effectif = map_effectif[time]
    output_effectif[time_offset.getTime()].date_effectif = time_d
    //}

    var past_month_offsets = [6,12,18,24]
    past_month_offsets.forEach(offset => {
      var time_past_offset = DateAddMonth(time_offset, offset)
      var variable_name_effectif = "effectif_past_" + offset
      // if (time_past_offset.getTime() in output_effectif){
      output_effectif[time_past_offset.getTime()] = output_effectif[time_past_offset.getTime()] || {}
      output_effectif[time_past_offset.getTime()][variable_name_effectif] = map_effectif[time]
    })
  })

  Object.keys(output_effectif).forEach(k => {
    if (output_effectif[k].effectif == null) {
      delete output_effectif[k]
      // delete output_array[index]
    }
  })
  print("-----------------------------------------")
  print(JSON.stringify(output_effectif, null, 2))
  return(output_effectif)
}
