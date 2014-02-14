package main

import (
	"fmt"
	"log"
	"net/http"
	// "net/url"
	"io/ioutil"	
)

func dumpHTTPResponse(res *http.Response) {
	page, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", page)
}