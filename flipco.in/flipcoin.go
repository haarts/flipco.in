package flipcoin

import (
    "fmt"
    "appengine"
    "appengine/datastore"
    "appengine/mail"
    "http"
    "time"
    "mustache.go"
    "os"
    "strings"
)

type Coinflip struct {
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
  /* static file serve */
  if len(r.URL.Path) != 1 {
    http.ServeFile(w, r, "./flipco.in/views" + r.URL.Path)
    return
  }

  /* the real root */
  c := appengine.NewContext(r)
  count, err := datastore.NewQuery("Coinflip").Count(c)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }

  /*long, very long line */
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body":mustache.RenderFile("./flipco.in/views/home.html", map[string]string{"title":"Awesome coin tosses - Flipco.in", "nr_of_flips":fmt.Sprint(count)})}))
}

func register(w http.ResponseWriter, r *http.Request) {
  /* not a GET request? redirect to home */
  if r.Method != "GET" {
    http.Redirect(w, r, "/", 302)
    return
  }

  context        := appengine.NewContext(r)
  coinflipKey, _ := datastore.DecodeKey(strings.Split(r.URL.Path, "/")[2])

  iterator := datastore.NewQuery("Participant").Ancestor(coinflipKey).Filter("Email =", r.FormValue("email")).Run(context)
  var found Participant
  key, err := iterator.Next(&found)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }
  found.Seen = datastore.SecondsToTime(time.Seconds())
  datastore.Put(context, key, &found)
  /*http.Redirect(w, r, "/show/" + coinflipKey.Encode(), 302)*/
}

func create(w http.ResponseWriter, r *http.Request) {
  /* not a POST request? redirect to root */
  if r.Method != "POST" {
    http.Redirect(w, r, "/", 302)
    return
  }

  c := appengine.NewContext(r)

  r.ParseForm()
  tail    := r.Form["tail"][0]
  head    := r.Form["head"][0]
  friends := r.Form["friends[]"]

  if tail == "" || head == "" || friends == nil {
    http.Redirect(w, r, "/", 302)
    return
  }

  coin := Coinflip {
    Head: head,
    Tail: tail,
    Done: false,
  }

  coinflipKey, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Coinflip", nil), &coin)
  if err != nil {
    http.Error(w, err.String(), http.StatusInternalServerError)
    return
  }

  for i := range friends {
    key := datastore.NewIncompleteKey(c, "Participant", coinflipKey)
    participant := Participant{Email:friends[i]}
    datastore.Put(c, key, &participant)
  }
  coin.mailParticipants(c, coinflipKey)

  http.Redirect(w, r, "/show/" + coinflipKey.Encode(), 302)
}

func show(w http.ResponseWriter, r *http.Request) {
  context        := appengine.NewContext(r)
  coinflipKey, _ := datastore.DecodeKey(strings.Split(r.URL.Path, "/")[2])
  coinflip, _    := find(coinflipKey, context)

  iterator := datastore.NewQuery("Participant").Ancestor(coinflipKey).Run(context)

  email_list := participantsMap(iterator, func(p Participant) map[string]string {
    return map[string]string{"email":p.Email, "seen_at":p.Seen.Time().Format(time.ANSIC)}
  })
  str_to_str   := map[string]string{"count":fmt.Sprint(len(email_list)),"head":coinflip.Head, "tail":coinflip.Tail}
  str_to_slice := map[string][]map[string]string{"participants":email_list}
  fmt.Fprint(w, mustache.RenderFile("./flipco.in/views/layout.html", map[string]string{"body":mustache.RenderFile("./flipco.in/views/show.html", str_to_str, str_to_slice)}))
}

func find(key *datastore.Key, context appengine.Context) (*Coinflip, os.Error) {
  coinflip := new(Coinflip)
  if err := datastore.Get(context, key, coinflip); err != nil {
    return nil, err
  }
  return coinflip, nil
}

func participantsMap(iterator *datastore.Iterator, f func(Participant) map[string]string) (mapped []map[string]string) {
  var participant Participant
  for _, err := iterator.Next(&participant); ; _, err = iterator.Next(&participant) {
    if err == datastore.Done {
      break
    }
    if err != nil {
      break
    }
    mapped = append(mapped, f(participant))
  }
  return mapped
}

/* passing the Context, again */
/* this is a function on a pointer to a Coinflip struct. Yet either Context OR a slice of Participant must be passed as an argument. */
func (coinflip *Coinflip) mailParticipants(context appengine.Context, coinflipKey *datastore.Key) {
  query := datastore.NewQuery("Participant").Ancestor(coinflipKey)

  for t := query.Run(context) ; ; {
    var participant Participant
    _, err := t.Next(&participant)
    if err == datastore.Done {
      break
    }

    msg := &mail.Message{
                  Sender:  "harm@flipco.in",
                  ReplyTo: "harm@flipco.in",
                  To:      []string{participant.Email},
                  Subject: "What will it be? " + coinflip.Head + " or " + coinflip.Tail + "?",
                  Body:    fmt.Sprintf(confirmMessage, "http://www.flipco.in/register/" + coinflipKey.Encode() + "?email=" + participant.Email),
          }
    if err := mail.Send(context, msg); err != nil {
            context.Errorf("Couldn't send email: %v", err)
    }
  }
}

const confirmMessage = `
Someone created a coin toss with you.
Please confirm your email address by clicking on the link below:

%s
`

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

