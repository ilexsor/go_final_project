# # Build stage
FROM golang:1.24 AS builder

WORKDIR /backend

COPY ./backend .

RUN go mod tidy

RUN cd ./cmd/app && CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o taskManager .

# # Release stage
FROM ubuntu:jammy

WORKDIR /go_final_project/backend/cmd/app/

COPY --from=builder /backend/cmd/app/taskManager .

COPY ./backend/.env ../../

COPY /web ../../../web

RUN chmod +x taskManager

EXPOSE 7540

CMD ["./taskManager"]