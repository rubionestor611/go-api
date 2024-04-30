FROM golang:1.22.2

WORKDIR /app

COPY . .

RUN go build -o main main.go

EXPOSE 80

CMD ["./main"]