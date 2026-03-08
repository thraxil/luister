import { createApp } from 'vue'
import App from './App.vue'
import Nav from './Nav.vue'
import Player from '@/components/Player.vue'
import router from './router'
import store from './store'

const app = createApp(App)
app.use(router)
app.use(store)
app.mount('#app')

const nav = createApp(Nav)
nav.use(store)
nav.mount('#nav')

const player = createApp(Player)
player.use(store)
player.mount('#player')