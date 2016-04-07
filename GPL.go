package main

import (
  "net/http" // router functionality (and so much more)
  "io/ioutil" // parsing out body of POST
  "log"
  "fmt" // general "console" and output utilities
  "encoding/json" // library to parse and unparse JSON

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
  fmt.Println("\n\nElasticSearch Connection Created")

  for report := range elastiChan {
    _, err := connection.Index(indexStr, "performance", "", nil, report)
    if err != nil {
      panic(err)
    }
    fmt.Println("\nReport Logged!")
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

// Middleware confirming only POST methods are being used
func postOnly(handle http.HandlerFunc) http.HandlerFunc {
  return func(rw http.ResponseWriter, req *http.Request) {
    if req.Method == "POST" {
      handle(rw, req)
      return
    }
    http.Error(rw, "post only", http.StatusMethodNotAllowed)
  }
}

// Routes the POST req.Body into the ElasticSearch channel to be indexed
func routeReport(rw http.ResponseWriter, req *http.Request) {
  // Read body of POST message
  body, _ := ioutil.ReadAll(req.Body)

  // Parse json string into Report struct
  var report Report
  err := json.Unmarshal(body, &report)
  if err != nil {
    panic(err)
  }

  // Middleware to manipulate data as desired
  returnStr := handleReportData(report)

  // Send stringified report to ElasticSearch via elastiChan channel
  elastiChan <- returnStr
}

// Middleware to handle different Report Types
func handleReportData(report Report) string {
  switch report.Type {
    case "envData":
      handleEnvData(report)
    case "adData":
      handleAdData(report)
    case "assetData":
      handleAssetData(report)
    case "eventLog":
      handleEventLog(report)
  }

  // Convert Report struct to json string
  returnStr, _ := json.Marshal(report)
  return string(returnStr)
}

// Handler function for Environment Data - no return
func handleEnvData (report Report) {
  fmt.Printf("\n\nReport Type: %v", report.EnvData.UserAgent)
}

// Handler function for Ad Data - no return
func handleAdData (report Report) {
  if report.AdData.HasPreroll {
    report.AdData.AdRequestDelta = report.AdData.AdRequestEnd - report.AdData.AdRequestStart
    fmt.Printf("AssetData AdRequestDelta : %v", report.AdData.AdRequestDelta)
  }
}

// Handler function for Asset Data - no return
func handleAssetData (report Report) {
  fmt.Println("ASSET DATA")
}

// Handler function for Event Logs - no return
func handleEventLog (report Report) {
  fmt.Println("EVENT LOG")
}

// Structs to parse raw data into and contain manipulated data before stringifying for Indexing
// Top Level Report Structure
type Report struct {
  Type        string            `json:"type"`
  Guid        string            `json:"guid"`
  // RawData     string            
  EnvData     *EnvDataReport    `json:"envData"`
  AdData      *AdDataReport     `json:"adData"`
  AssetData   *AssetDataReport  `json:"assetData"`
  EventLog    []*EventLogItem   `json:"eventLog"`
}

// Secondary Level Environment Data Structure
type EnvDataReport struct {
  UserAgent       string  `json:"userAgent"`
  PageURL         string  `json:"pageURL"`
  PlayerAdapter   string  `json:"playerAdapter"`
  PlayerVersion   string  `json:"playerVersion"`
  HasAdBlocker    bool    `json:"hasAdBlocker"`
  // Date string/float? `json:"date"`
}

// Secondary Level Ad Data Structure
type AdDataReport struct {
  AdRequestStart  float64   `json:"adRequestStart"`
  AdRequestEnd    float64   `json:"adRequestEnd"`
  AdRequestUrl    string  `json:"adRequestUrl"`
  HasPreroll      bool    `json:"hasPreroll"`
  PrerollData     struct {
    AdId                  string  `json:"_adId"`
    CreativeId            string  `json:"_creativeId"`
    CreativeRenditionId   string  `json:"_creativeRenditionId"`
    SlotCustomId          string  `json:"_slotCustomId"`
  }                       `json:"prerollData"`

  // any extra data from parsing middleware
  AdRequestDelta float64
}

// Secondary Level Asset Data Structure
type AssetDataReport struct {
  AssetURL    string  `json:"assetURL"`
  AssetMPXid  string  `json:"assetMPXid"`
  AssetGUID   string  `json:"assetGUID"`
  CcType      string  `json:"ccType"`
}

// Secondary Level Event Data Structure for Event Logs
type EventLogItem struct {
  Type  string  `json:"type"`
  Time  float64   `json:"time"`
  // Delta float this time - previous time
}