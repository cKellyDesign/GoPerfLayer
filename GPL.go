package main

import (
  "net/http" // router functionality (and so much more)
  "io/ioutil" // parsing out body of POST
  "log"
  "fmt" // general "console" and output utilities

  elastigo "github.com/mattbaird/elastigo/lib" // ElasticSearch Go Library
)

var elastiChan chan string
var indexStr = "reports"

func main() {
  go createElastic() // init ElasticSearch listener in seperate GO thread
  startWebServer() // init Web Server / Router
}

// Starts listenting to elastiChan chanel for reports to Index
func createElastic() {
  elastiChan = make(chan string, 8)
  connection := elastigo.NewConn()
  connection.Domain = "localhost"
  fmt.Printf("\n\nElasticSearch Connection Created at Domain : %s", connection.Domain)

  for report := range elastiChan {
    _, err := connection.Index(indexStr, "performance", "", nil, report)
    if err != nil {
      panic(err)
    }
  }
}

// Starts listenting to various routes on definied port
func startWebServer() {
  http.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request){
    fmt.Fprint(rw, "Hello World!")
  })
  http.HandleFunc("/perfReport", postOnly(routeReport))
  fmt.Printf("Starting Server on port %v", 8080)
  log.Fatal(http.ListenAndServe(":8080", nil))
}

// Routes the POST req.Body into the ElasticSearch channel to be indexed
func routeReport(rw http.ResponseWriter, req *http.Request) {
  body, _ := ioutil.ReadAll(req.Body)
  // fmt.Printf("\n\nRequest Body: %v", string(body))
  elastiChan <- string(body) // send stringified report to ElasticSearch via elastiChan channel
}

// confirm only POST methods are being used
func postOnly(handle http.HandlerFunc) http.HandlerFunc {
  return func(rw http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
      handle(rw, req)
      return
    }
    http.Error(rw, "post only", http.StatusMethodNotAllowed)
  }
}