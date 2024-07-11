FROM golang:alphine
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o shuttleslot
ENTRYPOINT [ "/app/shuttleslot" ]