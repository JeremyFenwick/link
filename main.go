package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	htmlFlag := flag.String("html", "N/A", "Enter the name opf the html document.")
	flag.Parse()
	if *htmlFlag == "N/A" {
		fmt.Println("Please enter a file name. Exiting...")
		os.Exit(0)
	}
	data, err := os.ReadFile(*htmlFlag)
	if err != nil {
		fmt.Println("Could not open file. Exited with code 1")
		os.Exit(1)
	}
	r := strings.NewReader(string(data))
	links, err := Parse(r)
	if err != nil {
		fmt.Printf("Something went wrong with the parser %v", err)
	}
	for _, link := range links {
		fmt.Printf("%v\n", link)

	}
}

// Link represents a link in a HTML document (<a href="...")
type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	var links []Link
	nodes := LinkNodes(doc)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func buildLink(n *html.Node) Link {
	var ret Link
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			ret.Href = attr.Val
		}
	}
	ret.Text = text(n)
	return ret

}

func text(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	if n.Type != html.ElementNode {
		return ""
	}
	var ret string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret += text(c) + " "
	}
	return strings.Join(strings.Fields(ret), " ")
}

func LinkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var ret []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		ret = append(ret, LinkNodes(c)...)
	}
	return ret
}

/*
var exampleHtml = `
<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
</body>
</html>
`*/
