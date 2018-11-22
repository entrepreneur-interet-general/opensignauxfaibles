<template>
  <v-app id="app">
    <v-toolbar
    class="toolbar"
    dense
    color="grey lighten-3"
    app
    >

    <v-menu v-if="login" open-on-click top offset-y transition="slide-x-transition">
      <v-btn flat slot="activator">
        <v-icon large color="grey darken-4">menu</v-icon>
      </v-btn>
      
      <v-list>
        <v-list-tile
        v-for="(item, index) in items"
        :key="index"
        @click="setMenu(index)"
        :to="item.action"
        style="vertical-align: middle">
          <v-list-tile-action>
            <v-icon :color="item.color">{{ item.icon }}</v-icon>
          </v-list-tile-action>
          <v-list-tile-title
          :prepend-icon="item.icon"
          style="vertical-align: middle">
            {{ item.title }}
          </v-list-tile-title>
        </v-list-tile>
      </v-list>
    </v-menu>
    <v-toolbar-title>

      <span class="span" id="blue">Signaux</span>
      <span class="span" id="red">Faibles</span>
      <span class="span gray">– Plateforme de détection des entreprises fragilisées</span>

    </v-toolbar-title>
    <v-spacer></v-spacer>
    
    <v-btn flat icon disabled>
      <v-tooltip bottom open-delay="1000" transition="fade-transition" v-if="login && dbstatus != null">
        {{ dbstatus }}
          <v-progress-circular slot="activator"
            indeterminate
            color="grey darken-4"
          ></v-progress-circular>
      </v-tooltip>
    </v-btn>

    <Uploads/>
      <v-dialog 
      lazy 
      v-model="dialog" 
      scrollable v-if="login">
        <v-btn flat lazy icon slot="activator" >
          <v-icon small color="grey darken-4">fa-bell</v-icon>
        </v-btn>
        <v-card
        style="min-height: 90vh;">
          <v-toolbar
          style="vertical-align: middle"
          class="headline"
          color="grey darken-4"
          dark
          dense
          card>
            <v-toolbar-title>
              Activité Serveur
            </v-toolbar-title>
            <v-spacer></v-spacer>
            <v-toolbar-items >
              <v-switch
              height="100%"
              v-for="(key,index) in avatar"
              :key="key.color"
              v-model="activityFilter" 
              :color="key.color"
              :prependIcon="key.icon" 
              :value="index"/>
              <v-spacer></v-spacer>
              <v-icon color="red" @click="dialog = false">fa-window-close</v-icon>
              </v-toolbar-items>
            </v-toolbar>
          <v-card-text >
            <v-list dense>
              <template 
              v-for="(m, index) in messages.filter(m => activityFilter.indexOf(m.priority) > -1)">
                <v-divider :key="'divide' + index"/>
                <v-list-tile
                  :key="index"
                  avatar
                >
                  <v-list-tile-avatar>
                    <v-icon :color="(avatar[m.priority]||{}).color">{{ (avatar[m.priority]||{}).icon }}</v-icon>
                  </v-list-tile-avatar> 

                  <v-list-tile-content>
                    <v-list-tile-title v-html="m.event"/>
                    <v-list-tile-sub-title class="caption" v-html="m.date.toLocaleString()"/>
                  </v-list-tile-content>
                </v-list-tile>
              </template>
            </v-list>
          </v-card-text>
        </v-card>
      </v-dialog>
      
      <v-tooltip bottom open-delay="1000" transition="fade-transition" v-if="login" >
      Déconnexion
          <v-btn icon @click="logout()" slot="activator">
            <v-icon small color="grey darken-4">fa-sign-out-alt</v-icon>
          </v-btn>
      </v-tooltip>

    </v-toolbar>
    <v-content>
        <Login v-if="!login" :state="login"/>
        <router-view v-if="login" />
    </v-content>

    <v-footer class="elevation-12">
      <v-btn 
        flat 
        icon 
        color="blue"
        href="https://github.com/entrepreneur-interet-general/opensignauxfaibles">
        <v-icon>fab fa-github</v-icon>
      </v-btn>
    </v-footer>


  </v-app>
</template>

<script>
import Login from '@/components/Login'
import Uploads from '@/components/Uploads'

export default {
  methods: {
    setMenu (key) {
      this.menu.title = this.items[key].title
      this.menu.color = this.items[key].color
    },
    logout () {
      this.$store.commit('logout')
    }
  },
  components: { Login, Uploads },
  computed: {
    uploads () {
      return this.$store.getters.getUploads
    },
    login () {
      return this.$store.state.token != null
    },
    batches () {
      return JSON.stringify(this.$store.state.batches, null, 2)
    },
    dbstatus () {
      return this.$store.state.dbstatus
    },
    messages () {
      return this.$store.getters.messages
    }
  },
  data () {
    return {
      fixed: false,
      drawer: false,
      dialog: false,
      activityFilter: ['info', 'warning', 'critical'],
      menu: {
        title: 'Accueil',
        color: 'light-green darken-4'
      },
      items: [
        { title: 'Accueil',
          action: '/',
          icon: 'fa-home',
          color: 'green darken-3' },
        { title: 'Détection',
          action: '/Browse',
          icon: 'fa-search',
          color: 'indigo darken-4' },
        { title: 'Données',
          action: '/data',
          icon: 'fa-database',
          color: 'red darken-4' }
      ],
      avatar: {
        debug: {icon: 'fa-cogs', color: 'blue'},
        info: {icon: 'fa-info', color: 'green'},
        warning: {icon: 'fa-exclamation-triangle', color: 'yellow'},
        critical: {icon: 'fa-sad-cry', color: 'red'}
      }
    }
  },
  mounted () {
    this.items.forEach((item, index) => {
      if (item.action === location.hash.substring(1)) {
        this.setMenu(index)
      }
    })
  },
  name: 'App'
}
</script>

<style scoped>
 @import url('https://fonts.googleapis.com/css?family=Quicksand');
body {
  font-family: 'Quicksand', sans-serif;
}
.gray {
  color: #606060;
  font-family: 'Quicksand', sans-serif;
}
#blue {
  color: #20449a;
  font-family: 'Quicksand', sans-serif;
}
#red {
  color: #e9222e;
  font-family: 'Quicksand', sans-serif;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
v-list {
  font-family: 'Quicksand', sans-serif;
}
</style>

<style>
.toolbar {
  background-image: url("/static/bgtoolbar.png");
  color: "black";
}
#app {
  min-height: 5000px;
  background-image: url("/static/bgapp.png");
}
.span {
  max-height: 10px
}
</style>