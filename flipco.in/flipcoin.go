package flipcoin

import (
    "fmt"
    "appengine"
    "appengine/datastore"
    "http"
    /*"time"*/
    "mustache.go"
    "os"
    "strings"
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
  if len(r.URL.Path) != 1 {
    http.ServeFile(w, r, "./flipco.in/views" + r.URL.Path)
    return
  }
  c := appengine.NewContext(r)
  count, err := datastore.NewQuery("Coinflip").Count(c)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
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
    participants[i] = Participant{Email: friends[i]}
  }

  coin := Coinflip {
    Head: head,
    Tail: tail,
    Participants: participants,
    Done: false,
  }

  key, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Coinflip", nil), &coin)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }

  http.Redirect(w, r, "/show/" + key.Encode(), 302)
}

func show(w http.ResponseWriter, r *http.Request) {
/*, Seen: datastore.SecondsToTime(time.Seconds())*/
  context  := appengine.NewContext(r)
  key_as_string := strings.Split(r.URL.Path, "/")[2]
  coinflip, _ := find(key_as_string, context)
  /*coinflip := Coinflip.find(uuid)*/ /* this is conceptually wrong */
  /*registerParticipant(email)*/
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/show.html", map[string]string{"participants":fmt.Sprint(len(coinflip.Participants))}))
}

func find(key_as_string string, context appengine.Context) (*Coinflip, os.Error) {
  coinflip := new(Coinflip)
  key, _ := datastore.DecodeKey(key_as_string)
  if err := datastore.Get(context, key, coinflip); err != nil {
    return nil, err
  }
  return coinflip, nil
}

func (p *Coinflip) mailParticipants() {

}

  /*fmt.Println(r.Form)*/
  /*for k, v := range r.Form {*/
    /*fmt.Println(k)*/
    /*fmt.Printf("%T\n",v)*/
    /*for _, sv := range v {*/
      /*fmt.Println(sv)*/
    /*}*/
  /*}*/

