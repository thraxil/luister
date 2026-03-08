<template>
    <div>
        <p>Index page</p>
    </div>
</template>

<script>
 import Rating from '@/components/Rating.vue'
 import SongLink from '@/components/SongLink.vue'
 import ArtistLink from '@/components/ArtistLink.vue'
 import AlbumLink from '@/components/AlbumLink.vue'  
 
 export default {
     name: 'Index',
     components: {
         Rating,
         SongLink,
         ArtistLink,
         AlbumLink,
     }
 }
</script>