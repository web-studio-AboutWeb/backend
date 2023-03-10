FROM golang:1.19-alpine3.16

RUN mkdir -p /app/app

COPY . /app/viotrina

WORKDIR /app/viotrina

RUN go build -o app cmd/main.go

EXPOSE 443

CMD ["./viotrina"]
