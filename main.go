package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/hyvs/rosetta/Godeps/_workspace/src/github.com/andybalholm/cascadia"
	"github.com/hyvs/rosetta/Godeps/_workspace/src/github.com/mattn/go-jsonpointer"
	"github.com/hyvs/rosetta/Godeps/_workspace/src/golang.org/x/net/html"
)

var url string

func determineListenAddress() (string, error) {
	port := os.Getenv("PORT")
	if port == "" {
		return "", fmt.Errorf("$PORT not set")
	}
	return ":" + port, nil
}

func main() {
	addr, err := determineListenAddress()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", rosetta)
	log.Printf("Listening on %s...\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}

func rosetta(w http.ResponseWriter, r *http.Request) {
	//url := "http://www.lemonde.fr/europe/article/2015/03/26/le-point-sur-le-verrouillage-de-la-porte-du-cockpit_4602066_3214.html"
	//url := "https://github.com/hyvs"
	url := r.FormValue("url")

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer res.Body.Close()

	// Parse the HTML into nodes
	root, err := html.Parse(res.Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	selector, err := cascadia.Compile("meta[property='og:title']")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	node := selector.MatchFirst(root)
	fmt.Printf("%v\n", node)

	fmt.Println("URL is ", url)

	domainConfig := mockGithubConfig()

	if rewriter := findRewriter(url, domainConfig); rewriter != nil {
		url = rewriter.Rewrite(url)
		fmt.Println(url)
		domainConfig = mockGithubConfig()
	}

	urlConfig := buildUrlConfig(url, &domainConfig)
	fmt.Printf("%v\n", urlConfig.parsingRules)

	//@TODO: http call to URL
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	//@TODO retrieve response headers.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Printf("%s\n", body)

	// parse unknow json
	var f interface{}
	err = json.Unmarshal(body, &f)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("%T\n", f)
	//@TODO: apply parsing rules
	for _, r := range urlConfig.parsingRules {
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("rule : %v\n", r)
		field, err := jsonpointer.Get(f, r.path)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("field value : %v\n", field)
	}
	//@TODO: return OG-like structure
}

func findRewriter(url string, config DomainConfig) *UrlRewriter {
	for _, rewriter := range config.rewriters {
		matched, err := regexp.MatchString(rewriter.OriginPattern, url)
		if matched {
			return &rewriter
		}
		if err != nil {
			//@TODO: do something with err (log.error)
			fmt.Printf("err %v", err)
			return nil
		}
	}

	return nil
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

func mockGithubConfig() DomainConfig {
	urlRewriter := UrlRewriter{
		"^http(?:s)?:\\/\\/github\\.com\\/([a-zA-Z0-9\\-_]+)$",
		"https://api.github.com/users/$1",
	}
	fmt.Printf("my urlConverter: %v\n", urlRewriter)
	rewriters := make([]UrlRewriter, 0, 2)
	rewriters = append(rewriters, urlRewriter)
	fmt.Printf("my urlConverters: %v\n", rewriters)

	rules := make([]ParsingRule, 0)
	ruleName := ParsingRule{
		"^https:\\/\\/api\\.github\\.com\\/users\\/",
		"json-pointer",
		[]string{"application/json"},
		"/name",
		"title",
	}
	ruleThumb := ParsingRule{
		"^https:\\/\\/api\\.github\\.com\\/users\\/",
		"json-pointer",
		[]string{"application/json"},
		"/avatar_url",
		"thumbs",
	}
	rules = append(rules, ruleName, ruleThumb)

	domainConfig := DomainConfig{
		rewriters,
		rules,
	}
	fmt.Printf("domainConfig: %v\n", domainConfig)

	return domainConfig
}
