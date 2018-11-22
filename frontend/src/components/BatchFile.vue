<template>
  <v-card class="elevation-6" >
      <v-toolbar
        class="elevation-2 toolbar"
        color="red darken-4"
        dark
        card
      >
      <v-toolbar-title class="headline">{{ currentType.text }}</v-toolbar-title>
    </v-toolbar>

      <v-tabs
      v-model="active"
      light
      slider-color="indigo darken-4"
      >
        <v-tab ripple>
          Fichiers
        </v-tab>
        <v-tab ripple>
          Téléverser
        </v-tab>
        <v-tab ripple>
          Relier    
        </v-tab>
        <v-tab-item>

          <v-list style="width: 100%">
            <v-list-tile style="width: 100%">
              <v-list-tile-content>
                <span class="strong">{{ currentFiles.length }} fichier{{ (currentFiles.length>0?"s":"") }}</span>
              </v-list-tile-content>
              <v-list-tile-action>
                <span class="light">{{ formatBytes(currentFiles.reduce((m,f) => m+f.size, 0)) }}</span>
              </v-list-tile-action> 
            </v-list-tile>

            <v-divider></v-divider>
            <v-list-tile style="width: 100%"
            v-for="(file, index) in currentFiles"
            :key="index"
            ripple
            @click="removeMark(file.name)"
            :class="removeFiles.includes(file.name)?'todelete':'tokeep'"
            >
              <v-list-tile-content>
                <span class="strong">{{ file.filename }}</span>
                <span class="light">{{ file.pathname }}</span>
              </v-list-tile-content>
              <v-list-tile-action>
                <span class="light">{{ file.psize }}</span>
                <span class="light">{{ file.pdate }}</span>
              </v-list-tile-action>
            </v-list-tile>
          </v-list>
          <v-btn 
          block 
          flat 
          v-if="removeFiles.length > 0"
          @click="remove()"
          >
          <v-icon >fa-trash</v-icon>
          </v-btn>
        
        </v-tab-item>
        <v-tab-item>
          <v-container fluid grid-list-md>
            <v-layout wrap>
              <v-flex xs12>
                <v-container fluid grid-list-md>
                  <v-data-iterator
                    :items="filesUploadArray"
                    :rows-per-page-items="rowsPerPageItems"
                    :pagination.sync="pagination"
                    content-tag="v-layout"
                    no-data-text="Ajouter des fichiers"
                    row
                    wrap
                  >
                    <v-flex
                      slot="item"
                      slot-scope="props"
                      xs4
                    >
                      <v-card class="elevation-3">
                        <v-card-title><h4>{{ props.item.file.name }}</h4></v-card-title>
                        <v-divider></v-divider>
                        <v-list dense>
                          <v-list-tile>
                            <v-list-tile-content>type:</v-list-tile-content>
                            <v-list-tile-content class="align-end">{{ props.item.type }}</v-list-tile-content>
                          </v-list-tile>
                          <v-list-tile>
                            <v-list-tile-content>taille:</v-list-tile-content>
                            <v-list-tile-content class="align-end">{{ formatBytes(props.item.file.size) }}</v-list-tile-content>
                          </v-list-tile>
                        </v-list>
                      </v-card>
                    </v-flex>
                  </v-data-iterator>
                </v-container>
              </v-flex>
              <v-flex xs12 text-xs-center>
                <v-btn
                flat
                @click="fileSelect()">
                <label 
                ref="input-file-id"
                for="input-file-id"  
                class="md-button md-raised md-primary">
                  <v-icon>fa-plus</v-icon>
                </label>
                </v-btn>
                <input
                style="display: none"
                id="input-file-id"
                ref="file"
                multiple
                v-on:change="handleFileUpload()"
                type="file"/>
                <v-btn flat :disabled="filesUploadArray.length == 0" color="success" v-on:click="submitFile()"><v-icon>fa-upload</v-icon></v-btn>
              </v-flex>
            </v-layout>
          </v-container>
        </v-tab-item>
        <v-tab-item
        class="pa-3 ma-3">
          <v-flex xs11>
            <v-combobox 
            multiple
            deletableChips
            chips
            flat
            v-model="addFiles"
            :items="files"
            item-text="filename"
            item-value="name"
            label="Relier un ancien fichier"
            ripple
            append-outer-icon="fa-plus-square"
            @click:append-outer="add"
            >
              <template
              slot="item"
              slot-scope="{index, item, parent}"
              >
                  <v-list-tile style="width: 100%">
                    <v-list-tile-content>
                      <span class="strong">{{ item.filename }}</span>
                      <span class="light">{{ item.pathname }}</span>
                    </v-list-tile-content>
                    <v-list-tile-action>
                      <span class="light">{{ item.psize }}</span>
                      <span class="light">{{ item.pdate }}</span>
                    </v-list-tile-action>
                  </v-list-tile>
              </template>
            </v-combobox>
          </v-flex>
        </v-tab-item>
      </v-tabs>
  </v-card>
