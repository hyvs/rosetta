ConfigLoader (behavior)
	-> configDirectory
	<- config (global)

DomainMatcher (behavior)
	-> config, URL
	<- domainConfig

URLRewriter (behavior)
	-> domainConfig, URL
	<- URL

URLMatcher (behavior)
	-> domainConfig, URL
	<- URLConfig

Request (data)
	-> URLConfig.credentials, URL
	<- Request

Requester (behavior)
	-> Request
	<- Response

Redirecter (behavior)
	-> Response
	<- URL (or nil)

Response (data)
	-> http.Response
	<- Response

ParsingRule (data)
	-> type, pattern, ...
	<- ParsingRule

ParsingRulesMatcher (behavior)
	-> Response.ContentType, URLConfig.ParsingRules
	<- ParsingRules

Parser (behavior) (typed : json pointer, css path, ...)
	-> Response.Body, ParsingRule.Path
	<- Value OU Field

Result (data)

Field (data)

