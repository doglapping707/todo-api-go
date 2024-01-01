FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/go-delve/delve/cmd/dlv@latest
EXPOSE 3001
CMD ["air", "-c", ".air.toml"]