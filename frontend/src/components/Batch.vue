<template>
    <div >
      
    <v-navigation-drawer
    class="elevation-6"
    absolute
    permanent
    style="z-index: 1"
    >
      <v-list dense class="pt-0">
          <!-- <v-list-tile >
            <v-list-tile-content class="text-xs-center">
                          <v-icon large>fa-database</v-icon>
            </v-list-tile-content>
          </v-list-tile> -->
        <v-list-group>
          <v-list-tile slot="activator" bgcolor="red">
            <v-list-tile-action>
              <v-icon>fa-cogs</v-icon>
            </v-list-tile-action>
            <v-list-tile-content class="title">
              Paramètres
            </v-list-tile-content>
          </v-list-tile>
          <v-list-tile
          v-for="param in parameters"
          :key="batchKey + param.key"
          ripple
          @click="setCurrentType(param.key)">
            <v-list-tile-content> 
              <v-list-tile-title
              :class="(param.key===currentType) ? 'selected': null"
              >{{ param.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
        <v-list-group>
          <v-list-tile slot="activator">
            <v-list-tile-action>
              <v-icon>fa-copy</v-icon>
            </v-list-tile-action> 
            <v-list-tile-title class="title">
              Fichiers
            </v-list-tile-title>
          </v-list-tile>
          <v-list-tile
          v-for="type in types"
          :key="type.text"
          @click="setCurrentType(type.type)"
          >
            <v-list-tile-content>
              <v-list-tile-title
              :class="(type.type==currentType) ? 'selected': null">
              {{ type.text }}
              </v-list-tile-title>
            </v-list-tile-content>
            <v-list-tile-action>
              <v-icon 
              @click="toggleComplete(type.type)"
              >{{ currentBatch.complete_types.includes(type.type)?'mdi-square-inc':'mdi-shape-square-plus' }}</v-icon>
            </v-list-tile-action>
          </v-list-tile>
          <v-divider></v-divider>
        </v-list-group>
        <v-list-group>
          <v-list-tile slot="activator">
            <v-list-tile-action>
              <v-icon>fa-microchip</v-icon>
            </v-list-tile-action> 
            <v-list-tile-title class="title">
              Traitements
            </v-list-tile-title>
          </v-list-tile>
          <v-divider></v-divider>
          <v-list-tile
          v-for="process in processes"
          :key="batchKey + process.key"
          @click="setCurrentType(process.key)"
          >
            <v-list-tile-content>
              <v-list-tile-title
              :class="(process.key==currentType) ? 'selected': null">
              {{ process.text }}</v-list-tile-title>
            </v-list-tile-content>
          </v-list-tile>
        </v-list-group>
      </v-list>
      
    </v-navigation-drawer>

    <div class="widget">
      <BatchDate 
      class="elevation-6"
      :key="batchKey + 'batchDate'"
      :date="currentType"
      :param="parameters.filter(p => p.key === currentType)[0]"
      v-if="parameters.map(p => p.key).includes(currentType)"
      />
      <BatchFile 
      :key="batchKey + 'batchFile'"
      :type="currentType"
      v-if="types.map(t => t.type).includes(currentType)"
      />
      <BatchProcess
      :key="batchKey + 'batchProcess'"
      :process="processes.filter(p => p.key === currentType)[0]"
      v-if="processes.map(p => p.key).includes(currentType)"
      />
    </div>

  </div>
</template>

<script>
import BatchFile from '@/components/BatchFile'
import BatchDate from '@/components/BatchDate'
import BatchProcess from '@/components/BatchProcess'

export default {
  props: ['batchKey'],
  data () {
    return {
      currentType: null,
      parameters: [
        {text: 'Date de début', key: 'dateDebut', prop: 'date_debut'},
        {text: 'Date de fin', key: 'dateFin', prop: 'date_fin'},
        {text: 'Date de fin effectifs', key: 'dateFinEffectif', prop: 'date_fin_effectif'}
      ],
      processes: [
        {text: 'Suppression',
          color: 'red',
          key: 'reset',
          img: '/static/poubelle.png',
          description: 'Retour au batch précédent',
          do (self) { self.$axios.get('/api/batch/revert') }
        },
        {text: 'Purger',
          color: 'blue',
          key: 'purge',
          img: '/static/gomme.svg',
          description: 'Retour au paramétrage',
          do (self) { self.$axios.get('/api/batch/purge') }
        },
        {text: 'Calcul Prédictions',
          color: 'green',
          key: 'predict',
          img: '/static/warning.png',
          description: 'Intégration des données et calcul des prédictions.',
          do (self) { self.$axios.get('/api/batch/process') }
        },
        {text: 'Clôture',
          color: 'black',
          key: 'close',
          img: '/static/warning.png',
          description: 'Clôture du batch en cours et création du suivant',
          do (self) { self.$axios.get('/api/batch/next') }
        }
      ]
    }
  },
  computed: {
    currentBatchKey () {
      return this.$store.state.currentBatchKey
    },
    currentBatch: {
      get () {
        if (this.$store.state.batches !== []) {
          return this.$store.state.batches[this.currentBatchKey]
        } else {
          return null
        }
      },
      set (batch) {
        this.$store.dispatch('saveBatch', batch).then(r => this.$store.dispatch('checkEpoch'))
      }
    },
    features () {
      return this.$store.state.features
    },
    types () {
      return this.$store.state.types
    }
  },
  methods: {
    setCurrentType (type) {
      this.currentType = type
    },
    toggleComplete (type) {
      let batch = this.currentBatch
      if (batch.complete_types.includes(type)) {
        batch.complete_types = batch.complete_types.filter(t => t !== type)
      } else {
        batch.complete_types = (batch.complete_types || []).concat(type)
      }
      this.currentBatch = batch
    }
  },
  components: { BatchFile, BatchDate, BatchProcess }
}
</script>

<style>
.selected {
  color: #700;
  font-size: 15px;
}
.widget {
  position: absolute;
  left: 320px;
  top: 20px; 
  right: 20px;
}
</style>