FROM golang:1.13.0-alpine3.10 as tester
ENV CGO_ENABLED=0
RUN apk add --no-cache git openssh-client && \
    git config --global url."git@bitbucket.org:".insteadOf "https://bitbucket.org/"
RUN mkdir /root/.ssh/
COPY id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa
RUN touch /root/.ssh/known_hosts && \
    ssh-keyscan bitbucket.org >> /root/.ssh/known_hosts
WORKDIR /go/src/bitbucket.org/gismart/{{Name}}
COPY . .
RUN go get github.com/swaggo/swag/cmd/swag@v1.6.2 && swag init
RUN go get ./... && go test ./... -covermode=atomic

FROM golang:1.13.0-alpine3.10 as builder
WORKDIR /go/src/bitbucket.org/gismart/{{Name}}
RUN mkdir /go/src/bitbucket.org/gismart/{{Name}}/docs
COPY --from=tester /go/src/bitbucket.org/gismart/{{Name}}/docs /go/src/bitbucket.org/gismart/{{Name}}/docs
COPY --from=tester /go/pkg/mod /go/pkg/mod
COPY . .
RUN go build -o ./cmd/{{Name}}
RUN cd ./cmd/migration && go build

FROM alpine:3.7
WORKDIR /go/bin
RUN apk add --no-cache ca-certificates shadow && \
    groupadd {{Name}} && \
    useradd -r -g {{Name}} {{Name}}
COPY --from=builder  /go/src/bitbucket.org/gismart/{{Name}}/cmd/{{Name}} {{Name}}
COPY --from=builder  /go/src/bitbucket.org/gismart/{{Name}}/cmd/migration/migration migration
RUN chown {{Name}}:{{Name}} /go/bin/*

USER {{Name}}
ENV SRV_PORT=8888
EXPOSE 8888
CMD ["./{{Name}}"]
