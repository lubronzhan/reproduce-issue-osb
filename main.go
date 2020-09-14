package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func (c *Client) callODBHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello osb %s\n", r.URL.Path[1:])
	res, err := c.httpClient.Get("https://localhost:8080/foo")
	if err != nil {
		log.Println("Error: ")
		log.Println(err)
		return
	}
	defer res.Body.Close()
}

func NewClient() *Client {
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

	return &Client{
		httpClient: httpClient,
	}
}

func main() {
	log.Println("ok")
	client := NewClient()

	http.HandleFunc("/foo", client.callODBHandler)

	log.Fatal(http.ListenAndServe(":8082", nil))
	log.Println("yes?")
}
