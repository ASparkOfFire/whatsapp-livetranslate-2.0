FROM golang:1.24.1

WORKDIR /app

COPY . .

RUN apt-get update && apt-get install -y gcc libc6-dev python3 python3-pip ffmpeg && \
    pip3 install --break-system-packages yt-dlp
RUN go mod tidy
RUN go build -o main main.go

EXPOSE 8080

CMD ["./main"]
