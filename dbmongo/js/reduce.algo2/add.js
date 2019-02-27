function add(obj, output){
      Object.keys(obj).forEach(periode => {
        if (periode in output){
          output[periode] = Object.assign(output[periode], obj[periode])
        }
      }) 
}
