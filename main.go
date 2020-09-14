package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

func callODB() error {
	httpClient := &http.Client{
		Timeout: time.Duration(60) * time.Second,
	}
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		// MaxIdleConns:        100,
		// IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		DisableKeepAlives: true,
	}
	httpClient.Transport = transport

	res, err := httpClient.Get("https://localhost:8080/foo")
	if err != nil {
		return err
	}
	defer res.Body.Close()

	log.Println(res)
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello osb %s\n", r.URL.Path[1:])
	err := callODB()
	log.Println("Error: ")
	log.Println(err)
}

func main() {
	log.Println("ok")
	http.HandleFunc("/foo", handler)

	log.Fatal(http.ListenAndServe(":8082", nil))
	log.Println("yes?")
}
