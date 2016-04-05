FROM golang:1.4.2-wheezy

ENV ES_PKG_NAME elasticsearch-1.5.2

WORKDIR /go/src/github.com/cKellyDesign/goTest

ADD . /go/src/github.com/cKellyDesign/goTest

# Install ElasticSearch.
RUN \
  cd / && \
  wget https://download.elastic.co/elasticsearch/elasticsearch/$ES_PKG_NAME.tar.gz && \
  tar xvzf $ES_PKG_NAME.tar.gz && \
  rm -f $ES_PKG_NAME.tar.gz && \
  mv /$ES_PKG_NAME /elasticsearch

RUN go run /go/src/github.com/cKellyDesign/goTest/test2.go

EXPOSE 8080 9200