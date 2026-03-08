<template>
    <div id="player">
        <div class="row">
            <div class="col-md-4">
                <audio preload="auto" controls=""
                       @ended="nextTrack()"
                       @play="onPlay()"
                       @pause="onPause()"
                       ref="audioPlayer">
                    <source v-if="current"
                            :src="current.URL">
                </audio>
            </div>
            <div class="col-md-8">
                <ul class="nav nav-pills" v-if="playlistMode">
                    <li role="presentation"><a href="#" @click.prevent="mode='hide'">Hide</a></li>
                    <li role="presentation" class="active"><a href="#">Playlist</a></li>
                    <li role="presentation"><a href="#" @click.prevent="mode='recent'">Recent</a></li>            
                </ul>

                <ul class="nav nav-pills" v-if="hideMode">
                    <li role="presentation" class="active"><a href="#">Hide</a></li>
                    <li role="presentation"><a href="#" @click.prevent="mode='playlist'">Playlist</a></li>
                    <li role="presentation"><a href="#" @click.prevent="mode='recent'">Recent</a></li>            
                </ul>

                <ul class="nav nav-pills" v-if="recentMode">
                    <li role="presentation"><a href="#" @click.prevent="mode='hide'">Hide</a></li>
                    <li role="presentation"><a href="#" @click.prevent="mode='playlist'">Playlist</a></li>
                    <li role="presentation" class="active"><a href="#">Recent</a></li>            
                </ul>
            </div>
        </div>

        <table class="table table-striped table-condensed" v-if="!hideMode">
            <thead>
                <tr>
                    <th></th>
                    <th>Song</th>
                    <th>Artist</th>
                    <th>Album</th>
                    <th>Rating</th>
                </tr>
            </thead>
            <tbody>

                <tr class="active info" v-if="current">
                    <td>
                        <span class="glyphicon glyphicon-fast-forward"
                              @click="nextTrack()"
                        ></span>
                    </td>
                    
                    <td>
                        <song-link :id="current.ID"
                                   :title="current.Title"></song-link>
                    </td>
                    
                    <td>
                        <artist-link :id="current.ArtistID"
                                     :name="current.Artist"></artist-link>
                    </td>
                    
                    <td>
                        <album-link :id="current.AlbumID"
                                    :name="current.Album"></album-link>
                    </td>
                    
                    <td>
                        <rating :id="current.ID"></rating>
                    </td>
                </tr>

                <template v-if="playlistMode">
                  <tr v-for="(song, index) in playlist" :key="song.ID">
                    <td>
                        <span class="glyphicon glyphicon-minus" @click="remove(index)"></span>
                        <span class="glyphicon glyphicon-chevron-up" @click="toTheTop(index)"
                              v-if="index > 0"
                        ></span>
                    </td>
                    
                    <td>
                        <song-link :id="song.ID"
                                   :title="song.Title"></song-link>
                    </td>
                    
                    <td>
                        <artist-link :id="song.ArtistID"
                                     :name="song.Artist"></artist-link>
                    </td>
                    
                    <td>
                        <album-link :id="song.AlbumID"
                                    :name="song.Album"></album-link>
                    </td>
                    
                    <td>
                        <rating :id="song.ID"></rating>
                    </td>
                  </tr>
                </template>

                <template v-if="recentMode">
                  <tr v-for="song in recentlyPlayed" :key="song.ID">
                    <td></td>
                    
                    <td>
                        <song-link :id="song.ID"
                                   :title="song.Title"></song-link>
                    </td>
                    
                    <td>
                        <artist-link :id="song.ArtistID"
                                     :name="song.Artist"></artist-link>
                    </td>
                    
                    <td>
                        <album-link :id="song.AlbumID"
                                    :name="song.Album"></album-link>
                    </td>
                    
                    <td>
                        <rating :id="song.ID"></rating>
                    </td>
                  </tr>
                </template>

            </tbody>
        </table>
        
    </div>
</template>

<script>
 import axios from 'axios' 
 import Rating from '@/components/Rating.vue'
 import SongLink from '@/components/SongLink.vue'
 import ArtistLink from '@/components/ArtistLink.vue'
 import AlbumLink from '@/components/AlbumLink.vue'  

 export default {
     name: 'Player',
     data () {
         return {
             'isPaused': false,
             'mode': 'playlist'
         }
     },
     computed: {
         playlistMode () {
             return this.mode === 'playlist'
         },
         hideMode () {
             return this.mode === 'hide'
         },
         recentMode () {
             return this.mode === 'recent'
         },
         current () {
             return this.$store.state.current
         },
         playlist () {
             return this.$store.state.playlist
         },
         recentlyPlayed () {
             return this.$store.state.recent
         }
     },
     components: {
         Rating,
         SongLink,
         ArtistLink,
         AlbumLink,
     },
     methods: {
         getRecent() {
             const recentPath = `/api/recentlyPlayed/`
             axios.get(recentPath).then(response => {
                 this.$store.commit('setRecent', response.data.Plays)
             }).catch(error => {
                 console.log(error)
             })
         },
         getData() {
             this.getRecent()
             const playlistPath = `/api/randomPlaylist/`
             axios.get(playlistPath).then(response => {
                 var songs = response.data.Songs

                 this.$store.commit('setCurrent', songs.shift())
                 this.$store.commit('setPlaylist', songs)
             }).catch(error => {
                 console.log(error)
             })
         },
         nextTrack() {
             this.isPaused = false

             this.$store.commit('nextTrack')
             if (this.$refs.audioPlayer) {
                this.$refs.audioPlayer.src = this.$store.state.current.URL
                this.$refs.audioPlayer.play()
             }
             
             // fetch a new random track to add to the end of the playlist
             this.addRandomTrack()
         },
         onPlay() {
             if (!this.isPaused) {
                 this.logPlay()
             }
             this.isPaused = false
         },
         onPause() {
             this.isPaused = true
         },
         addRandomTrack() {
             const path = `/api/random/`
             axios.get(path).then(response => {
                 this.$store.commit('appendTrack', response.data)
             })
         },
         logPlay() {
             axios.get(this.$store.state.current.PlayURL)
         },
         remove(idx) {
             this.$store.commit('remove', idx)
         },
         toTheTop(idx) {
             this.$store.commit('toTheTop', idx)
         }
     },
     created () {
         this.getData()
     }
 }
</script>

<style>
</style>