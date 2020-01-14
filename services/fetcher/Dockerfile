FROM base

RUN mkdir /fetcher

WORKDIR /fetcher

ADD . /fetcher

RUN go build ./main.go

CMD ["./main"]
