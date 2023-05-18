FROM golang:1.19-alpine3.16

RUN mkdir -p /app/web-studio

COPY . /app/web-studio

WORKDIR /app/web-studio

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN mkdir -p api/docs
RUN swag init -q --ot yaml -o api/docs -g cmd/main.go

RUN go build -o web-studio cmd/main.go

RUN apk add --update nodejs npm
RUN npm i -g redoc-cli
RUN redoc-cli build -o api/docs/api.html --title "API Docs" api/docs/swagger.yaml

EXPOSE 80

CMD ["./web-studio"]
