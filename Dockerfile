FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go install github.com/cosmtrek/air@latest
EXPOSE 3001
CMD ["air", "-c", ".air.toml"]