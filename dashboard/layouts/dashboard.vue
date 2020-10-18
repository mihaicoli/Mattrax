<template>
  <div>
    <Header />
    <Sidebar />
    <main>
      <Nuxt />
    </main>
  </div>
</template>

<style>
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
</style>
