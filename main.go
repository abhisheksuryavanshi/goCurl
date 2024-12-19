/*
Example usage:

	./gocurl http://eu.httpbin.org/get
	./gocurl http://eu.httpbin.org/get:80
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var supportedProtocols = []string{"http"}

func isProtocolSupported(protocolName string) bool {
	for _, v := range supportedProtocols {
		if v == protocolName {
			return true
		}
	}
	return false
}

// Finds url from the command line args and returns parsed url.
// If url is not found returns nil
func getUrl() *url.URL {
	// url must be provided.
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./goCurl <URL>")
		os.Exit(0)
	}
	// url must be the last argument.
	rawUrl := os.Args[len(os.Args)-1]
	if !strings.Contains(rawUrl, "://") {
		rawUrl = fmt.Sprintf("https://%s", rawUrl)
	}
	// try to parse the url.
	parsedUrl, err := url.ParseRequestURI(rawUrl)
	if err != nil {
		println("Provided url is malformed: ", err)
	}
	// if port wasn't specified, append port according to the scheme.
	if len(parsedUrl.Port()) == 0 {
		switch parsedUrl.Scheme {
		case "http":
			parsedUrl.Host = parsedUrl.Host + ":80"
		case "https":
			parsedUrl.Host = parsedUrl.Host + ":443"
		}
	}
	return parsedUrl
}

func makeGetRequest(userUrl *url.URL) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", userUrl.String(), nil)
	req.Header.Set("Connection", "close")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	// Print all the headers we sent.
	for key, values := range req.Header {
		for _, value := range values {
			fmt.Printf("> %s: %s\n", key, value)
		}
	}

	// Print all the headers we recieved.
	for key, values := range response.Header {
		for _, value := range values {
			fmt.Printf("< %s: %s\n", key, value)
		}
	}
	content, _ := io.ReadAll(response.Body)
	// Response that we got.
	println(string(content))
}

func main() {
	fmt.Println("Welcome to goCurl")
	makeGetRequest(getUrl())
}
