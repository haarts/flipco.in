package flipcoin

import (
    "fmt"
    /*"appEngine"*/
    "appengine/datastore"
    "http"
    /*"time"*/
    "mustache.go"
)

type Coinflip struct {
  Participants []Participant
  Head         string
  Tail         string
}

type Participant struct {
  Email string
  Seen  datastore.Time
}

func init() {
  http.HandleFunc("/", root)
  http.HandleFunc("/show/", show)
  http.HandleFunc("/create", create)
}

func root(w http.ResponseWriter, r *http.Request) {
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/home.html", map[string]string{"title":"Awesome coin tosses - Flipco.in", "nr_of_flips":"123"}))
}

func create(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Redirect(w, r, "/", 302)
  }

  r.ParseForm()
  fmt.Println(r.Form)
  for k, v := range r.Form {
    fmt.Println(k)
    fmt.Printf("%T\n",v)
    for _, sv := range v {
      fmt.Println(sv)
    }
  }
  http.Redirect(w, r, "/show/test_hash", 302)
}

func show(w http.ResponseWriter, r *http.Request) {

}
