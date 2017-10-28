<template>
    <div>
        <h2>Recently Played</h2>

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
                    <td><a v-bind:href="'/s/' + play.Song.ID + '/'">{{play.Song.Title}}</a></td>
                    <td><a v-bind:href="'/ar/' + play.Song.Artist.ID + '/'">{{play.Song.Artist.Name}}</a></td>
                    <td><a v-bind:href="'/al/' + play.Song.Album.ID + '/'">{{play.Song.Album.Name}}</a></td>
                    <td>
                        <rating v-bind:id="play.Song.ID" v-bind:initial-rating="play.Song.Rating"></rating>
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
 
 export default {
     name: 'Index',
     data () {
         return {
             'recentlyPlayed': []
         }
     },
     components: {
         'rating': Rating
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
