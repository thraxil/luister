<template>
    <a :href="'/al/' + id + '/'">{{normalizedName}}</a>
</template>

<script>
 export default {
     name: 'AlbumLink',
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