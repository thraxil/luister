package main

const indexTemplate = `
<html>
<head>
<title>{{.Title}}</title>
<link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.1/css/bootstrap.min.css" />
</head>
<body>
<div class="container">
Total Songs: {{.TotalSongs}}

<h2>Recently Added</h2>
<table class="table table-striped table-condensed">
<thead>
<tr><th>Song</th><th>Artist</th><th>Album</th></tr>
</thead>
<tbody>
{{ range .RecentSongs }}
<tr>
<td>{{.Title}}</td>
<td>{{.Artist.Name}}</td>
<td>{{.Album.Name}}</td>
</tr>
{{ end }}
</tbody>
</table>
</div>
</body>
</html>
`
