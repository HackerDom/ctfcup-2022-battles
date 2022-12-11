package main

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var flag string

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandStr() string {
	n := 5
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func GetRandMetric() string {
	return fmt.Sprintf("metric_%s %f", GetRandStr(), rand.Float32())
}

func GetMetrics() string {
	var metrics []string
	for i := 0; i < rand.Intn(10); i++ {
		metrics = append(metrics, GetRandMetric())
	}
	metrics = append(metrics, fmt.Sprintf("%s %f", flag, rand.Float64()))

	rand.Shuffle(len(metrics), func(i, j int) { metrics[i], metrics[j] = metrics[j], metrics[i] })

	return strings.Join(metrics, "\n")
}

func getMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method is not supported.", http.StatusNotFound)
		return
	}

	if strings.Split(r.RemoteAddr, ":")[0] != "10.13.37.13" {
		http.Error(w, "Incorrect IP", http.StatusUnauthorized)
		return
	}

	if r.Header.Get("Authorization") != fmt.Sprintf("%x", md5.Sum([]byte(flag))) {
		http.Error(w, "User unauthorized", http.StatusUnauthorized)
		return
	}

	_, err := fmt.Fprint(w, GetMetrics())
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
}

func main() {
	port, exist := os.LookupEnv("PORT")
	if !exist {
		log.Fatal("PORT not found")
	}

	flag, exist = os.LookupEnv("FLAG")
	if !exist {
		log.Fatal("FLAG not found")
	}

	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/metrics", getMetrics)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Metrics Exporter</title></head>
			<body>
			<h1>Metrics Exporter</h1>
			<p><a href="/metrics">Metrics</a></p>
			</body>
			</html>`))
	})

	fmt.Printf("Starting server at port :%s\n", port)
	if err := http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil); err != nil {
		log.Fatal(err)
	}
}
