FROM golang:alpine
WORKDIR /build
COPY . .
RUN go build -o secret-santa main.go
CMD ["./secret-santa"]