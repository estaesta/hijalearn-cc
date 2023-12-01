FROM debian:latest

WORKDIR /app
COPY testbuild /app

CMD ["./testbuild"]