</template>

<script>
export default {
  props: ['type'],
  data () {
    return {
      active: null,
      filesUpload: {},
      filesUploadArray: [],
      uploadVisible: false,
      addFiles: [],
      removeFiles: [],
      rowsPerPageItems: [16],
      pagination: {
        rowsPerPage: 9
      }
    }
  },
  methods: {
    handleFileUpload () {
      this.filesUploadArray = this.filesUploadArray.concat(Object.keys(this.$refs.file.files).map(k => {
        return {
          file: this.$refs.file.files[k],
          type: this.currentType.type,
          batch: this.currentBatch.id.key
        }
      }))
    },
    comboFilter (item, queryText, itemText) {
      return item.name.toLowerCase().includes(queryText.toLowerCase())
    },
    submitFile () {
      this.filesUploadArray.forEach(file => {
        this.$store.dispatch('upload', file)
      })
      this.filesUploadArray = []
    },
    formatBytes (a, b) {
      if (a === 0) return '0 Bytes'
      var c = 1024
      var d = b || 2
      var e = ['Octets', 'Kio', 'Mio', 'Gio', 'Tio', 'Pio', 'Eio', 'Zio', 'Yio']
      var f = Math.floor(Math.log(a) / Math.log(c))
      return parseFloat((a / Math.pow(c, f)).toFixed(d)) + ' ' + e[f]
    },
    formatDate (d) {
      return d.getFullYear() + '/' + d.getMonth() + '/' + d.getDate()
    },
    add () {
      var batch = this.currentBatch
      this.addFiles.forEach(f => {
        if (!(batch.files[this.type] || []).includes(f.name)) {
          batch.files[this.type] = (batch.files[this.type] || []).concat(f.name)
        }
      })
      this.addFiles = []
      this.currentBatch = batch
    },
    fileDetail (file) {
      var candidates = (this.files || []).filter(f => f.name === file)
      if (candidates.length === 1) {
        return candidates[0]
      } else {
        var arrayFile = file.split('/')
        var lengthArrayFile = arrayFile.length
        return {
          name: file,
          pathname: arrayFile.slice(0, lengthArrayFile - 1).join('/') + '/',
          filename: arrayFile[lengthArrayFile - 1],
          psize: 'n/a',
          pdate: 'n/a',
          date: 'n/a'
        }
      }
    },
    fileSelect () {
      this.$refs['input-file-id'].click()
    },
    removeMark (file) {
      if (this.removeFiles.includes(file)) {
        this.removeFiles = this.removeFiles.filter(f => f !== file)
      } else {
        this.removeFiles.push(file)
      }
    },
    remove () {
      var batch = this.currentBatch
      batch.files[this.type] = batch.files[this.type].filter(f => !this.removeFiles.includes(f))
      this.removeFiles = []
      this.currentBatch = batch
    }
  },
  computed: {
    uploads () {
      return this.$store.getters.getUploads
    },
    currentType () {
      return this.$store.state.types.filter(t => t.type === this.type)[0]
    },
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
    currentFiles () {
      if (this.$store.state.batches !== []) {
        return (this.currentBatch.files[this.type] || []).map(f => this.fileDetail(f))
      }
    },
    files () {
      var files = (this.$store.state.files || []).map(f => {
        var arrayFile = f.name.split('/')
        var lengthArrayFile = arrayFile.length
        f = {
          name: f.name,
          size: f.size,
          pathname: arrayFile.slice(0, lengthArrayFile - 1).join('/') + '/',
          filename: arrayFile[lengthArrayFile - 1],
          psize: this.formatBytes(f.size),
          pdate: f.date.substring(0, 10),
          date: new Date(f.date)
        }
        return f
      })
      return files
    },
    filterFiles () {
      return this.files.filter(f => {
        return !(this.currentFiles.includes(f))
      })
    }
  }
}
</script>

<style scoped>
.todelete {
  background-color: #FCE4EC;
}
.tokeep {
  background-color: white;
}
.light {
  color: #888;
  font-size: 10px;
}
.string {
  color: #000;
  font-size: 11px;
}
</style>