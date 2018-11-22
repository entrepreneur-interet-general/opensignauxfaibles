function map() { 
    try{
        emit(this.value.siret, this.value) 
    } catch(error) {
        print(this.value.siret)
    }
}