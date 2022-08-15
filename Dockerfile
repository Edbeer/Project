FROM golang:latest

COPY ./ ./

RUN go build -o ./cmd/api/main.go

CMD [ "./main" ]
