FROM golang as builder

WORKDIR /app

COPY go.* .
COPY *.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o simpsons .

FROM golang:alpine

COPY --from=builder /app/simpsons .
COPY public ./public

CMD ./simpsons