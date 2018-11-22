<template>
  <div>
    <v-navigation-drawer
    class="elevation-6"
    absolute
    permanent
    :mini-variant = "mini"
    style="z-index: 1"
    >
    <v-list two-line>
      <v-list-group
      v-model="nomini">
        <v-list-tile  slot="activator" @click="mini=!mini">
          <v-list-tile-action>
            <v-icon>fa-filter</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>Filtrage</v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-icon>fa-industry</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-select
              :items="naf1"
              v-model="naf"
              label="Secteur d'activité"
            ></v-select>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
            <v-icon>fa-users</v-icon>
          </v-list-tile-action>
          <v-list-tile-content>
            <v-select
              :items="effectifClass"
              v-model="minEffectif"
              label="Effectif minimum"
            ></v-select>
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
            v-model="entrepriseConnue">   
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            Entreprise non suivie
          </v-list-tile-content>

        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
            v-model="horsCCSF">   
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            hors CCSF
          </v-list-tile-content>
        </v-list-tile>
        <v-list-tile>
          <v-list-tile-action>
          <v-checkbox
            v-model="horsProcol">   
          </v-checkbox>
          </v-list-tile-action>
          <v-list-tile-content>
            hors Procédure Collective
          </v-list-tile-content>
        </v-list-tile>
        </v-list-group>
    </v-list>
    </v-navigation-drawer>

    <v-card :class="mini?'widget_large':'widget_tiny'">
    <v-data-table
    v-model="selected"
    :headers="headers"
    :items="predictionFiltered"
    :pagination.sync="pagination"
    select-all
    item-key="name"
    class="elevation-1"
    loading="false"
    :rows-per-page-items="[10]"
    >
      <template slot="headers" slot-scope="props">
        <tr>
          <th/>
          <th
            v-for="header in props.headers"
            :key="header.text"
            :class="['column sortable', pagination.descending ? 'desc' : 'asc', header.value === pagination.sortBy ? 'active' : '']"
            @click="changeSort(header.value)"
          >
          <v-icon small>arrow_upward</v-icon>
          {{ header.text }}
          </th>
        </tr>
      </template>
      <template slot="items" slot-scope="props">
        <tr 
        :active="props.selected"
        >
          <td>
            <v-icon
            @click.left="open(props.item, true)"
            @click.middle="open(props.item, false)"
            >
              fa-address-card
            </v-icon>
          </td>
          <td ><v-tooltip left>
            <div slot="activator">{{ props.item.raison_sociale }}</div>
            {{ props.item._id.siret }}
            </v-tooltip> </td>
          <td class="text-xs-center"><widgetPrediction :prob="props.item.prob" :diff="props.item.diff"/></td>
          <td class="text-xs-right">
            {{ props.item.effectif }}
          </td>
          <td class="text-xs-right">{{ props.item.default_urssaf?"oui":"non" }}</td>
          <td>
            <IEcharts
              style="height: 40px"
              :loading="chart"
              :option="getMargeOption(
                (props.item.bdf || []).map(b => {
                return {'x': b.annee, 'y': b.taux_marge}
                })
              )"
            />
          </td>
          <td>
            <IEcharts
              style="height: 40px"
              :loading="chart"
              :option="getMargeOption(
                (props.item.bdf || []).map(b => {
                return {'x': b.annee, 'y': b.poids_frng}
                })
              )"
            />
          </td>
          <td>
            <IEcharts
              style="height: 40px"
              :loading="chart"
              :option="getMargeOption(
                (props.item.bdf || []).map(b => {
                return {'x': b.annee, 'y': b.financier_court_terme}
                })
              )"
            />
          </td>
        </tr>
      </template>
    </v-data-table>

    </v-card>
  </div>
</template>

