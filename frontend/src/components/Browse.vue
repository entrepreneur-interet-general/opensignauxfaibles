<template>
<div>
  <v-tabs
      v-model="activeTab"
      color="indigo darken-4"
      dark
      slider-color="red accent-2"
      lazy
    >
      <v-tab
        v-for="(tab, index) in tabs"
        :key="index"
        ripple
      >
       {{ tab.param }}
       <div style="width: 10px"/>

        <v-icon
        v-if="tab.type === 'Etablissement'"
        color="red accent-1"
        style="font-size: 15px"
        @click="close(index)">
          fa-times
        </v-icon>
      </v-tab>
      <v-tab-item
        style="min-height: 5000vh;"
        v-for="(tab,index) in tabs"
        :key="index"
      >
        <PredictionTable v-if="tab.type==='Prediction'" :batch="tab.batch"/>
        <Etablissement v-if="tab.type==='Etablissement'" :siret="tab.siret" :batch="tab.batch"/>
      </v-tab-item>
    </v-tabs>

</div>
</template>

<script>
import PredictionTable from '@/components/PredictionTable'
import Etablissement from '@/components/Etablissement'

export default {
  data () {
    return {
      active: 0
    }
  },
  mounted () {
    this.$store.commit('updateBatches')
  },
  methods: {
    setCurrentBatchKey (batchKey) {
      this.currentBatchKey = batchKey
    },
    close (tabIndex) {
      this.activeTab = Math.min(this.activeTab, (this.tabs.length - 2))
      this.tabs = this.tabs.filter((tab, index) => index !== tabIndex)
    }
  },
  computed: {
    tabs: {
      get () { return this.$store.getters.getTabs },
      set (tabs) { this.$store.dispatch('updateTabs', tabs) }
    },
    activeTab: {
      get () { return this.$store.getters.activeTab },
      set (activeTab) { this.$store.dispatch('updateActiveTab', activeTab) }
    },
    currentBatchKey: {
      get () {
        return this.$store.state.currentBatchKey
      },
      set (value) {
        this.$store.commit('setCurrentBatchKey', value)
      }
    },
    batches () {
      return this.$store.state.batches.filter(b => b.readonly === true).map(batch => batch.id.key)
    }
  },
  components: { PredictionTable, Etablissement },
  name: 'Browse'
}
</script>

