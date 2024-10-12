package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

const base = "https://dl1.vadapav.mov/8b63c507-1d93-4369-958b-fea40a6c35cd/"

func main() {

	html_string, err := SendRequests(base)

	if err != nil {
		log.Fatalln(err)
	}
	links, err := LinkGrabber(html_string)

	if err != nil {
		log.Fatalln(err)
	}
	for _, l := range links {

		cmd := exec.Command("wget", "https://dl1.vadapav.mov"+l)
		stdout, err := cmd.Output()

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		fmt.Println(string(stdout))
	}
}

func LinkGrabber(value string) ([]string, error) {
	var links []string

	doc, err := html.Parse(strings.NewReader(value))
	if err != nil {
		return nil, err
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					if strings.Contains(a.Val, "/f/") {
						links = append(links, a.Val)
					}
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
	fmt.Println(links)
	return links, nil
}

func SendRequests(value string) (string, error) {
	client := &http.Client{}

	request, err := http.NewRequest("GET", value, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}
