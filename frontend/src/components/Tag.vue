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
                            <song-link v-bind:id="song.ID"
                                       v-bind:title="song.Title"></song-link>
                        </td>
                        <td>
                            <artist-link v-bind:id="song.Artist.ID"
                                         v-bind:name="song.Artist.Name"></artist-link>
                        </td>
                        <td>
                            <album-link v-bind:id="song.Album.ID"
                                        v-bind:name="song.Album.Name"></album-link>
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
                 this.songs = response.data.Songs
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
