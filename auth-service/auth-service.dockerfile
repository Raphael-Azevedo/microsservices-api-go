# The base go-image
FROM golang:latest as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY . /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o authApp ./cmd

RUN chmod +x /app/authApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/authApp /app

COPY ./cmd/.env ./.env

COPY ./migrations ./migrations

# Run the server executable
CMD [ "/app/authApp" ]



