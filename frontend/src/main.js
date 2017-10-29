// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import Nav from './Nav'
import Player from '@/components/Player'
import router from './router'

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  template: '<App/>',
  components: { App }
})

new Vue({
  el: '#nav',
  template: '<Nav/>',
  components: { Nav }
})

new Vue({
    el: '#player',
    template: '<Player/>',
    components: { Player }
})
