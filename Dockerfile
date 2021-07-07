FROM golang:alpine AS build-env
RUN  apk add --no-cache git mercurial 

WORKDIR /go/src/github.com/middlewaregruppen/probear
COPY ./ .

RUN go get -d -v  ./... && \ 
CGO_ENABLED=0 GOOS=linux go build -o ./out/probear cmd/probear/main.go

FROM scratch

COPY --from=build-env /go/src/github.com/middlewaregruppen/probear/out/probear ./probear

EXPOSE 8080

CMD ["./probear"]