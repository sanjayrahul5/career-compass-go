FROM golang:1.21-alpine3.18

RUN mkdir /career-compass

ADD . /career-compass

WORKDIR /career-compass

RUN go build -o app .

CMD ["/career-compass/app"]
