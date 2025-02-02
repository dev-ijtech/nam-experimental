FROM golang:1.23

WORKDIR /app

COPY . .

RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o /namd ./cmd/namd/*.go

EXPOSE 8080

CMD [ "/namd" ]
