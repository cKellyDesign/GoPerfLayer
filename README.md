# GoPerfLayer
The function of the Go Performance Layer is toreceive, parse, store, retrieve, and analyse performance data reported by the NDP Performance Module. The GPL acts as a router application to receive and handle incoming data to then store it an ElasticSearch database. Once we have Kibana Dashboard set up the GPL will fetch data to be displayed via Kibana in a human readable format.

### Installation

Set $GOPATH `$ export GOPATH=[path/to/repo]/GoPerfLayer`

Get Elastigo `$ go get github.com/mattbaird/elastigo`
