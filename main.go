package main

import (
    "flag"
    "fmt"
)

var url string

type UrlConverter struct {
    origin string
    destination string
}

type UrlConfig struct {
    urlConverter UrlConverter
}

func init() {
    flag.StringVar(&url, "url", "", "Requested URL")
    flag.Parse()
}

func main() {
    fmt.Println("URL is ", url)
}
