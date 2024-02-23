FROM golang:1.22-alpine AS build
ARG SERVICE_NAME
ARG PORT
RUN echo "service_name: ${SERVICE_NAME}"
RUN echo "port: ${PORT}"

WORKDIR /github.com/AdityaP1502/Instant-Messanging/api/
COPY . .
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app service/$SERVICE_NAME/main.go

FROM alpine:latest AS final
ARG SERVICE_NAME
ARG PORT
WORKDIR /app
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/api/service/$SERVICE_NAME/config/app.config.json config/app.config.json

EXPOSE $PORT
CMD ["./app"]