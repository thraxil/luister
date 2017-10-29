<template>
    <div>
        <h2>Recently Played <span class="glyphicon glyphicon-refresh pull-right" v-on:click="getData()"></span></h2>

        <table class="table table-striped table-condensed" v-if="recentlyPlayed.length">
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
                <tr v-for="play in recentlyPlayed">
                    <td></td>
                    <td>
                        <song-link v-bind:id="play.Song.ID"
                                   v-bind:title="play.Song.Title"></song-link>
                    </td>
                    <td>
                        <artist-link v-bind:id="play.Song.Artist.ID"
                                     v-bind:name="play.Song.Artist.Name"></artist-link>
                    </td>
                    <td>
                        <album-link v-bind:id="play.Song.Album.ID"
                                     v-bind:name="play.Song.Album.Name"></album-link>
                    </td>
                    <td>
                        <rating v-bind:id="play.Song.ID"
                                v-bind:initial-rating="play.Song.Rating"></rating>
                    </td>
                </tr>
            </tbody>

        </table>

        <p v-else>No recently played songs</p>

    </div>
</template>

<script>
 import axios from 'axios' 
 import Rating from '@/components/Rating'
 import SongLink from '@/components/SongLink'
 import ArtistLink from '@/components/ArtistLink'
 import AlbumLink from '@/components/AlbumLink'  
 
 export default {
     name: 'Index',
     data () {
         return {
             'recentlyPlayed': []
         }
     },
     components: {
         'rating': Rating,
         'song-link': SongLink,
         'artist-link': ArtistLink,
         'album-link': AlbumLink,
     },
     methods: {
         getData() {
             const path = `/api/recentlyPlayed/`
             axios.get(path).then(response => {
                 this.recentlyPlayed = response.data.Plays
             }).catch(error => {
                 console.log(error)
             })
         }

     },
     created () {
         this.getData()
     }
 }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style>
</style>
