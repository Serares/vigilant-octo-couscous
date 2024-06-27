# Start a new stage from scratch
FROM alpine:latest  

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY ./bin/playerData .

# Command to run the executable
CMD ["playerData"]