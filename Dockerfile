FROM golang:1.13.5

RUN go get -t github.com/EemeliSaari/turso
RUN go get -t github.com/mmcdole/gofeed
RUN go get -t github.com/streadway/amqp
