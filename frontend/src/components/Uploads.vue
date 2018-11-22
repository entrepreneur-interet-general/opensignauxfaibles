<template>
  <v-dialog v-if="uploads.length > 0"
    width="500">
    <v-btn
    flat
    icon=""
    slot="activator">
      <v-progress-circular 
      color="indigo"
      :width="2"
      v-model="globalUploads">
      <v-icon color="indigo darken-4" small>fa-upload</v-icon>
      </v-progress-circular>
    </v-btn>
    <v-card>
      <v-toolbar>
        <v-btn
        flat
        @click="resetUploads()">
          masquer les téléversements terminés
        </v-btn>
      </v-toolbar>
      <v-card-text 
      v-for="(upload, index) in uploads" 
      :key="index">
        {{ upload.name }} 
        <v-progress-linear :value="upload.amount"></v-progress-linear>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script>

export default {
  computed: {
    uploads () {
      return this.$store.getters.getUploads
    },
    globalUploads () {
      let amount = this.uploads.reduce((accu, upload) => {
        if (upload.amount < 100) {
          accu.total += 100
          accu.amount += upload.amount
        }
        return accu
      }, {total: 0, amount: 0})
      if (amount.total === 0) return 100
      else return parseInt(100 * amount.amount / amount.total)
    }
  },
  methods: {
    resetUploads () {
      this.$store.dispatch('resetUploads')
    }
  }
}
</script>

<style>
  #toolbar {
    font-size: 25px;
    color: white;
  }
</style>