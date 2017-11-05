<template>
    <div id="player">
        <div class="row">
            <div class="col-md-4">
                <audio preload="auto" controls=""
                       v-on:ended="nextTrack()"
                       v-on:play="onPlay()"
                       v-on:pause="onPause()">
                    <source v-if="current"
                            v-bind:src="current.URL">
                </audio>
            </div>
            <div class="col-md-8">
                <ul class="nav nav-pills" v-if="playlistMode">
                    <li role="presentation"><a href="#" v-on:click.prevent="mode='hide'">Hide</a></li>
                    <li role="presentation" class="active"><a href="#">Playlist</a></li>
                    <li role="presentation"><a href="#" v-on:click.prevent="mode='recent'">Recent</a></li>            
                </ul>

                <ul class="nav nav-pills" v-if="hideMode">
                    <li role="presentation" class="active"><a href="#">Hide</a></li>
                    <li role="presentation"><a href="#" v-on:click.prevent="mode='playlist'">Playlist</a></li>
                    <li role="presentation"><a href="#" v-on:click.prevent="mode='recent'">Recent</a></li>            
                </ul>

                <ul class="nav nav-pills" v-if="recentMode">
                    <li role="presentation"><a href="#" v-on:click.prevent="mode='hide'">Hide</a></li>
                    <li role="presentation"><a href="#" v-on:click.prevent="mode='playlist'">Playlist</a></li>
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
                              v-on:click="nextTrack()"
                        ></span>
                    </td>
                    
                    <td>
                        <song-link v-bind:id="current.ID"
                                   v-bind:title="current.Title"></song-link>
                    </td>
                    
                    <td>
                        <artist-link v-bind:id="current.ArtistID"
                                     v-bind:name="current.Artist"></artist-link>
                    </td>
                    
                    <td>
                        <album-link v-bind:id="current.AlbumID"
                                    v-bind:name="current.Album"></album-link>
                    </td>
                    
                    <td>
                        <rating v-bind:id="current.ID"
                                v-bind:initial-rating="current.Rating"
                                :key="current.ID"
                        ></rating>
                    </td>
                </tr>

                
                <tr v-for="(song, index) in playlist" v-if="playlistMode" :key="song.ID">
                    <td>
                        <span class="glyphicon glyphicon-minus" v-on:click="remove(index)"></span>
                        <span class="glyphicon glyphicon-chevron-up" v-on:click="toTheTop(index)"
                              v-if="index > 0"
                        ></span>
                    </td>
                    
                    <td>
                        <song-link v-bind:id="song.ID"
                                   v-bind:title="song.Title"></song-link>
                    </td>
                    
                    <td>
                        <artist-link v-bind:id="song.ArtistID"
                                     v-bind:name="song.Artist"></artist-link>
                    </td>
                    
                    <td>
                        <album-link v-bind:id="song.AlbumID"
                                    v-bind:name="song.Album"></album-link>
                    </td>
                    
                    <td>
                        <rating v-bind:id="song.ID"
                                v-bind:initial-rating="song.Rating"
                                :key="song.ID"
                        ></rating>
                    </td>
                </tr>

                <tr v-for="song in recentlyPlayed" v-if="recentMode" :key="song.ID">
                    <td></td>
                    
                    <td>
                        <song-link v-bind:id="song.ID"
                                   v-bind:title="song.Title"></song-link>
                    </td>
                    
                    <td>
                        <artist-link v-bind:id="song.ArtistID"
                                     v-bind:name="song.Artist"></artist-link>
                    </td>
                    
                    <td>
                        <album-link v-bind:id="song.AlbumID"
                                    v-bind:name="song.Album"></album-link>
                    </td>
                    
                    <td>
                        <rating v-bind:id="song.ID"
                                v-bind:initial-rating="song.Rating"
                                :key="song.ID"
                        ></rating>
                    </td>
                </tr>

            </tbody>
        </table>
        
    </div>
</template>

<script>
 import axios from 'axios' 
 import Rating from '@/components/Rating'
 import SongLink from '@/components/SongLink'
 import ArtistLink from '@/components/ArtistLink'
 import AlbumLink from '@/components/AlbumLink'  

 export default {
     name: 'Player',
     data () {
         return {
             'isPaused': false,
             'audio': undefined,
             'mode': 'playlist'
         }
     },
     computed: {
         'playlistMode': function () {
             return this.mode === 'playlist'
         },
         'hideMode': function () {
             return this.mode === 'hide'
         },
         'recentMode': function () {
             return this.mode === 'recent'
         },
         'current': function () {
             return this.$store.state.current
         },
         'playlist': function () {
             return this.$store.state.playlist
         },
         'recentlyPlayed': function () {
             return this.$store.state.recent
         }
     },
     components: {
         'rating': Rating,
         'song-link': SongLink,
         'artist-link': ArtistLink,
         'album-link': AlbumLink,
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
             this.audio.src = this.$store.state.current.URL
             
             this.audio.play()
             
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
     },
     mounted () {
         this.audio = this.$el.querySelectorAll('audio')[0]
     }
 }
</script>

<style>
</style>
