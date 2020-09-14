# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder
# Add Maintainer Info
LABEL maintainer="Suthisak Chuenjit <suthisak.ch@gmail.com>"
ENV TZ Asia/Bangkok
# Set the Current Working Directory inside the container
WORKDIR /go/src/runex
# Build Args
ARG LOG_DIR=/go/src/runex/logs
# Create Log Directory
RUN mkdir -p ${LOG_DIR}

# Environment Variables
ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log 

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
######## Start a new stage from scratch #######
FROM alpine:latest
#RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
RUN apk --no-cache add ca-certificates --no-cache tzdata\
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" >  /etc/timezone

WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/runex/main .
COPY --from=builder /go/src/runex/runex.co.crt .
COPY --from=builder /go/src/runex/runex.co.key .
COPY --from=0 /go/src/runex/templates ./templates

# Copy config file
COPY config-prd.yaml ./config.yaml

# Declare volumes to mount
# VOLUME $PWD/runex:/go/src/runex/upload
# VOLUME $PWD/runex:/go/src/runex/template

# Expose port 8080 to the outside world
EXPOSE 3006

# Command to run the executable
CMD ["./main"]