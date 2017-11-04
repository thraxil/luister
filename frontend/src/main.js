// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import Vuex from 'vuex'
import App from './App'
import Nav from './Nav'
import Player from '@/components/Player'
import router from './router'

Vue.use(Vuex)
Vue.config.productionTip = false

window.bus = new Vue();

const store = new Vuex.Store({
    state: {
        current: undefined,
        playlist: [],
        recent: []
    },
    mutations: {
        setCurrent(state, c) {
            state.current = c
        }
    }
})

/* eslint-disable no-new */
new Vue({
    el: '#app',
    store,
    router,
    template: '<App/>',
    components: { App }
})

new Vue({
    el: '#nav',
    store,
    template: '<Nav/>',
    components: { Nav }
})

new Vue({
    el: '#player',
    store,
    template: '<Player/>',
    components: { Player }
})
