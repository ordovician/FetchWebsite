package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	// "net/http"
	// "net/url"
	"code.google.com/p/go.net/html"
	"net/url"
	"path"
	"bytes"
	"flag"
	// "net/http/cookiejar"
)

const baseURL = "https://spwebservicebm.reaktor.no"

func GetBanners(base *url.URL, client *CSRFClient) {
	// Get URL to list of banners, so we can pull out info about each banner	
	base.Path = "/admin/campaign/list"
	resp, err := client.client.Get(base.String())
	defer resp.Body.Close()
	if err != nil {
		 log.Fatal(err)		
	}
	
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		 log.Fatal(err)		
	}		
	
	err = ioutil.WriteFile("campaignlist.html", body, 0666)
	if err != nil {
		 log.Fatal(err)		
	}
	
	bufferReader := bytes.NewBuffer(body)	
	doc, err := html.Parse(bufferReader)
	if err != nil {
		 log.Fatal(err)		
	}
	
	// Download images
	for _, imgURL := range(GetBannerImgURLs(base, doc)) {
//		fmt.Println(imgURL.String())		

		resp, err := client.client.Get(imgURL.String())
		defer resp.Body.Close()
		if err != nil {
			 log.Fatal(err)		
		}	

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			 log.Fatal(err)		
		}	

		err = ioutil.WriteFile(path.Base(imgURL.Path), body, 0666)
		if err != nil {
			 log.Fatal(err)		
		}	
	}
	
	// Download linked banner info pages
	for _, bannerInfoURL := range(GetBannerInfoURLs(base, doc)) {
//		fmt.Println(imgURL.String())		

		resp, err := client.client.Get(bannerInfoURL.String())
		defer resp.Body.Close()
		if err != nil {
			 log.Fatal(err)		
		}	

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			 log.Fatal(err)		
		}	

		err = ioutil.WriteFile("banner_" + path.Base(bannerInfoURL.Path) + ".html", body, 0666)
		if err != nil {
			 log.Fatal(err)		
		}	
	}				
}

func GetBannerImgURLs(base *url.URL, doc *html.Node) []*url.URL {
	imgURLs := make([]*url.URL, 0, 10)
	
	bannerlist := FindNodeWithID("admin", doc)
	nodes := NodesInPath([]string{"div", "div", "a"}, bannerlist, false)

	// links := attributesInPath("src", []string{"div", "div", "a", "img"}, bannerlist, false)	
	// fmt.Printf("no. links %d\n", len(links))
	for _, n := range(nodes) {
		imgNode := FindTag("img", n)
		if imgNode == nil {
			continue
			// log.Fatal("Did not find tag <img> under tag <" + n.Data + ">")
		}
				
		imgURL, err := url.Parse(valueOfAttr(imgNode, "src"))

		if err != nil {
			log.Fatal(err)		
		}
						
		if imgURL.Host == "" {
			path := imgURL.Path
			*imgURL = *base 
			imgURL.Path = path
		}
		
		imgURLs = append(imgURLs, imgURL)
	}
	return imgURLs
}

func GetBannerInfoURLs(base *url.URL, doc *html.Node) []*url.URL {
	imgURLs := make([]*url.URL, 0, 10)
	
	bannerlist := FindNodeWithID("admin", doc)
	nodes := NodesInPath([]string{"div", "div", "a"}, bannerlist, false)

	// links := attributesInPath("src", []string{"div", "div", "a", "img"}, bannerlist, false)	
	// fmt.Printf("no. links %d\n", len(links))
	for _, n := range(nodes) {				
		if valueOfAttr(n, "id") == "bttn_campaign" {
			continue
		}
		
		imgURL, err := url.Parse(valueOfAttr(n, "href"))

		if err != nil {
			log.Fatal(err)		
		}
						
		if imgURL.Host == "" {
			path := imgURL.Path
			*imgURL = *base 
			imgURL.Path = path
		}
		
		imgURLs = append(imgURLs, imgURL)
	}
	return imgURLs
}

var Usage = func() {
    fmt.Fprintf(os.Stderr, "%s is a tool for downloading and storing content of a website\n\n", os.Args[0])
    fmt.Fprintf(os.Stderr, "Usage:\n")
    fmt.Fprintf(os.Stderr, "%s username password\n", os.Args[0])
}

func main() {
	flag.Parse()   // Scans the arg list and sets up flags
	if flag.NArg() != 2 {
    	fmt.Fprintf(os.Stderr, "Expected %d arguments got %d\n", 2, flag.NArg())
		Usage()
		os.Exit(0)
	}
	
	username  := flag.Arg(0)
	password  := flag.Arg(1)		
	
	client := &CSRFClient{}
	
	resp, err := client.Login(username, password)
	defer resp.Body.Close()
	if err != nil {
	     log.Fatal(err)
	}
	
	doc, err := html.Parse(resp.Body)
	if err != nil {
		 log.Fatal(err)		
	}	

	banklist := FindNodesWithClass("banklist", doc)[0]
	links := attributesInPath("href", []string{"ul", "li", "a"}, banklist, true)
	// fmt.Printf("no. links %d\n", len(links))

	baseURL := "https://spwebservicebm.reaktor.no"

	url, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, link := range(links) {
		url.Path = link
		resp, err := client.client.Get(url.String())
		defer resp.Body.Close()
		if err != nil {
			 log.Fatal(err)		
		}	
		
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			 log.Fatal(err)		
		}	

		err = ioutil.WriteFile(path.Base(link) + ".html", body, 0666)
		if err != nil {
			 log.Fatal(err)		
		}	

		// fmt.Printf("%s%s\n", baseURL, link)
	}

	GetBanners(url, client)
	
	// dumpHTTPResponse(resp)	
}



func test_main() {
	
	
	filename := "/Users/erikengheim/Development/SB1-MobilbankBM-Admin/list.html"
	file, err := os.OpenFile(filename, os.O_RDONLY, 0)

	// bytes, err := ioutil.ReadFile(filename)
	// if err != nil {
	// 	 log.Fatal(err)		
	// }	
	// fmt.Println(string(bytes))	

	base, err := url.Parse(baseURL)

	doc, err := html.Parse(file)
	if err != nil {
		 log.Fatal(err)		
	}	

	for _, imgURL := range(GetBannerImgURLs(base, doc)) {
		fmt.Println(imgURL.String())		
	}

}