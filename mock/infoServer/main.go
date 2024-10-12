package main

import (
	"encoding/json"
	"log"
	"math/rand/v2"
	"net/http"
)

type ResponseAudioInfo struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

const (
	song1 = `In the quiet of the evening, we’re tangled in the sheets,
The moonlight dances softly, where our heartbeats meet.
There's a story in the silence, a melody unsung,
Every breath a promise, every feeling young.

So let’s chase the shadows, where the sweet echoes roam,
In the whispers of the night, we’ve found a place called home.
With every glance, you pull me closer, like gravity unseen,
Love’s a canvas we’re painting, in a shade of evergreen.

Coffee on the counter, your laughter fills the air,
We’re weaving through the moments, with a spark that’s always there.
You’re the line in my sketchbook, the tune in my guitar,
In a world of endless changes, you’re my constant star.`
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		num := rand.IntN(3)
		resp := ResponseAudioInfo{
			ReleaseDate: "23.09.2023",
			Text:        song1,
			Link:        "https://youtu.be/dQw4w9WgXcQ",
		}
		if num == 0 {
			e := json.NewEncoder(w)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := e.Encode(resp); err != nil {
				log.Print(err)
			}
		} else if num == 1 {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
	log.Print(http.ListenAndServe(":8088", mux))

}
