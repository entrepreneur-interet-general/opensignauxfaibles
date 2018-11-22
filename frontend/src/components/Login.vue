<template>
<v-container fluid >
  <v-layout align-center justify-center>
    <v-flex xs4>
      <v-card class="elevation-6">
        <v-form @submit="login">
        <v-toolbar class="toolbar" id="toolbar" card color="indigo darken-1">
          <v-spacer> </v-spacer>
            Identification
          <v-spacer> </v-spacer>
        </v-toolbar>
        <v-card-text>
          <v-text-field label="Nom d'utilisateur" v-model="username" outline></v-text-field >
          <v-text-field type="password" label="Mot de passe" v-model="password" outline></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-btn block outline  color="indigo darken-4" flat type="submit">LOGIN</v-btn>
        </v-card-actions>
                  </v-form>

      </v-card>
    </v-flex>
  </v-layout>
</v-container>
</template>

<script>

export default {
  computed: {
    username: {
      get (username) {
        return this.$store.state.credentials.username
      },
      set (username) {
        this.$store.commit('setUser', username)
      }
    },
    password: {
      get (password) {
        return this.$store.state.credentials.password
      },
      set (password) {
        this.$store.commit('setPassword', password)
      }
    },
    token () {
      return this.$store.state.token
    }
  },
  methods: {
    login () {
      this.$store.commit('login')
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