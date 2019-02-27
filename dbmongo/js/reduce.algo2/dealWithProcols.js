function deal_with_procols(data_source, altar_or_procol, output_indexed){
  var codes  =  Object.keys(data_source).reduce((events,hash) => {
    var the_event = data_source[hash]

    if (altar_or_procol == "altares")
      var etat = f.altaresToHuman(the_event.code_evenement);
    else if (altar_or_procol == "procol")
      var etat = f.procolToHuman(the_event.action_procol, the_event.stade_procol);

    if (etat != null)
      events.push({"etat": etat, "date_proc_col": new Date(the_event.date_effet)})

    return(events)
  },[]).sort(
    (a,b) => {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
  )

  codes.forEach(
    event => {
      let periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
      var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})

      time_til_last.forEach(time => {
        if (time in output_indexed) {
          output_indexed[time].etat_proc_collective = event.etat
          output_indexed[time].date_proc_collective = event.date_proc_col
          if (event.etat != "in_bonis")
            output_indexed[time].tag_failure = true 
        }
      })
    }
  )
}
