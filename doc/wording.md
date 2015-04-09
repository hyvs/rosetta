Wording
=======

### Request

Composée de l'input URL et éventuellement de crédentials. Elle va être exécutée par le Requester.

### Response

Résultat du Requester, contient divers info (content-type, http status, content-length, body, input URL)

### Message (token, request, message, query, log)

Message qui transite tout au long de la chaîne.
Ne contient pas la configuration.

### Web resource configuration

Configuration rattachée à une URL. Elle comprend :
- conversion de l'url
- les règles de parsing
- les crédentials
- et les post-processors

### Parsing rule

Peut être de type :
- json pointer
- XPATH
- CSS path
- open graph
- html meta

Contient la logique d'extraction des infos.

### Parsers

Prend une Response en entrée et applique les parsing rules correspondantes.
Renvoie une structure open graph.

### Fields

Les informations retournées en output. Respecte [le protocole OpenGraph](ogp.me)

### Url converter (router ?)

Prend en entrée une URL public (site web) et peut la traduire en une URL alternative mieux adaptée à la récupération de données. Exemples : url d'API, url oembed, ...

### Post processor

Traite les données avant de les restituer. Exemples : trim whitespaces, add/remove prefix, ...

### URLs

- input URL : url initialement fournie en entrée.
- canonical URL : url obtenue après redirection ou après conversion.
