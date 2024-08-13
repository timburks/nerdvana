package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"google.golang.org/appengine/v2/user"
)

func main() {
	http.HandleFunc("/", indexHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	var err error

	page := &Page{}
	page.Title = "Hello, world!"
	page.Header = r.Header
	page.User = fmt.Sprintf("%+v", user.Current(r.Context()))
	page.LoginURL, err = user.LoginURL(r.Context(), r.URL.Path)
	if err != nil {
		page.LoginURL = err.Error()
		log.Printf("LoginURL error: %s", err)
	}
	page.LoginURLFederated, err = user.LoginURLFederated(r.Context(), r.URL.Path, "gmail.com")
	if err != nil {
		page.LoginURLFederated = err.Error()
		log.Printf("LoginURLFederated error: %s", err)
	}
	page.LogoutURL, err = user.LogoutURL(r.Context(), "http://"+r.Host+"/")
	if err != nil {
		page.LogoutURL = err.Error()
		log.Printf("LogoutURL error: %s", err)
	}

	t, err := template.New("page").Parse(tmpl)
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, page)
	if err != nil {
		panic(err)
	}

}

type Page struct {
	Title             string
	User              string
	LoginURL          string
	LoginURLFederated string
	LogoutURL         string
	Header            map[string][]string
}

const tmpl = `
<html>
<body>
<h1>{{ .Title }}</h1>

<hr>

<table>
<tr><td>User</td><td>{{ .User }}</td></tr>
<tr><td>LoginURL</td><td>{{ .LoginURL }}</td></tr>
<tr><td>LoginURLFederated</td><td>{{ .LoginURLFederated }}</td></tr>
<tr><td>LogoutURL</td><td>{{ .LogoutURL }}</td></tr>
</table>

<hr>

<table>
{{ range $k,$v := .Header }}
<tr><td>{{ $k }}</td><td>{{ $v }}</td></tr>
{{ end }}
</table>


</body>
</html>
`
