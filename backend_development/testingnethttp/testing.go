package main

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	// // https://go.dev/play/p/eE32qPmuDeS
	// const method = "GET"
	// const url = "https://eblog.fly.dev/index.html"
	// var body io.Reader = nil
	// req, err := http.NewRequestWithContext(context.TODO(), method, url, body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// req.Header.Add("Accept-Encoding", "gzip")
	// req.Header.Add("Accept-Encoding", "deflate")
	// req.Header.Set("User-Agent", "eblog/1.0")
	// req.Header.Set("some-key", "a value")   // will be canonicalized to Some-Key
	// req.Header.Set("SOMe-KEY", "somevalue") // will overwrite the above since we used Set rather than Add
	// req.Write(os.Stdout)
	//
	const method = "GET"
	v := make(url.Values)
	v.Add("q", `"of Emrakul"`)
	v.Add("order", "released")
	v.Add("dir", "asc")
	const path = "https://scryfall.com/search"
	dst := path + "?" + v.Encode() // Encode() will escape the values for us. Remember the '?' seperator
	req, err := http.NewRequestWithContext(context.TODO(), method, dst, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Write(os.Stdout)
}
