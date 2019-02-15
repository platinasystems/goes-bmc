FROM golang:1.11.5-alpine
COPY . /
CMD ["go build"]

