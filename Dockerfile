FROM golang:alpine

RUN apk update && apk add --no-cache git

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o websocket-client

EXPOSE 1322

CMD [ "/app/websocket-client" ]