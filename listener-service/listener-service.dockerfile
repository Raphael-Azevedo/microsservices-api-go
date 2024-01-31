# The base go-image
FROM golang:latest as builder

# create a directory for the app
RUN mkdir /app

# copy all files from the current directory to the app directory
COPY . /app

# set working directory
WORKDIR /app

# build executable
RUN CGO_ENABLED=0 go build -o listenerApp ./cmd

RUN chmod +x /app/listenerApp

# create a tiny image for use
FROM alpine:latest
RUN mkdir /app

COPY --from=builder /app/listenerApp /app

# Run the server executable
CMD [ "/app/listenerApp" ]



