# Start from golang base image
FROM golang:1.15-alpine as builder
 
# Set maintainer label: maintainer=[YOUR-EMAIL]
LABEL maintainer='andreas.roither@protonmail.com'

# install git for dependencies
RUN apk add --no-cache git

ENV CGO_ENABLED=0

# Set working directory: `/app`
WORKDIR /app/
 
# Copy files to the working directory
COPY . .

 # List items in the working directory (ls)
RUN ls -l
 
# Build the GO app as binary
RUN go build -o .

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /app/

# Copy the Pre-built binary file from the previous stage.
COPY --from=builder /app/snippets .

EXPOSE 8010
 
# Run the service myapp when a container of this image is launched
ENTRYPOINT ["./snippets"]