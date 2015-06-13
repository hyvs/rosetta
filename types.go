package main

import "net/http"

type Config []DomainConfig

// config for domain (converters, parsing rules, credentials, etc.)
type DomainConfig struct {
	rewriters    []UrlRewriter
	parsingRules []ParsingRule
}

// config for input URL ; parsing rules here are a filtered subset of domain's parsing rules
type UrlConfig struct {
	parsingRules []ParsingRule
}

type Request struct {
	credentials string
	url         string
}

type ParsingRule struct {
	urlPattern   string
	ruleType     string
	contentTypes []string
	path         string
	field        string
}

type ConfigLoader interface {
	Load(dirname string) *Config
}

type DomainMatcher interface {
	Match(c *Config, url string) *DomainConfig
}

type UrlRewriter struct {
	OriginPattern      string
	DestinationPattern string
}

func (UrlRewriter) Rewrite(url string) string {
	return "mock-url"
}

type UrlMatcher interface {
	Match(c *DomainConfig, url string) *UrlConfig
}

type Requester interface {
	Get(c *Request) *http.Response
}

type Redirecter interface {
	Redirect(r *http.Response) (url string)
}
