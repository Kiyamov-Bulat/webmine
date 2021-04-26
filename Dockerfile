FROM alpine:latest

EXPOSE 80

RUN apk add --update --no-cache go \
		&& apk add gcc musl-dev nodejs npm docker openrc \
			&& rc-update add docker boot

# Configure Go
ENV GOROOT /usr/lib/go
ENV GOPATH /go
ENV PATH /go/bin:$PATH

RUN mkdir -p ${GOPATH}/src ${GOPATH}/bin

WORKDIR $GOPATH/src/webmine
COPY . .

WORKDIR $GOPATH/src/webmine/frontend
RUN npm ci --only=production && npm run-script build

WORKDIR $GOPATH/src/webmine
RUN go mod tidy && go build .

CMD ["./webmine"]