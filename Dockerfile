FROM golang:1.17.3

# Install golang-migrate
# https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
RUN export PATH=$PATH:/usr/local/go/bin

