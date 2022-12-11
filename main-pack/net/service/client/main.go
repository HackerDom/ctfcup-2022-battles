package main

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const timeout = 5 * time.Second

var flag string

func getMetrics(byteSlice []byte) map[string]float64 {
	metrics := map[string]float64{}
	s := string(byteSlice)
	paramsList := strings.Fields(s)
	for i := 0; i < len(paramsList); i += 2 {
		metrics[paramsList[i]], _ = strconv.ParseFloat(paramsList[i+1], 64)
	}

	return metrics
}

func check(url string) bool {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	httpClient := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return false
	}

	req.Header.Set("Authorization", fmt.Sprintf("%x", md5.Sum([]byte(flag))))

	res, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}

	if res.StatusCode != 200 {
		log.Println(fmt.Sprintf("Incorrect status code. Got %s, wait 200", res.Status))
		return false
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	metrics := getMetrics(body)
	if _, ok := metrics[flag]; !ok {
		log.Println("Can't get flag from metrics")
		return false
	}
	log.Println("Success")
	return true
}

func main() {
	url, exist := os.LookupEnv("URL")
	if !exist {
		log.Fatal("URL not found")
	}

	flag, exist = os.LookupEnv("FLAG")
	if !exist {
		log.Fatal("FLAG not found")
	}

	for {
		if !check(url) {
			log.Println("ERROR: check isn't correct")
		}
		time.Sleep(timeout)
	}
}
