<template>
    <div id="player">
        <audio preload="auto" controls=""
               v-on:ended="nextTrack()"
               v-on:play="onPlay()"
               v-on:pause="onPause()"
        >
            <source v-if="current"
                    v-bind:src="current.URL">
        </audio>

        <table class="table table-striped table-condensed">
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
                    <td></td>
                    
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
                                :key="current"
                        ></rating>
                    </td>
                </tr>

                
                <tr v-for="song in playlist">
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
                                v-bind:initial-rating="song.Rating"></rating>
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
             'recentlyPlayed':  [],
             'playlist': [],
             'current': undefined,
             'isPaused': false,
             'audio': undefined,
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
                 this.recentlyPlayed = response.data.Plays
             }).catch(error => {
                 console.log(error)
             })
         },
         getData() {
             this.getRecent()
             const playlistPath = `/api/randomPlaylist/`
             axios.get(playlistPath).then(response => {
                 var songs = response.data.Songs
                 
                 this.current = songs.shift()
                 
                 this.playlist = songs
             }).catch(error => {
                 console.log(error)
             })
         },
         nextTrack() {
             this.isPaused = false
             // append current to recentlyPlayed and pop the end off that
             // this.getRecent()
             
             // new current track
             var newCurrent = this.playlist.shift()
             
             this.current = newCurrent
             this.audio.src = this.current.URL
             
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
                 this.playlist.push(response.data);
             })
         },
         logPlay() {
             axios.get(this.current.PlayURL)
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
