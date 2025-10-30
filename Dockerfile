FROM golang:1.21-alpine AS build_base
RUN apk add --no-cache make 
WORKDIR /app
COPY go.mod go.sum /app/
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download 
COPY . /app/
RUN make static

FROM alpine:3
RUN apk add --no-cache graphviz
WORKDIR /app
COPY --from=build_base /app/pprof-webserver pprof-webserver
EXPOSE 8080
CMD ["/app/pprof-webserver"]