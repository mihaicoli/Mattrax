<template>
  <div class="container">
    <Breadcrumb />
    <div>
      <h1 class="title">Mattrax Dashboard</h1>
      <div class="links">
        <a
          href="https://github.com/mattrax/Mattrax"
          target="_blank"
          class="button--grey"
        >
          GitHub
        </a>
        <a
          href="https://github.com/mattrax/Mattrax/issues/new"
          target="_blank"
          class="button--grey"
        >
          Report Bug
        </a>
        <a
          href="mailto:oscar@otbeaumont.me?subject=Mattrax Support:"
          target="_blank"
          class="button--grey"
        >
          Support
        </a>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'

export default Vue.extend({
  layout: 'dashboard',
  middleware: ['auth'],
  meta: {
    breadcrumb: [{ name: 'Dashboard' }],
  },
})
</script>

<style>
.container {
  margin: 0 auto;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  text-align: center;
}

.title {
  font-family: 'Quicksand', 'Source Sans Pro', -apple-system, BlinkMacSystemFont,
    'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  display: block;
  font-weight: 300;
  font-size: 100px;
  color: #35495e;
  letter-spacing: 1px;
}

.subtitle {
  font-weight: 300;
  font-size: 42px;
  color: #526488;
  word-spacing: 5px;
  padding-bottom: 15px;
}

.links {
  padding-top: 15px;
}

.button--green {
  display: inline-block;
  border-radius: 4px;
  border: 1px solid #3b8070;
  color: #3b8070;
  text-decoration: none;
  padding: 10px 30px;
}

.button--green:hover {
  color: #fff;
  background-color: #3b8070;
}

.button--grey {
  display: inline-block;
  border-radius: 4px;
  border: 1px solid #35495e;
  color: #35495e;
  text-decoration: none;
  padding: 10px 30px;
  margin-left: 15px;
}

.button--grey:hover {
  color: #fff;
  background-color: #35495e;
}
</style>
