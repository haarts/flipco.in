package flipcoin

import (
    "fmt"
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
  Done         bool
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
  c := appengine.NewContext(r)
  count, err := datastore.NewQuery("Coinflip").Count(c)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
  fmt.Println(count)
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/home.html", map[string]string{"title":"Awesome coin tosses - Flipco.in", "nr_of_flips":fmt.Sprint(count)}))
}

func create(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Redirect(w, r, "/", 302)
    /*return*/
  }

  c := appengine.NewContext(r)

  r.ParseForm()
  tail    := r.Form["tail"][0]
  head    := r.Form["head"][0]
  friends := r.Form["friends"]

  /*if tail == "" || head == "" || friends == nil {*/
    /*http.Redirect(w, r, "/", 302)*/
    /*return*/
  /*}*/

  participants := make([]Participant, len(friends))
  for i := range friends {
    participants[i] = Participant{Email: friends[i], Seen: datastore.SecondsToTime(time.Seconds())}
  }

  coin := Coinflip {
    Head: head,
    Tail: tail,
    Participants: participants,
    Done: false,
  }

  _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Coinflip", nil), &coin)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }


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
