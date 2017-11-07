<template>
    <table class="table">
        <tr>
            <th>Tag</th>
            <td>{{tag.Name}}</td>
        </tr>

        <tr>
            <th>Songs</th>
            <td>
                <table class="table table-striped table-condensed" id="songs">
                    <tr v-for="song in songs">
                        <td>
                            <span class="glyphicon glyphicon-plus" v-on:click="addToPlaylist(song)"></span>
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
                                    v-bind:initial-rating="song.Rating"></rating>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>

    </table>
</template>

<script>
 import axios from 'axios'
 import Rating from '@/components/Rating'
 import SongLink from '@/components/SongLink'
 import ArtistLink from '@/components/ArtistLink'
 import AlbumLink from '@/components/AlbumLink'  


 export default {
     name: 'Tag',
     data () {
         return {
             tag: {},
             songs: []
         }
     },
     computed: {
     },
     components: {
         'rating': Rating,
         'song-link': SongLink,
         'artist-link': ArtistLink,
         'album-link': AlbumLink,
     },
     methods: {
         getData() {
             const path = `/api/tags/` + this.$route.params.tagname + `/`
             axios.get(path).then(response => {
                 this.tag = response.data.Tag
                 for (var i=0; i < response.data.Songs.length; i++) {
                     this.$store.commit('setRating', response.data.Songs[i])
                 }
                 this.songs = response.data.Songs
             }).catch(error => {
                 console.log(error)
             })
         },
         addToPlaylist(song) {
             this.$store.commit('appendTrack', song)
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
