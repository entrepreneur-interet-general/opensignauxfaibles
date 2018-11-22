// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import router from './router'
import Vuetify from 'vuetify'
import 'vuetify/dist/vuetify.min.css'
import store from './store'
import axios from 'axios'

Vue.use(Vuetify)

Vue.config.productionTip = false

// Prod
// npm run build
// cp dist/* ../dbmongo/static -r

Vue.prototype.$axios = axios.create(
  {
    headers: {
      'Content-Type': 'application/json'
    }
  }
)

Vue.prototype.$axios.interceptors.request.use(
  config => {
    config.baseURL = 'http://localhost:3000'
    if (store.state.token != null) config.headers['Authorization'] = 'Bearer ' + store.state.token
    return config
  }
)

Vue.prototype.$store = store

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  components: { App },
  template: '<App/>'
})
