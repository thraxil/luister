<template>
    <table class="table" v-if="tag">
        <tr>
            <th>Tag</th>
            <td>{{tag.Name}}</td>
        </tr>

        <tr>
            <th>Songs</th>
            <td>
                <table class="table table-striped table-condensed" id="songs">
                    <tr v-for="song in songs" :key="song.ID">
                        <td>
                            <span class="glyphicon glyphicon-plus" @click="addToPlaylist(song)"></span>
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
                </table>
            </td>
        </tr>

    </table>
</template>

<script>
 import axios from 'axios'
 import Rating from '@/components/Rating.vue'
 import SongLink from '@/components/SongLink.vue'
 import ArtistLink from '@/components/ArtistLink.vue'
 import AlbumLink from '@/components/AlbumLink.vue'  

 export default {
     name: 'Tag',
     data () {
         return {
             tag: null,
             songs: []
         }
     },
     components: {
         Rating,
         SongLink,
         ArtistLink,
         AlbumLink,
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
     },
     watch: {
        '$route.params.tagname': 'getData'
     }
 }
</script>