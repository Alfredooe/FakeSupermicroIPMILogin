FROM golang:1.23

WORKDIR /app

COPY ./app /app

RUN go build -v -o main .

CMD ["./main"]