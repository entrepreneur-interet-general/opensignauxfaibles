function map() { 
    emit(
    {'siren': this.value.siren,
     'batch': actual_batch,
     'algo': 'algo2'},
    this.value)
} 