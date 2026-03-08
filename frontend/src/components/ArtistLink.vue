<template>
    <a :href="'/ar/' + id + '/'">{{normalizedName}}</a>
</template>

<script>
 export default {
     name: 'ArtistLink',
     props: ["id", "name"],
     computed: {
         normalizedName () {
             let t = (this.name || '').trim().substring(0, 256)
             if (t === '') {
                 return '[missing]'
             }
             return t
         }
     }
 }
</script>