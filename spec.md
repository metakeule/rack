# rack bibliothek  

[HTML]($GOWORK/rack2.html)  
[Text (Markdown)]($GOWORK/rack2.md)  
[JSON]($GOWORK/rack2.json)  

Projekt `rack2`  
Firma `Know GmbH`  

Personen  

- `Marc René Arns` (metakeule)  

angefordert von

- metakeule  

## Einleitung

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `in Planung` 
| **zuletzt geändert am**          | `12.12.2013` 
| **ID**                           | `7ad64054-148f-4bb8-bdaa-86ddc37ec6d6` 


Inspiration:

- <https://github.com/metakeule/rack>
- <https://github.com/gocraft/web>
- <https://github.com/codegangsta/martini>
- <https://github.com/gorilla/mux>

Benchmarks:

<https://github.com/cypriss/golang-mux-benchmark>

******
******

## Anwendungsfälle


### S1 Integration mit fat

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `bc554f46-458d-4f2e-b5bd-d86dcfe9558c` 


- Es soll möglich sein aus fat-Structs URLs zu generieren und
zu extrahieren.

- Es soll möglich sein, REST routen automatisch aus fat Structs zu
erzeugen

- Es soll möglich sein, Abfragen mit Paging und Sorting aus fat Structs
zu erzeugen

******
******

### S2 Integration mit net/http

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `ad255354-061d-473b-9d6a-f1726dadfe0f` 


- Ein rack soll als http.Handler fungieren.
- http.Handler sollen als Middleware fungieren können

******
******

### S3 Stackability / Mountability

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `d67b99da-3840-49f5-b09b-03748f537603` 


Ein Rack soll selber als Middleware in einem anderen Rack 
fungieren können, so dass die Routen angepasst werden.

So dass ein rack wie eine unabhängige App irgendwo reingemountet werden
kann.

******
******

### S4 Performance

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `8352250c-f899-451f-b9f0-9d591b7744b9` 


rack soll mindestens so schnell sein, wie herkömmliche muxer.
Durch verzicht auf reguläre Ausdrück soll dies möglich sein.

Ein benchmark, der erweitert und geteils werden könnte ist:

<https://github.com/cypriss/golang-mux-benchmark>

Für die Performance ist ein tree-basierter Ansatz für das Routing
interessant. Hierbei kann man sich an 

<https://github.com/gocraft/web> orientieren

******
******

### S5 before / after / routing middleware

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `11.12.2013` 
| **ID**                           | `34ffae8d-e93d-4df5-8cb7-d17738db265a` 


es soll möglich sein, middleware vor und nach dem routing zu haben, sowie
spezielle routing middleware.

es soll routing middleware geben, die nach:

- schema (http/https)
- verb (GET, POST, PUT, löschen, PATCH, OPTIONS)
- host
- pfad

trennt.

einiges an middleware kann vom martini-contrib repo portiert werden:

<https://github.com/codegangsta/martini-contrib>

für die routing middleware, die im wesentlichen filter und verteiler
darstellen, kann der gorilla muxer als vorbild genommen werden:

<https://github.com/gorilla/mux>

der middleware aufruf / ansatz selbst kann von 

<https://github.com/metakeule/rack>

verwendet werden (allerdings in der richtigen reihenfolge)

******
******

### S6 context

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `783c25e3-d95f-4f35-b461-ebf81253121e` 


es soll möglich sein, context zwischen middlewares, routen usw. zu teilen.

Ideen für die Teilung von context können bei martinis service injection
genommen werden:

    Handlers are invoked via reflection. Martini makes use of Dependency 
    Injection to resolve dependencies in a Handlers argument list
    
siehe <https://github.com/codegangsta/martini>

******
******

### S7 integration mit muxer aus net/http package

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `b274e146-c0f1-4db8-ad2a-3e8a2c72d0f0` 


Es sollten standart muxer und server aus dem net/http packet integriert
werden können. ebenso soll ein rack als standart muxer verwendet 
werden können

## offene Fragen


### U1 Mögliches Design für rack

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `in Planung` 
| **zuletzt geändert am**          | `09.12.2013` 
| **ID**                           | `736e5d84-b4b1-4468-a2ab-c1f04750acb8` 


Da der net/http.Handler ohnehin ein interface ist, ist es denkbar, dass
man die "Services" à la martini einfach als Structs und die Middlewares
als Methoden setzt. Um dann an den Context zu kommen wird getypecastet.

Der Handler wird nur durchgereicht und hält alle Kontexte.

Typische Kontexte werden parallel zu Middleware in einem extra Repository
bereitgestellt. Der eigentliche Handler erbt von diesen ganzen Kontexten.

## Funktionen


### F1 Respektiere http Verben (GET/POST usw)

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `fertig` 
| **zuletzt geändert am**          | `12.12.2013` 
| **ID**                           | `85094389-f6ee-4f51-9edb-bf22d11669f2` 


zur Zeit werden nur die Pfade respektiert.

Es soll, ähnlich wie bei goh4 eine "Maske" von zweier Potenzen
gebildet werden

******
******

### F2 URLs verwalten

| **Eigenschaft**                  | **Wert**             | 
| :------------------------------- | :------------------- | 
| **Verantwortlicher**             | `metakeule` 
| **Zustand**                      | `vereinbart` 
| **zuletzt geändert am**          | `12.12.2013` 
| **ID**                           | `7063bbfb-c049-41c7-a472-f8edfa579e9b` 


Folgende Features sind wichtig:


Bei Aufruf von .Handle(), .GET() usw soll eine Route zurückgeliefert
werden.

Diese Route kann dann nach dem mounten verwendet werden, um die 
tatsächliche URL zu bekommen.

Außerdem soll ein Router alle seine URLs zurückgeben können.
Außerdem soll eine URL relativ zu einem Router abgefragt werden können.
Außerdem brauchen wir die Möglichkeit eine URL mit Parametern zu
füllen (für die Platzhalter).

Außerdem brauchen wir die Möglichkeit, die Parameter einer aktuellen
URL in einem Handler zu bekommen.

Außerdem brauchen wir die Möglichkeit, structs zu verwenden (mit tags),
die zum konstruieren einer URL befüllt werden können und in die 
Parameter einer URL eingelesen werden können.

******
******
******
******


