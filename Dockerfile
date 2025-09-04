FROM golang:1.24.1

WORKDIR /app

COPY . .

RUN DEBIAN_FRONTEND=noninteractive apt-get update && apt-get install -y gcc libc6-dev python3 python3-pip ffmpeg wget
RUN wget https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp_linux_aarch64 -O /usr/local/bin/yt-dlp \
    && chmod a+rx /usr/local/bin/yt-dlp
RUN go mod tidy
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
