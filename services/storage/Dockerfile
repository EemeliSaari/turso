FROM base

RUN mkdir /storage

WORKDIR /storage

ADD . /storage

RUN go build ./main.go

RUN rm ./main.go

CMD ["./main"]
