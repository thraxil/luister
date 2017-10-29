<template>
    <div>
        <p>Index page</p>
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
