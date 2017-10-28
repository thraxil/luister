<template>
    <div>
        <h2>Tags</h2>
        <div class="list-group" v-if="tags.length">
        <a v-for="tag in tags" v-bind:href="'/vue/#/t/' + tag.Name + '/'" class="list-group-item">{{tag.Name}}</a>
    </div>
</template>

<script>
 import axios from 'axios'
 
 export default {
     name: 'Tags',
     data () {
         return {
             tags: []
         }
     },
     computed: {
     },
     methods: {
         getTags() {
             const path = `/api/tags/`
             axios.get(path)
                  .then(response => {
                      this.tags = response.data.Tags
                  })
                  .catch(error => {
                      console.log(error)
                  })
             
         }
     },
     created () {
         this.getTags()
     }
 }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