<script>
  import IEcharts from 'vue-echarts-v3/src/lite.js'
  import 'echarts/lib/chart/line'
  import 'echarts/lib/component/title'
  import widgetPrediction from '@/components/widgetPrediction'
  export default {
    props: ['batchKey'],
    components: {
      IEcharts,
      widgetPrediction
    },
    data: () => ({
      effectifClass: [10, 20, 50, 100],
      selected: [],
      mini: true,
      loading: false,
      chart: false,
      pagination: {
        sortBy: 'name'
      },
      naf1: [
        'Tous',
        'Activités spécialisées, scientifiques et techniques',
        'Activités de services administratifs et de soutien',
        'Industrie manufacturière',
        'Hébergement et restauration',
        'Construction',
        'Transports et entreposage',
        'Commerce ; réparation d\'automobiles et de motocycles',
        'Santé humaine et action sociale',
        'Autres activités de services',
        'Arts, spectacles et activités récréatives',
        'Industries extractives',
        'Production et distribution d\'eau ; assainissement, gestion des déchets et dépollution',
        'Information et communication',
        'Activités financières et d\'assurance',
        'Activités immobilières',
        'Agriculture, sylviculture et pêche',
        'Production et distribution d\'électricité, de gaz, de vapeur et d\'air conditionné',
        'Activités extra-territoriales'
      ],
      headers: [
        {
          text: 'raison sociale',
          align: 'left',
          value: 'raison_sociale'
        },
        {text: 'détection', value: 'prob'},
        {text: 'emploi', value: 'effectif'},
        {text: 'Défault urssaf', value: 'default_urssaf'},
        {text: 'Taux de marge', value: 'taux_marge'},
        {text: 'Fond de roulement', value: 'fond_roulement'},
        {text: 'Financier court terme', value: 'financier_court_terme'}
      ],
      prediction: [],
      naf: 'Industrie manufacturière',
      minEffectif: 20,
      entrepriseConnue: true,
      horsCCSF: true,
      horsProcol: true
    }),
    computed: {
      nomini: {
        get () { return !this.mini },
        set (mini) { this.mini = !mini }
      },
      predictionFiltered () {
        return this.prediction.filter(p => this.applyFilter(p))
      },
      tabs: {
        get () { return this.$store.getters.getTabs },
        set (tabs) { this.$store.dispatch('updateTabs', tabs) }
      },
      activeTab: {
        get () { return this.$store.getters.activeTab },
        set (activeTab) { this.$store.dispatch('updateActiveTab', activeTab) }
      }
    },
    mounted () {
      this.getPrediction()
    },
    methods: {
      open (etab, focus) {
        if (this.tabs.findIndex(t => t.siret === etab._id.siret) === -1) {
          let i = this.tabs.push({
            'type': 'Etablissement',
            'param': etab.raison_sociale,
            'siret': etab._id.siret,
            'batch': '1802'
          })
          if (focus) { this.activeTab = i - 1 }
        }
      },
      applyFilter (p) {
        return (this.naf === 'Tous' || p.naf1 === this.naf) &&
        (p.effectif >= this.minEffectif) &&
        (p.connu === false || this.entrepriseConnue === false) &&
        (p.ccsf === false || this.horsCCSF === false) &&
        (p.procol === 'in_bonis' || this.horsProcol === false)
      },
      getNAF () {
        var self = this
        this.$axios.get('/api/data/naf').then(response => { self.naf = response.data })
      },
      toggleAll () {
        if (this.selected.length) this.selected = []
        else this.selected = this.desserts.slice()
      },
      changeSort (column) {
        if (this.pagination.sortBy === column) {
          this.pagination.descending = !this.pagination.descending
        } else {
          this.pagination.sortBy = column
          this.pagination.descending = false
        }
      },
      getMargeOption (marge) {
        return {
          title: {
            text: null
          },
          tooltip: {
            trigger: 'axis',
            axisPointer: {
              type: 'cross',
              label: {
                backgroundColor: '#283b56'
              }
            }
          },
          toolbox: {
            show: true
          },
          xAxis: {
            show: true,
            type: 'category',
            axisTick: false,
            data: marge.map((m) => m.x)
          },
          yAxis: {
            type: 'value',
            show: false,
            min: -150,
            max: 150
          },
          series: [{
            color: 'indigo',
            smooth: true,
            name: 'taux marge',
            type: 'line',
            data: marge.map((m) => m.y)
          }]
        }
      },
      getPrediction () {
        this.loading = true
        var self = this
        this.$axios.get('/api/data/prediction').then(response => {
          var prediction = response.data
          console.log(prediction)
          prediction.forEach(p => {
            p.bdf = Object.keys((p.bdf || {})).map(b => p.bdf[b]).sort((a, b) => a.annee < b.annee)
          })
          console.log(prediction)
          this.prediction = prediction
          self.loading = false
        })
      }
    }
  }
</script>

<style scoped>
.echarts {
  width: 400px
}
.widget_tiny {
  position: absolute;
  left: 320px;
  top: 20px; 
  right: 20px;
}
.widget_large {
  position: absolute;
  left: 100px;
  top: 20px;
  right: 20px;
}
.pointer:hover {
  cursor: hand;
}
</style>