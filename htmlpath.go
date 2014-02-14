package main

import (
	"fmt"
	// "log"
	// "bufio"
	// "os"
	"code.google.com/p/go.net/html"
	// "io/ioutil"		
)
	
		
func dumpPath(strings []string)  {
	fmt.Print("path: ");
	for _, s := range(strings) {
		fmt.Printf("%s ", s)
	}
	fmt.Println();
}
	
func dumpAttributes(n *html.Node) {
	fmt.Printf("Attributes of %s: ", n.Data)
	for _, attr := range(n.Attr) {
		fmt.Printf("%s = %s,  ", attr.Key, attr.Val)
	}
	fmt.Println()
}
func NodesInPath(path []string, n *html.Node, include_siblings bool) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	var addMatching func(path []string, n *html.Node) bool
	addMatching = func(path []string, n *html.Node) bool {
		// dumpPath(path)
		if len(path) == 0 {
			return false
		}
		
		tagName := path[0]
		if n.Type != html.ElementNode || n.Data != tagName {			
			return false
		}

		if len(path) == 1 {
			nodes = append(nodes, n)
			return true
		}


		subpath := path[1:]
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if addMatching(subpath, c) && !include_siblings {
				// fmt.Printf("Not adding path: ")
		// 		dumpPath(subpath)
				break
			}
		}
		return false
	}
	addMatching(path, n)
	return nodes
}	
			
func stringsInPath(path []string, n *html.Node) []string {
	nodes := NodesInPath(path, n, true)
	strings := make([]string, 0, 10)
	for _, c := range(nodes) {
		strings = append(strings, c.Data)
	}
	return strings
}

func attributesInPath(attrName string, path []string, n *html.Node, include_siblings bool) []string {
	nodes := NodesInPath(path, n, include_siblings)
	strings := make([]string, 0, 10)
	for _, c := range(nodes) {
		// dumpAttributes(c)
		strings = append(strings, valueOfAttr(c, attrName))
	}
	return strings
}


func findFirstNodeWithAttribute(key, value string, n *html.Node) *html.Node {
	var find func(n *html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		if n.Type == html.ElementNode {
			attrValue := valueOfAttr(n, key)
			if attrValue == value {
				return n
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			 node := find(c)
			 if node != nil {
				 return node
			 }
		}
		return nil		
	}
	return find(n)
}

func findAllNodesWithAttribute(key, value string, n *html.Node) []*html.Node {
	nodes := make([]*html.Node, 0, 10)
	var collectMatching func(n *html.Node)
	collectMatching = func(n *html.Node) {
		if n.Type == html.ElementNode {
			attrValue := valueOfAttr(n, key)
			if attrValue == value {
				nodes = append(nodes, n)
			}
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			 collectMatching(c)
		}
	}
	collectMatching(n)
	return nodes
}

func FindTag(name string, n *html.Node) *html.Node {
	var find func(n *html.Node) *html.Node
	find = func(n *html.Node) *html.Node {
		// fmt.Printf("Comparing %s with %s\n", n.Data, name)
		if n.Type == html.ElementNode && n.Data == name {
			// fmt.Println("Found node!")
			return n
		}
		
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			 node := find(c)
			 if node != nil {
				 return node
			 }
		}
		return nil		
	}
	// result := find(n)
	// if result == nil {
	// 	fmt.Println("Found nothing!")		
	// }

	
	return find(n)
}

func FindNodeWithID(id string, n *html.Node) *html.Node {
	return findFirstNodeWithAttribute("id", id, n)
}

func FindNodesWithClass(class string, n *html.Node) []*html.Node {
	return findAllNodesWithAttribute("class", class, n)
}


// func main() {
// 	filename := "/Users/erikengheim/Development/SB1-MobilbankBM-Admin/nb.html"
// 	file, err := os.OpenFile(filename, os.O_RDONLY, 0)
// 	
// 	// bytes, err := ioutil.ReadFile(filename)
// 	// if err != nil {
// 	// 	 log.Fatal(err)		
// 	// }	
// 	// fmt.Println(string(bytes))	
// 	
// 	doc, err := html.Parse(file)
// 	if err != nil {
// 		 log.Fatal(err)		
// 	}	
// 
// 	banklist := FindNodesWithClass("banklist", doc)[0]
// 	links := attributesInPath("href", []string{"ul", "li", "a"}, banklist)
// 	fmt.Printf("no. links %d\n", len(links))
// 	for _, link := range(links) {
// 		fmt.Printf("%s%s\n", baseURL, link)
// 	}
// 
// }