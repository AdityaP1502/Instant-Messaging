FROM golang:1.22-alpine AS build
ARG PORT
RUN echo "port: ${PORT}"

WORKDIR /github.com/AdityaP1502/Instant-Messanging/reverse_proxy
COPY . .
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app main.go

FROM alpine:latest AS final
ARG SERVICE_NAME
ARG PORT
WORKDIR /app
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/reverse_proxy/config/app.config.json config/app.config.json

EXPOSE $PORT
CMD ["./app"]