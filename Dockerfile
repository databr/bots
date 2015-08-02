FROM golang

ENV BOT_PATH /go/src/github.com/databr/metrosp-bot

ADD . ${BOT_PATH}

RUN cd ${BOT_PATH}

CMD go run ${BOT_PATH}/runner.go
