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
                            <a v-bind:href="'/s/' + song.ID + '/'">{{song.Title}}</a>
                        </td>
                        <td>
                            <a v-bind:href="'/ar/' + song.Artist.ID + '/'">{{song.Artist.Name}}</a>
                        </td>
                        <td>
                            <a v-bind:href="'/al/' + song.Album.ID + '/'">{{song.Album.Name}}</a>
                        </td>
                        <td>
                            <rating v-bind:id="song.ID" v-bind:initial-rating="song.Rating"></rating>
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
         'rating': Rating
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
