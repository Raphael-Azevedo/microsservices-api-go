# The base go-image
FROM golang:latest as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY . /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o loggerServiceApp ./cmd

RUN chmod +x /app/loggerServiceApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app
RUN mkdir /templates

COPY --from=builder /app/loggerServiceApp /app

COPY ./cmd/.env ./.env

# Run the server executable
CMD [ "/app/loggerServiceApp" ]