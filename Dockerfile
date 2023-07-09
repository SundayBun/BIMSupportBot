FROM golang:1.16-alpine as builder
ENV config=docker
WORKDIR /app
COPY ./ /app
RUN go mod download

FROM golang:1.20-alpine as runner
COPY --from=builder ./app ./app
ENV config=docker
WORKDIR /app

RUN go install github.com/githubnemo/CompileDaemon@latest
EXPOSE 5000
EXPOSE 7070
ENTRYPOINT CompileDaemon --build="go build main.go" --command=./main
