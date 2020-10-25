<template>
  <div>
    <Header />
    <Sidebar />
    <main>
      <div v-if="$store.state.dashboard.error">
        <h1>An Error Occured</h1>
        <p>{{ $store.state.dashboard.error }}</p>
      </div>
      <Nuxt v-else />
    </main>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  middleware: ['auth', 'administrators-only'],
  updated() {
    if (
      this.$store.state.dashboard.error !== null &&
      this.$store.state.dashboard.error.name === 'AuthError'
    ) {
      this.$store
        .dispatch('authentication/logout')
        .then(() => this.$router.push('/login'))
        .catch(console.error)
    }
  },
})
</script>

<style>
@import url('https://fonts.googleapis.com/css?family=Raleway');

:root {
  --primary-color: #0082c8;
  --primary-color-accent: #1d75b4;
  --secondary-color: #353435;
  --secondary-color-accent: #232528;
  --light-text-color: white;
}

body {
  font-family: Raleway, sans-serif;
  font-weight: 300;
  height: 100vh;
  overflow: hidden;
  background-color: #f2f2f2;
}

h1 {
  font-weight: 400;
}

h1,
h2,
h5,
p {
  margin: 0;
  padding: 5px;
}

.brand {
  font-size: 28px;
  font-weight: 400;
  color: inherit;
  letter-spacing: 1px;
  text-transform: uppercase;
  text-decoration: none;
  line-height: 50px;
}

main {
  height: 100%;
  margin: 50px 0 0 250px;
  padding: 5px;
  overflow-y: scroll;
}

.loading {
  margin: 10px;
}

.filter-panel {
  margin: 0 auto;
  width: 100%;
  border-radius: 10px;
  background-color: #fff;
  margin: 5px;
  padding: 5px;
}

.panel {
  margin: 0 auto;
  width: 100%;
  border-radius: 10px;
  background-color: #fff;
  margin: 5px;
  padding: 5px;
}

.panel-head {
  border-bottom: 1px solid grey;

  font-weight: 700;
}

.panel-head svg {
  vertical-align: middle;
  margin-right: 7px;
}

.panel-head h1 {
  display: inline-block;
  vertical-align: middle;
  line-height: normal;
}

.panel-body {
  padding: 5px;
}

.subtitley {
  font-size: 1em;
}
</style>
