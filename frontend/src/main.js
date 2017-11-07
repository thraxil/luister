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

const store = new Vuex.Store({
    state: {
        current: undefined,
        playlist: [],
        recent: [],
        ratings: {}
    },
    mutations: {
        setCurrent(state, c) {
            state.current = c
        },
        setPlaylist(state, p) {
            for (var i=0; i < p.length; i++) {
                state.ratings[p[i].ID] = p[i].Rating
            }
            state.playlist = p
        },
        setRecent(state, r) {
            for (var i=0; i < r.length; i++) {
                state.ratings[r[i].ID] = r[i].Rating
            }
            state.recent = r
        },
        nextTrack(state) {
            // new current track
            var newCurrent = state.playlist.shift()
             
            // move the old one over to the recently played list
            state.recent.unshift(state.current)
            // then trim it
            state.recent.splice(-1, 1)
            state.current = newCurrent
        },
        appendTrack(state, track) {
            Vue.set(state.ratings, track.ID, track.Rating)
            state.playlist.push(track);
        },
        remove(state, idx) {
            state.playlist.splice(idx, 1)
        },
        toTheTop(state, idx) {
            var s = state.playlist[idx]
            state.playlist.splice(idx, 1)
            state.playlist.unshift(s)
        },
        setRating(state, song) {
            Vue.set(state.ratings, song.ID, song.Rating)
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
