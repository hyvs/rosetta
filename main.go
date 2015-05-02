package main

import (
    "flag"
    "fmt"
    "regexp"
)

var url string

type UrlConverter struct {
	origin      string
	destination string
}

type ParsingRule struct {
	urlPattern   string
	ruleType     string
	contentTypes []string
	path         string
	field        string
}

type DomainConfig struct {
	converters   []UrlConverter
    parsingRules []ParsingRule
}

func init() {
	flag.StringVar(&url, "url", "", "Requested URL")
	flag.Parse()
}

func main() {
	fmt.Println("URL is ", url)

	urlConverter := UrlConverter{
		"^http(?:s)?:\\/\\/github\\.com\\/([a-zA-Z0-9\\-_]+)$",
		"https://api.github.com/users/$1",
	}
	fmt.Printf("my urlConverter: %v\n", urlConverter)
	converters := make([]UrlConverter, 0, 2)
	converters = append(converters, urlConverter)
	fmt.Printf("my urlConverters: %v\n", converters)

    rules := make([]ParsingRule, 0)
    rule := ParsingRule{
        "^https:\\/\\/api\\.github\\.com\\/users\\/",
        "json-pointer",
        []string{"application/json"},
        "/name",
        "title",
    }
    rules = append(rules, rule);

	domainConfig := DomainConfig{
		converters,
        rules,
	}
	fmt.Printf("domainConfig: %v\n", domainConfig)

    fmt.Printf("Converter for URL: %v", findConverter(url, domainConfig))
}

func findConverter(url string, config DomainConfig) *UrlConverter {
   for _, converter := range config.converters {
        matched, err := regexp.MatchString(converter.origin, url)
        if matched {
            return &converter
        }
        if err != nil {
            fmt.Printf("err %v", err)
            return nil
        }
   }

   return nil
}
