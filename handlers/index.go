package handlers

import (
	"html/template"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	indexHTML := `<!DOCTYPE HTML>
	<html>
		<head>
			<meta charset="utf-8">
			<title>WebMine</title>
			<link href="/static/stylesheets/styles.css" rel="stylesheet">
		</head>
		<body>
	  		<div id="root"></div>
			<script src="/frontend/dist/app.js" type="text/javascript"></script>
		</body>
	</html>
	`
	indexTemplate := template.Must(template.New("index").Parse(indexHTML))
	indexTemplate.Execute(w, nil)
}
