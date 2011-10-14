package flipcoin

import (
    "appengine"
    "appengine/datastore"
    "http"
    "time"
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
}

func root(w http.ResponseWriter, r *http.Request) {
 
}

func show(w http.ResponseWriter, r *http.Request) {
}
