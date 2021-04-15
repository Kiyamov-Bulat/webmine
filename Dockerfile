FROM alpine:latest

EXPOSE 80

RUN apk add --update --no-cache go
RUN apk add gcc musl-dev
RUN apk add nodejs npm
RUN apk add docker openrc
RUN rc-update add docker boot

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR $GOPATH/src/webmine
COPY . .

WORKDIR $GOPATH/src/webmine/frontend
RUN npm install 
RUN npm build .

WORKDIR $GOPATH/src/webmine
RUN go mod tidy
RUN go build .

CMD ["./webmine"]