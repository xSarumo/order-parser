FROM golang
WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN go build -o ./out/server ./cmd/main.go

EXPOSE 8081
CMD [ "./out/server" ]
