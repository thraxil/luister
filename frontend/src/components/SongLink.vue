<template>
    <a v-bind:href="'/s/' + id + '/'">{{normalizedTitle}}</a>
</template>

<script>
 export default {
     name: 'SongLink',
     props: ["id", "title"],
     data () {
         return {}
     },
     computed: {
         normalizedTitle: function () {
             let t = this.title.trim().substring(0, 256)
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
