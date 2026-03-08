<template>
    <div>
        <h2>Tags</h2>
        <div class="list-group" v-if="tags.length">
            <router-link v-for="tag in tags" :key="tag.Name" :to="'/t/' + tag.Name + '/'" class="list-group-item">
                {{tag.Name}}
            </router-link>
        </div>
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