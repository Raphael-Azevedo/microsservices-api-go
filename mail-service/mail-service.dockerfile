# The base go-image
FROM golang:latest as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY . /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o mailServiceApp ./cmd

RUN chmod +x /app/mailServiceApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY templates /templates

COPY --from=builder /app/mailServiceApp /app

# Run the server executable
CMD [ "/app/mailServiceApp" ]