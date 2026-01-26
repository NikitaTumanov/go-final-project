FROM ubuntu:latest

RUN apt-get update
RUN apt-get install -y wget git gcc

ENV GO_VERSION=1.25.6
RUN wget -P /tmp "https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"
RUN rm "/tmp/go${GO_VERSION}.linux-amd64.tar.gz"

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY internal internal
COPY web web
COPY main.go .

ENV TODO_PORT=8080
ENV TODO_PASSWORD=12345
ENV TODO_DBFILE=/app/data/scheduler.db

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o application main.go

CMD [ "./application" ]

EXPOSE 8080
