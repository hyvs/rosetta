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

type UrlConfig struct {
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

    if converter := findConverter(url, domainConfig); converter != nil {
        url = convertUrl(url, converter)
        fmt.Println(url)
        //@TODO: reload new DomainConfig for url
    }

    urlConfig := buildUrlConfig(url, &domainConfig)
    fmt.Printf("%v\n", urlConfig.parsingRules)
}

// @TODO: a method of the UrlConverter struct ?
func findConverter(url string, config DomainConfig) *UrlConverter {
   for _, converter := range config.converters {
        matched, err := regexp.MatchString(converter.origin, url)
        if matched {
            return &converter
        }
        if err != nil {
            //@TODO: do something with err (log.error)
            fmt.Printf("err %v", err)
            return nil
        }
   }

   return nil
}

// @TODO: a method of the UrlConverter struct ?
func convertUrl(url string, converter *UrlConverter) string {
    re := regexp.MustCompile(converter.origin);
    return re.ReplaceAllString(url, converter.destination)
}

func buildUrlConfig(url string, domainConfig *DomainConfig) *UrlConfig {
    rules := make([]ParsingRule, 0)
    for _, rule := range domainConfig.parsingRules {
        fmt.Printf("%s\n", rule.urlPattern)
        matched, err := regexp.MatchString(rule.urlPattern, url)
        if matched {
            rules = append(rules, rule)
        }
        if err != nil {
            //@TODO: do something with err (log.error)
            fmt.Printf("err %v", err)
            return nil
        }
    }

    urlConfig := &UrlConfig{rules}
    return urlConfig
}

