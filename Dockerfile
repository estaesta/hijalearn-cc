# FROM debian:latest
FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY bin /app

# USER nonroot:nonroot

EXPOSE 8080

CMD ["./bin"]
