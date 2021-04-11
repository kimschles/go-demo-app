package main

import (
	"fmt"
	"html/template"
	"time"

	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func sayHello(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("Hello World! 🌍"))
	w.WriteHeader(200)
}

type GifInfo struct {
	Title       string
	URL         string
	Description string
	Display     bool
}

type GifList struct {
	PageTitle string
	Gifs      []GifInfo
}

func showGifs(w http.ResponseWriter, req *http.Request) {
	gifTemplate := template.Must(template.ParseFiles("index.html"))
	data := GifList{
		PageTitle: "Kim's Favorite Gifs",
		Gifs: []GifInfo{
			{Title: "Community Fire", URL: "https://giphy.com/embed/ZFzpCzdWotBIdiYoLA", Description: "A scene from Community where Troy hold three pizza boxes, opens the door to an apartment ready for a party and sees flames from a fire and his friends in a chaotic situation.", Display: true},
			{Title: "Technically Correct", URL: "https://giphy.com/embed/1hMk0bfsSrG32Nhd5K", Description: "A scene from Futurama where an older man is testifying in front of a committee. He reads off a page saying 'You are technically correct, the best kind of correct.'", Display: true},
			{Title: "Evil Elmo", URL: "https://giphy.com/embed/yr7n0u3qzO9nG", Description: "A knockoff version of Elmo has its hands raised and is in front of the flames of hell.", Display: true},
		},
	}
	gifTemplate.Execute(w, data)
}

func checkHealth(w http.ResponseWriter, r *http.Request) {
	// This code is from https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-startup-probes/
	start := time.Now()
	t := time.Now()
	duration := t.Sub(start)
	if duration.Seconds() > 10 {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf("error: %v", duration.Seconds())))
	} else {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}

func main() {

	helloHandler := http.HandlerFunc(sayHello)
	http.Handle("/", helloHandler)

	gifHandler := http.HandlerFunc(showGifs)
	http.Handle("/gifs", gifHandler)

	healthHandler := http.HandlerFunc(checkHealth)
	http.Handle("/healthz", healthHandler)

	http.Handle("/metrics", promhttp.Handler())

	http.ListenAndServe(":4000", nil)
}
