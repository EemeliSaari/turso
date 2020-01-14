FROM base

RUN mkdir /parser

WORKDIR /parser

ADD . /parser

RUN go build ./main.go

CMD ["./main"]
