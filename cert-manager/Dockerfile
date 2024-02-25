FROM golang:1.22-alpine AS build
ARG PORT
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH
ARG PASSPHRASE
RUN apk update && \
    apk add --no-cache openssl

RUN echo "port: ${PORT}"

WORKDIR /github.com/AdityaP1502/Instant-Messanging/cert_manager
COPY . .
RUN go mod download
RUN --mount=type=cache,target=/root/.cache/go-build go build -o app main.go
RUN export PASSPHRASE="$(openssl rand -base64 128)"
RUN scripts/init_key_cert.sh

FROM alpine:latest AS final

ARG PORT
ARG PASSPHRASE
ARG CERT_FILE_PATH
ARG PRIVATE_KEY_PATH

ENV PASSPHRASE ${PASSPHRASE}
ENV CERT_FILE_PATH ${CERT_FILE_PATH}
ENV PRIVATE_KEY_PATH ${PRIVATE_KEY_PATH}
ENV HOST "0.0.0.0:${PORT}"

WORKDIR /app
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/cert_manager/app .
COPY --from=build /github.com/AdityaP1502/Instant-Messanging/cert_manager/cert/ cert/
EXPOSE $PORT
CMD ["./app"]