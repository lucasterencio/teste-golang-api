FROM golang:1.26-alpine
WORKDIR /app

COPY . .
RUN go get -d -v ./...
RUN go build -o go-api .

EXPOSE 3000

CMD [ "./go-api" ]