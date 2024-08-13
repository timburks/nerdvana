package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/timburks/nerdvana/pkg/user"
	official "google.golang.org/appengine/v2/user"
)

func main() {
	http.HandleFunc("/login", loginHandler)
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

func loginHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", 303)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	page := &Page{}
	page.Title = "Hello, world!"
	page.Header = r.Header
	page.OfficialUser = fmt.Sprintf("%+v", official.Current(r.Context()))
	page.OfficialLoginURL, err = official.LoginURL(r.Context(), r.URL.Path)
	if err != nil {
		page.OfficialLoginURL = err.Error()
		log.Printf("LoginURL error: %s", err)
	}
	page.OfficialLoginURLFederated, err = official.LoginURLFederated(r.Context(), r.URL.Path, "gmail.com")
	if err != nil {
		page.OfficialLoginURLFederated = err.Error()
		log.Printf("LoginURLFederated error: %s", err)
	}
	page.OfficialLogoutURL, err = official.LogoutURL(r.Context(), "http://"+r.Host+"/")
	if err != nil {
		page.OfficialLogoutURL = err.Error()
		log.Printf("LogoutURL error: %s", err)
	}

	hackuser := user.Current(r)

	page.HackUser = fmt.Sprintf("%+v", hackuser)
	page.HackLoginURL = user.LoginURL()
	page.HackLogoutURL = user.LogoutURL()
	if hackuser != nil {
		page.Title = "Hello, " + strings.Title(hackuser.Nickname)
		page.Prompt = "sign out"
		page.PromptURL = user.LogoutURL()
	} else {
		page.Title = "Hello, world"
		page.Prompt = "sign in"
		page.PromptURL = user.LoginURL()
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
	Title                     string
	OfficialUser              string
	OfficialLoginURL          string
	OfficialLoginURLFederated string
	OfficialLogoutURL         string
	HackUser                  string
	HackLoginURL              string
	HackLogoutURL             string
	Prompt                    string
	PromptURL                 string
	Header                    map[string][]string
}

const tmpl = `
<html>
<body>
<h1>{{ .Title }}</h1>
<h2><a href="{{ .PromptURL }}">{{ .Prompt }}</a></h2>

<hr>

<table>
<tr><td>User</td><td>{{ .OfficialUser }}</td></tr>
<tr><td>LoginURL</td><td>{{ .OfficialLoginURL }}</td></tr>
<tr><td>LoginURLFederated</td><td>{{ .OfficialLoginURLFederated }}</td></tr>
<tr><td>LogoutURL</td><td>{{ .OfficialLogoutURL }}</td></tr>
</table>

<hr>

<table>
<tr><td>User</td><td>{{ .HackUser }}</td></tr>
<tr><td>LoginURL</td><td>{{ .HackLoginURL }}</td></tr>
<tr><td>LogoutURL</td><td>{{ .HackLogoutURL }}</td></tr>
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
