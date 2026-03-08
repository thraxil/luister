<template>
    <a :href="'/s/' + id + '/'">{{normalizedTitle}}</a>
</template>

<script>
 export default {
     name: 'SongLink',
     props: ["id", "title"],
     computed: {
         normalizedTitle () {
             let t = (this.title || '').trim().substring(0, 256)
             if (t === '') {
                 return '[missing]'
             }
             return t
         }
     }
 }
</script>