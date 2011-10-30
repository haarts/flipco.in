package flipcoin

import (
    "fmt"
    "appengine"
    "appengine/datastore"
    "http"
    "time"
    "mustache.go"
    "os"
    "strings"
)

type Coinflip struct {
  Participants []*datastore.Key
  Head         string
  Tail         string
  Done         bool
}

type Participant struct {
  Email string
  Seen  datastore.Time
}

func init() {
  http.HandleFunc("/", home)
  http.HandleFunc("/show/", show)
  http.HandleFunc("/create", create)
  http.HandleFunc("/register/", register)
}

func home(w http.ResponseWriter, r *http.Request) {
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
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body":mustache.RenderFile("./flipco.in/views/home.html", map[string]string{"title":"Awesome coin tosses - Flipco.in", "nr_of_flips":fmt.Sprint(count)})}))
}

func register(w http.ResponseWriter, r *http.Request) {
  if r.Method != "GET" {
    http.Redirect(w, r, "/", 302)
    return
  }

  context  := appengine.NewContext(r)
  key_as_string := strings.Split(r.URL.Path, "/")[2]
  coinflip, _ := find(key_as_string, context)
  participants, _ := coinflip.fetchParticipants(context)

  var found *Participant
  for i := 0; i < len(participants) && found == nil; i++ {
    if participants[i].Email == r.FormValue("email") {
      found = &participants[i]
    }
  }
  (*found).Seen = datastore.SecondsToTime(time.Seconds())
  /* Crap now how do I save it? */
  fmt.Println(*found)
  /*http.Redirect(w, r, "/show/" + key_as_string, 302)*/
}

func create(w http.ResponseWriter, r *http.Request) {
  if r.Method != "POST" {
    http.Redirect(w, r, "/", 302)
    return
  }

  c := appengine.NewContext(r)

  r.ParseForm()
  tail    := r.Form["tail"][0]
  head    := r.Form["head"][0]
  friends := r.Form["friends[]"]

  /*if tail == "" || head == "" || friends == nil {*/
    /*http.Redirect(w, r, "/", 302)*/
    /*return*/
  /*}*/
  participants, err := storeParticipants(friends, c)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
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
  context  := appengine.NewContext(r)
  key_as_string := strings.Split(r.URL.Path, "/")[2]
  coinflip, _ := find(key_as_string, context)
  participants, _ := coinflip.fetchParticipants(context)

  /*registerParticipant(email)*/
  str_to_str   := map[string]string{"count":fmt.Sprint(len(coinflip.Participants))}
  str_to_slice := map[string][]map[string]string{"participants":participantsMap(participants, func(p Participant) map[string]string {
    return map[string]string{"email":p.Email}
  })}
  /*str_to_slice := map[string][]map[string]string{"participants":{{"email":"a"},{"email":"b"},{"email":"c"}}}*/
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/show.html", str_to_str, str_to_slice))
}

func (coinflip Coinflip) fetchParticipants(context appengine.Context) ([]Participant, os.Error) {
  participants := make([]Participant, len(coinflip.Participants))
  for i := range participants {
    if err := datastore.Get(context, coinflip.Participants[i], &participants[i]); err != nil {
      return nil, err
    }
  }
  return participants, nil
}

func find(key_as_string string, context appengine.Context) (*Coinflip, os.Error) {
  coinflip := new(Coinflip)
  key, _ := datastore.DecodeKey(key_as_string)
  if err := datastore.Get(context, key, coinflip); err != nil {
    return nil, err
  }
  return coinflip, nil
}

func participantsMap(participants []Participant, f func(Participant) map[string]string) (mapped []map[string]string) {
  for v := range participants {
    mapped = append(mapped, f(participants[v]))
  }
  return mapped
}

func (p *Coinflip) mailParticipants() {

}

/*surely passing around Context all the time is ugly as hell*/
func storeParticipants(emails []string, context appengine.Context) ([]*datastore.Key, os.Error) {
  participants := make([]*datastore.Key, len(emails))
  for i := range emails {
    p := Participant{Email: emails[i]}
    key, err := datastore.Put(context, datastore.NewIncompleteKey(context, "Participant", nil), &p)
    if err != nil {
      return nil, err
    }
    participants[i] = key
  }
  return participants, nil
}

  /*fmt.Println(r.Form)*/
  /*for k, v := range r.Form {*/
    /*fmt.Println(k)*/
    /*fmt.Printf("%T\n",v)*/
    /*for _, sv := range v {*/
      /*fmt.Println(sv)*/
    /*}*/
  /*}*/

