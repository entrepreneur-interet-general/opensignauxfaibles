import Vue from 'vue'
import Vuex from 'vuex'
import axios from 'axios'
import createPersistedState from 'vuex-persistedstate'
import VueNativeSock from 'vue-native-websocket'

Vue.use(Vuex)

const vm = new Vue()

var axiosClient = axios.create(
  {
    headers: {
      'Content-Type': 'application/json'
    },
    baseURL: 'http://localhost:3000'
  }
)

axiosClient.interceptors.request.use(
  config => {
    if (store.state.token != null) config.headers['Authorization'] = 'Bearer ' + store.state.token
    return config
  }
)

const store = new Vuex.Store({
  plugins: [createPersistedState({storage: window.sessionStorage})],
  state: {
    credentials: {
      username: null,
      password: null
    },
    token: null,
    types: null,
    features: null,
    files: [],
    batches: [],
    dbstatus: null,
    currentBatchKey: 0,
    currentType: null,
    epoch: 0,
    socket: {
      isConnected: false,
      message: [],
      reconnectError: false
    },
    uploads: [],
    activeTab: null,
    tabs: [{
      'type': 'Prediction',
      'param': '09/2018',
      'batch': '1802'
    }
    ]
  },
  mutations: {
    updateActiveTab (state, activeTab) {
      state.activeTab = activeTab
    },
    SOCKET_ONOPEN (state, event) {
      Vue.prototype.$socket = event.currentTarget
      state.socket.isConnected = true
    },
    SOCKET_ONCLOSE (state, event) {
      state.socket.isConnected = false
    },
    SOCKET_ONERROR (state, event) {
      console.error(state, event)
    },
    // default handler called for all methods
    SOCKET_ONMESSAGE (state, message) {
      if ('journalEvent' in message) {
        let m = message.journalEvent
        state.socket.message.unshift(m)
        if (state.socket.message.length > 250) {
          state.socket.message.pop()
        }
      }
      if ('batches' in message) {
        state.batches = message.batches.reverse()
      }
      if ('files' in message) {
        state.files = message.files
      }
    },
    // mutations for reconnect methods
    SOCKET_RECONNECT (state, count) {
      console.info(state, count)
    },
    SOCKET_RECONNECT_ERROR (state) {
      wsConnect(state)
    },
    login (state) {
      axiosClient.post('/login', state.credentials).then(response => {
        state.token = response.data.token
        wsConnect(state)
        store.commit('updateRefs')
        store.commit('updateBatches')
        store.commit('updateDbStatus')
        store.commit('updateLogs')
      })
    },
    refreshToken (state) {
      axiosClient.get('/api/refreshToken').then(response => {
        state.token = response.data.token
      })
    },
    logout (state) {
      // vm.$disconnect()
      state.credentials.username = null
      state.credentials.password = null
      state.token = null
      state.types = null
      state.features = null
      state.files = null
      state.batches = []
      state.epoch = 0
      state.socket.message = []
      state.uploads = []
    },
    setUser (state, username) {
      state.credentials.username = username
    },
    setPassword (state, password) {
      state.credentials.password = password
    },
    updateBatches (state) {
      axiosClient.get('/api/admin/batch').then(response => {
        state.batches = response.data
      })
    },
    updateDbStatus (state) {
      axiosClient.get('/api/admin/status').then(response => {
        state.dbstatus = response.data
      })
    },
    setEpoch (state, epoch) {
      state.epoch = epoch
    },
    updateRefs (state) {
      axiosClient.get('/api/admin/types').then(response => { state.types = response.data.sort((a, b) => a.text.localeCompare(b.text)) })
      axiosClient.get('/api/admin/features').then(response => { state.features = response.data })
      axiosClient.get('/api/admin/files').then(response => { state.files = response.data.sort((a, b) => a.name.localeCompare(b.name)) })
    },
    updateLogs (state) {
      axiosClient.get('/api/admin/getLogs').then(response => { state.socket.message = (response.data || []) })
    },
    setCurrentBatchKey (state, key) {
      state.currentBatchKey = key
    },
    setCurrentType (state, type) {
      state.currentType = type
    },
    updateUploads (state, status) {
      let index = state.uploads.findIndex(s => s.name === status.name)
      state.uploads[index].amount = status.amount
    },
    addUpload (state, status) {
      if (state.uploads.findIndex(s => s.name === status.name) === -1) state.uploads.push(status)
    },
    resetUploads (state) {
      state.uploads = state.uploads.filter(u => u.amount < 100)
    },
    updateTabs (state, tabs) {
      state.tabs = tabs
    }
  },
  actions: {
    updateActiveTab (context, activeTab) {
      context.commit('updateActiveTab', activeTab)
    },
    updateTabs (state, tabs) {
      state.commit('updateTabs', tabs)
    },
    saveBatch (state, batch) {
      axiosClient.post('/api/admin/batch', batch).then(r => { state.currentBatch = batch })
    },
    resetUploads (context) {
      context.commit('resetUploads')
    },
    addFile (context, file) {
      axiosClient.post('/api/admin/batch/addFile', file)
    },
    upload (context, file) {
      let formData = new FormData()
      let filename = '/' + file.batch + '/' + file.type + '/' + file.file.name

      formData.append('file', file.file)
      formData.append('batch', file.batch)
      formData.append('type', file.type)

      var status = {
        'amount': 0,
        'name': filename
      }

      context.commit('addUpload', status)
      axiosClient.post(
        '/api/admin/files',
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data'
          },
          onUploadProgress: function (progressEvent) {
            var newStatus = {
              'amount': parseInt(Math.round((progressEvent.loaded * 100) / progressEvent.total)),
              'name': filename
            }
            context.commit('updateUploads', newStatus)
            context.commit('updateBatches')
          }
        }
      ).then(response => {
        let postData = {
          filename: filename,
          type: file.type,
          batch: file.batch
        }
        context.dispatch('addFile', postData)
      }).catch(function (response) {
      })
    },
    checkEpoch (context) {
      if (context.state.token != null) {
        axiosClient.get('/api/admin/epoch').then(response => {
          if (response.data !== store.state.epoch) {
            context.commit('setEpoch', response.data)
            context.commit('updateRefs')
            context.commit('updateBatches')
            context.commit('updateDbStatus')
          }
        }).catch(error => {
          if (error.response.status === 401) {
            context.commit('logout')
          }
        })
      }
    }
  },
  getters: {
    axiosConfig (state) {
      return {headers: {Authorization: 'Bearer ' + state.token}}
    },
    messages (state) {
      return state.socket.message.map(m => {
        m.date = new Date(m.date)
        return m
      })
    },
    getUploads (state) {
      return state.uploads
    },
    getTabs (state) {
      return state.tabs
    },
    activeTab (state) {
      return state.activeTab
    }
  }
})

function wsConnect (state) {
  let index = Vue._installedPlugins.indexOf(VueNativeSock)
  if (index > -1) {
    Vue._installedPlugins.splice(index, 1)
  }
  Vue.use(VueNativeSock, 'ws://localhost:3000/ws/' + state.token, {
    store: store,
    format: 'json',
    connectManually: true,
    reconnection: true,
    reconnectionAttempts: -1,
    reconnectionDelay: 3000
  })
  vm.$connect()
}

if (store.state.token != null) {
  wsConnect(store.state)
  store.commit('refreshToken')
  store.commit('updateRefs')
  store.commit('updateDbStatus')
  store.commit('updateBatches')
  store.commit('updateLogs')
}

setInterval(
    function () {
      if (store.state.token != null) {
        store.commit('refreshToken')
      }
    },
    180000)

// setInterval(
//       function () {
//         if (store.state.token != null && !store.state.socket.isConnected === false) {
//           store.dispatch('checkEpoch')
//         }
//       },
//       1000)

export default store
