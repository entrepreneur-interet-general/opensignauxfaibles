function cotisationsdettes(v, output_array, output_indexed) {
  var value_cotisation = {}

  var offset_cotisation = 1

  Object.keys(v.cotisation).forEach(function (h) {
    var cotisation = v.cotisation[h]
    var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.forEach(date_cotisation => {
      let date_offset = DateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
  })


  var value_dette = {}

  Object.keys(v.debit).forEach(function (h) {
    var debit = v.debit[h]

    var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})
    let date_limite = date_fin //new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
    date_traitement_debut = new Date(
      Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
    )

    date_traitement_fin = new Date(
      Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
    )

    let periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
    let periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)

    generatePeriodSerie(periode_debut, periode_fin).map(date => {
      let time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale, "montant_majorations": debit.montant_majorations}])
    })
  })    

  var numeros_compte = Array.from(new Set(
    Object.keys(v.cotisation).map(function (h) {
      return(v.cotisation[h].numero_compte)
    })))


  Object.keys(output_indexed).forEach(function (time) {
    output_indexed[time].numero_compte_urssaf = numeros_compte
    if (time in value_cotisation){
      output_indexed[time].cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
    }

    if (time in value_dette) {
      output_indexed[time].debit_array = value_dette[time]
    }
  })

  Object.keys(output_indexed).forEach(time => {
    let time_d = new Date(parseInt(time))
    var val = output_indexed[time]

    val.montant_dette = (val.debit_array || []).reduce(function (m, dette) {
      m.part_ouvriere += dette.part_ouvriere
      m.part_patronale += dette.part_patronale
      m.montant_majorations += dette.montant_majorations
      return m
    }, {"part_ouvriere": 0, "part_patronale": 0, "montant_majorations": 0})

    val.montant_part_ouvriere = val.montant_dette.part_ouvriere
    val.montant_part_patronale = val.montant_dette.part_patronale
    val.montant_majorations = val.montant_dette.montant_majorations

    let past_month_offsets = [1,2,3,6,12]
    past_month_offsets.forEach(offset => {
      let time_offset = DateAddMonth(time_d, offset)      
      let variable_name_part_ouvriere = "montant_part_ouvriere_past_" + offset
      let variable_name_part_patronale = "montant_part_patronale_past_" + offset
      if (time_offset.getTime() in output_indexed){
        let val_offset = output_indexed[time_offset.getTime()]
        val_offset[variable_name_part_ouvriere] = val.montant_part_ouvriere
        val_offset[variable_name_part_patronale] = val.montant_part_patronale
      }
    })

    let future_month_offsets = [0, 1, 2, 3, 4, 5]
    if (val.montant_part_ouvriere + val.montant_part_patronale > 0){
      future_month_offsets.forEach(offset => {
        let time_offset = DateAddMonth(time_d, offset)
        let val_offset = output_indexed[time_offset.getTime()]
        if (time_offset.getTime() in output_indexed){
          val_offset.interessante_urssaf = false    
        }
      })
    }
    delete val.montant_dette
    delete val.debit_array
  })
}
