<template>
  <nav class="nav">
    <NuxtLink to="/" exact class="brand">
      Mattrax - <span>{{ $store.state.authentication.user.org }}</span>
    </NuxtLink>

    <div class="navRight">
      <!-- <span class="navNotifications"><NotificationIcon /></span> -->

      <div class="dropdown">
        <span class="navUser">
          {{ $store.state.authentication.user.name }}
          <CaretIcon />
        </span>
        <div class="dropdown-content">
          <NuxtLink to="/settings/users"> Edit Account </NuxtLink>
          <!-- <a href="#" @click.prevent="router.">Edit Account</a> -->
          <a href="#" @click.prevent="logout()">Logout</a>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
// import NotificationIcon from '@/assets/icon/notification.svg?inline'
import CaretIcon from '@/assets/icon/caret.svg?inline'

export default {
  components: { /* NotificationIcon, */ CaretIcon },
  methods: {
    logout() {
      this.$store
        .dispatch('authentication/logout')
        .then(() => {
          this.$router.push('/login')
        })
        .catch((err) => {
          console.error(err)
        })
    },
  },
}
</script>

<style>
.nav {
  height: 50px;
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  background: var(--primary-color);
  color: var(--light-text-color);
  padding: 0 15px;
}

.brand span {
  font-size: 14px;
}

.navRight {
  float: right;
  height: 100%;
  line-height: 50px;
}

.navRight svg {
  display: inline-block;
  vertical-align: middle;
}

.navNotifications {
  float: left;
  padding: 0 7px;
}

.navUser {
  /* float: right; */
  padding: 0 5px 0 7px;
  cursor: default;
}

.dropdown {
  position: relative;
  display: inline-block;
}

.dropdown-content {
  display: none;
  position: absolute;
  background-color: #f9f9f9;
  min-width: 160px;
  box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
  z-index: 1;
}

.dropdown-content a {
  color: black;
  padding: 0px 12px;
  text-decoration: none;
  display: block;
}

.dropdown-content a:hover {
  background-color: #f1f1f1;
}

.dropdown:hover .dropdown-content {
  display: block;
}
</style>
