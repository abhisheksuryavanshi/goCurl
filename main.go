/*
Example usage:

	gocurl http://eu.httpbin.org/get
	gocurl http://eu.httpbin.org/get:80
*/
package main

import (
	"fmt"
	"io"
	"net/http"
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

func getRequestParams(userURL string) []string {
	userParams := []string{}
	splitURL := strings.Split(userURL, "://")
	// Protocol at index 0
	userParams = append(userParams, splitURL[0])
	var userHost string
	for i := 0; i < len(splitURL[1]); i++ {
		if splitURL[1][i] == '/' || splitURL[1][i] == ':' {
			break
		}
		userHost += string(splitURL[1][i])
	}
	// Host at index 1
	userParams = append(userParams, userHost)
	var port string
	var httpMethod string
	for i := len(userURL) - 1; i >= 0; i-- {
		if len(port) == 0 && userURL[i] == ':' {
			var customPort string
			for j := i + 1; j < len(userURL) && userURL[j] != '/'; j++ {
				customPort += string(userURL[j])
			}
			port = customPort
		}
		if len(httpMethod) == 0 && userURL[i] == '/' {
			var customHttpMethod string
			for j := i + 1; j < len(userURL); j++ {
				customHttpMethod += string(userURL[j])
			}
			httpMethod = customHttpMethod
		}
	}
	if len(port) == 0 {
		port = "80"
	}

	if len(httpMethod) == 0 {
		httpMethod = "get"
	}
	// Port at index 2
	userParams = append(userParams, port)
	// HTTP method at index 3
	userParams = append(userParams, httpMethod)

	return userParams
}

func verifyURL(userURL string) bool {
	splitURL := strings.Split(userURL, "://")
	if len(splitURL) < 2 {
		fmt.Println("URL malformed.")
		return false
	}
	protocolName := splitURL[0]
	if !isProtocolSupported(protocolName) {
		fmt.Println(protocolName, "not supported.")
		return false
	}
	return true
}

func makeGetRequest(userUrl string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", userUrl, nil)
	req.Header.Set("Connection", "close")
	response, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	// Print all the headers we sent.
	for key, values := range response.Header {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}
	content, _ := io.ReadAll(response.Body)
	// Response that we got.
	println(string(content))
}

func main() {
	fmt.Println("Welcome to goCurl")
	commandLineArgs := os.Args[1:]
	fmt.Println(commandLineArgs)

	userUrl := commandLineArgs[len(commandLineArgs)-1]
	if !verifyURL(userUrl) {
		panic("exiting")
	}
	cliargs := getRequestParams(userUrl)
	fmt.Println(cliargs)
	makeGetRequest(userUrl)
}
