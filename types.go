package main

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

type UrlRewriter interface {
	//	patternOrigin string
	//	patternDestination string
	Rewrite(url string) (url string)
}

type UrlMatcher interface {
	Match(c *DomainConfig, url string) *UrlConfig
}

type Requester interface {
	Get(c *Request) *http.response
}

type Redirecter interface {
	Redirect(r *http.response) (url string)
}
