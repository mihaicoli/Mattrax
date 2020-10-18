<template>
  <nav class="nav">
    <NuxtLink to="/" exact class="brand">
      Mattrax - <span>{{ $store.state.authentication.user.org }}</span>
    </NuxtLink>

    <div class="navRight">
      <!-- <span class="navNotifications"><NotificationIcon /></span> -->
      <span class="navUser">
        {{ $store.state.authentication.user.name }}
        <CaretIcon />
      </span>
    </div>
  </nav>
</template>

<script>
// import NotificationIcon from '@/assets/icon/notification.svg?inline'
import CaretIcon from '@/assets/icon/caret.svg?inline'

export default {
  components: { /* NotificationIcon, */ CaretIcon },
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
  float: right;
  padding: 0 5px 0 7px;
  cursor: default;
}

.navUser:hover,
.navNotifications:hover {
  background-color: var(--primary-color-accent);
}
</style>
