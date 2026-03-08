<template>
  <nav class="navbar navbar-default" id="nav">
    <div class="navbar-header">
      <router-link class="navbar-brand" to="/">
        <span class="glyphicon glyphicon-home"></span>
      </router-link>
    </div>
   
    <form class="navbar-form navbar-right" action="/vue/#/search/" method="get" >
      <div class="form-group">
        <input type="text" class="form-control" placeholder="search" name="q" />
      </div>
    </form>
               
    <ul class="nav navbar-nav">
      <li><a href="/ar/">Artists</a></li>
      <li><router-link to="/t/">Tags</router-link></li>                    
    </ul>
  </nav>
</template>

<script>
export default {
  name: 'Nav'
}
</script>