# Adding Base Image for Compiling
FROM golang:1.23.2-alpine

# Setting Default Workdir
WORKDIR /app

# Setting Default Argument
ARG PUBLISH_PORT

# Copy All files into Image
COPY . .

# Run go mod tidy for install 
RUN go mod tidy

# Running for auto migration
RUN go run main.go migration

# Running for auto seeder
RUN go run main.go seeder

# Setting for build binary
RUN go build -o main

# Parsing argument into expose port
EXPOSE $PUBLISH_PORT

# Depend file binary into runner
CMD ["./main", "serve"]