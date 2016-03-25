FROM golang

ENV BOT_PATH /go/src/github.com/databr/ibge-bot

RUN go get github.com/databr/bots/go_bot/parser

RUN apt-get update && apt-get install rsyslog -y && rsyslogd

COPY . ${BOT_PATH}

WORKDIR $BOT_PATH

RUN go get

CMD rsyslogd && go run ${BOT_PATH}/runner.go
