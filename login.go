package main

import (
	// "fmt"
	"log"
	"net/http"
	"errors"
	"net/url"
	"code.google.com/p/go.net/html"	
)

type CSRFClient struct {
	client *http.Client
}

func (c *CSRFClient) Login(username, password string) (resp *http.Response, err error) {
	csrfToken, found := c.findCSRFTokenOnLoginPage();
	if !found {
		return nil, errors.New("Did not find CSRF token, which we need to login")
	}

	return c.postLoginForm(username, password, csrfToken)
}

func (c *CSRFClient) postLoginForm(username, password, csrfToken string) (resp *http.Response, err error) {
	values := make(url.Values)
    values.Set("signin[username]", username)
    values.Set("signin[password]", password)	                                   
	values.Set("signin[_csrf_token]", csrfToken)	
		
	return c.client.PostForm("https://spwebservicebm.reaktor.no/admin/login", values)
}

func (c *CSRFClient) findCSRFTokenOnLoginPage() (csrfToken string, found bool) {
	c.client = &http.Client{}
	jar := &myjar{}
	jar.jar = make(map[string] []*http.Cookie)
	c.client.Jar = jar	
	
	resp, err := c.client.Get("https://spwebservicebm.reaktor.no/admin")		
	defer resp.Body.Close()
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return findCSRFToken(doc)		
}

func valueOfAttr(n *html.Node, key string) string {	
	for _, attr := range(n.Attr) {
		// fmt.Printf("comparing: %s with %s\n", attr.Key, key)
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func findCSRFToken(n *html.Node) (csrfToken string, found bool) {
	if n.Type == html.ElementNode && n.Data == "input" {
		attrValue := valueOfAttr(n, "name")
		if attrValue == "signin[_csrf_token]" {
			return valueOfAttr(n, "value"), true			
		}
	}
	
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		 if csrfToken, found = findCSRFToken(c); found {
			 return 
		 }		 
	}
	
	return "", false
}
