<template>
  <nav class="navbar navbar-default" id="nav">
    <div class="navbar-header">
      <a class="navbar-brand" href="/vue/#/">
        <span class="glyphicon glyphicon-home"></span>
      </a>
    </div>
   
    <form class="navbar-form navbar-right" action="/vue/#/search/" method="get" >
      <div class="form-group">
        <input type="text" class="form-control" placeholder="search" name="q" />
      </div>
    </form>
               
    <ul class="nav navbar-nav">
      <li><a href="/ar/">Artists</a></li>
      <li><a href="/vue/#/t/">Tags</a></li>                    
    </ul>
  </nav>
</template>

<script>
export default {
  name: 'nav'
}
</script>

<style>
</style>
