<template>
  <div>
    <v-icon large v-if="detection.prob=='high'" color="red">fa-exclamation-triangle</v-icon> 
    <v-icon large v-if="detection.prob=='low'" color="amber">fa-question</v-icon>
    <v-icon large v-if="detection.prob=='none'" color="green">fa-check-circle</v-icon> 
    <v-icon small v-if="detection.diff=='up'" color="red">fa-arrow-up</v-icon> 
    <v-icon small v-if="detection.diff=='down'" color="blue">fa-arrow-down</v-icon>
  </div>
</template>

<script>
export default {
  props: ['prob', 'diff'],
  computed: {
    detection () {
      return this.detectionClass(this.prob, this.diff)
    }
  },
  methods: {
    detectionClass (prob, diff) {
      let probClass = ''
      let diffClass = ''

      if (prob < this.seuils.low) {
        probClass = 'none'
      } else if (prob < this.seuils.high) {
        probClass = 'low'
      } else {
        probClass = 'high'
      }

      if (Math.abs(diff) <= 0.1) {
        diffClass = 'none'
      } else if (diff < 0) {
        diffClass = 'down'
      } else {
        diffClass = 'up'
      }

      return {
        'prob': probClass,
        'diff': diffClass
      }
    }
  },
  data () {
    return {
      seuils: {
        low: 0.13,
        high: 0.37
      }
    }
  }
}
</script>