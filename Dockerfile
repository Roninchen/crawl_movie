FROM golang
MAINTAINER yida
COPY . ./go/src/crawl_movie
EXPOSE 8080
EXPOSE 50052
VOLUME /go/src/crawl_movie/log
WORKDIR /go/src/crawl_movie/
RUN ["go build"]
CMD ["/bin/bash", "crawl_movie"]