This is intended to eventually be a general purpose tool for fetching the content of a website which requires login. At the moment it designed for fetching the content of [https://spwebservicebm.reaktor.no/admin/](https://spwebservicebm.reaktor.no/admin/).

## Why not use CURL?

To my best of knowledge you can use curl for websites which require login using a form. This tool was developed to be able to download the content of a website which required you to be authenticated by posting a form containing username and password together with a CSRF token. CSRF tokens are sent to you as a hidden field on the form by the webserver and needs to be included when logging in.

## Overview of code

The error handling and generalization of the code is limited.

| file                | purpose                                                | 
|---------------------|--------------------------------------------------------|
| cookiejar.go        | for storing HTTP cookies. Vital for CSRF to work.      |
| htmlpath.go         | pull data out of HTML document given a path of tags    |
| login.go            | login using POST form with a CSRF token                |
| fetchwebsite.go     | main and functions for storing retrieved data to disk  |


## Disclaimer

Don't use this code as an example of how to do anything in Go. It is terrible and I hope to clean it up in the future.