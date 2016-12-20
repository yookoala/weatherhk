FROM golang:1.7
MAINTAINER Koala Yeung <koalay@gmail.com>

# basic environment for building
ENV GOPATH /gopath
ENV PATH ${PATH}:${GOPATH}/bin
WORKDIR ${GOPATH}/src/github.com/yookoala/weatherhk

# copy source files to build
COPY "." "./"

# build the server
RUN go install github.com/yookoala/weatherhk/cmd/weatherhk-server

# export port for public
EXPOSE 8080

ENTRYPOINT ["./scripts/docker-entrypoint.sh"]
CMD ["start"]
