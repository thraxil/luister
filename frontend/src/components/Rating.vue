<template>
    <nobr>
        <span v-bind:class="'rating rating-' + rating">
            <span v-for="n in 5"  v-bind:class="ratingIcon(n)" v-on:click="setRating(n)"></span>
        </span>
    </nobr>
</template>

<script>
 export default {
     name: 'Rating',
     props: ['id'],
     data () {
         return {
             'id': this.id
         }
     },
     computed: {
         'rating': function () {
             return this.$store.state.ratings[this.id]
         }
     },
     methods: {
         ratingIcon: function(level) {
             if (this.rating >= level) {
                 return "glyphicon glyphicon-star"
             } else {
                 return "glyphicon glyphicon-star-empty"
             }
         },
         setRating: function(level) {
             this.$store.commit('setRating', {'ID': this.id, 'Rating': level})
             const path = `/r/` + this.id + `/`
             var data = new FormData()
             data.append('rating', level)
             var request = new XMLHttpRequest();
             request.open('POST', path);
             request.send(data);
         }
     }
 }
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
 .rating-0 { color: #ccc; }
 .rating-1 { color: #999; }
 .rating-2 { color: #666; }
 .rating-3 { color: #333; }
 .rating-4 { color: #c76; }
 .rating-5 { color: #f60; }
 .rating-0 .glyphicon-star-empty { color: #ccc; }
 .rating-1 .glyphicon-star-empty { color: #ccc; }
 .rating-2 .glyphicon-star-empty { color: #ccc; }
 .rating-3 .glyphicon-star-empty { color: #ccc; }
 .rating-4 .glyphicon-star-empty { color: #ccc; }
 .rating-5 .glyphicon-star-empty { color: #ccc; }         
</style>
