FROM golang:alpine
WORKDIR /build
COPY . .

# Install goose for db migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN go build -o secret-santa main.go
CMD ["sh", "-c", "goose -dir /build/db/migrations postgres $DATABASE_URL up && ./secret-santa"]