function map() { 
    try{
        emit(this.value.siren, this.value) 
    } catch(error) {
        print(this.value.siren)
    }
}