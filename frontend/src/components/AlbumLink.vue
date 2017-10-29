<template>
    <a v-bind:href="'/al/' + id + '/'">{{normalizedName}}</a>
</template>

<script>
 export default {
     name: 'AlbumLink',
     props: ["id", "name"],
     data () {
         return {}
     },
     computed: {
         normalizedName: function () {
             let t = this.name.trim().substring(0, 256)
             if (t === '') {
                 return '[missing]'
             }
             return t
         }
     }
 }

</script>

<style>
</style>
