FROM alpine:latest

RUN apk add --no-cache go

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

# Install Glide
#RUN go get -u github.com/Masterminds/glide/...

EXPOSE 8080

RUN mkdir $GOPATH/src/webmine

WORKDIR $GOPATH/src/webmine

COPY webmine.go .
COPY go.mod .
RUN go build .

CMD ["./webmine"]