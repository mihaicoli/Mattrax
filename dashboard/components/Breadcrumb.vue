<template>
  <div class="breadcrumb">
    <span v-for="(breadcrumb, idx) in breadcrumbList" :key="idx"
      >{{ breadcrumb.name }} >
    </span>
  </div>
</template>

<script lang="ts">
export default {
  data() {
    return {
      breadcrumbList: [],
    }
  },
  mounted() {
    // this.updateList()
  },
  // watch: {
  //   $route() {
  //     this.updateList()
  //   },
  // },
  methods: {
    updateList() {
      // console.log(this.$route.meta)
      // this.breadcrumbList = this.$route.meta.breadcrumb
    },
  },
}
</script>

<style>
.breadcrumb {
  font-size: 0.7em;
}

ul.breadcrumb {
  margin: 0;
  padding: 5px 10px;
  list-style: none;
  /* background-color: #eee; */
}

ul.breadcrumb li {
  display: inline;
  font-size: 0.9em;
}

ul.breadcrumb li + li:before {
  /* padding: 2px; */
  color: black;
  content: '/\00a0';
}

ul.breadcrumb li a {
  color: var(--primary-color);
  text-decoration: none;
}

ul.breadcrumb li a:hover {
  color: var(--secondary-color-accent);
  text-decoration: underline;
}
</style>
