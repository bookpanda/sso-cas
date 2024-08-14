FROM golang:1.22.4-alpine3.20 as builder
WORKDIR /app

COPY ./backend/go.mod ./backend/go.sum ./

RUN go mod download

COPY ./backend/. .

RUN go build -o server ./cmd/main.go


FROM alpine AS runner
WORKDIR /app

COPY --from=builder /app/server ./

ENV GO_ENV production

EXPOSE 3000

CMD ["./server"]