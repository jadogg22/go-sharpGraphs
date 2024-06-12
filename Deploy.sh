#!/bin/bash

# Define variables
CONTAINER_NAME="my-go-server"
IMAGE_NAME="my-go-server"
HOST_DATA_DIR="/home/jaden/goServe/Data"
CONTAINER_DATA_DIR="/app/Data"

# Stop and remove the existing container if it's running
echo "Stopping and removing existing container..."
docker stop $CONTAINER_NAME
docker rm $CONTAINER_NAME

# Navigate to the project directory
echo "Navigating to project directory..."
cd /home/jaden/goServe/go-sharpGraphs

# Pull the latest changes from the Git repository
echo "Pulling latest changes from Git..."
git pull

# Check if the pull was successful
if [ $? -ne 0 ]; then
        echo "Git pull failed"
        exit 1
fi

# Build the Docker image
echo "Building Docker image..."
docker build -t $IMAGE_NAME .

if [ $? -ne 0 ]; then
        echo "Image build failed"
        exit 1
fi

# Run the Docker container with the specified volume mount
echo "Running Docker container..."
docker run -d -p 5000:5000 --name $CONTAINER_NAME -v $HOST_DATA_DIR:$CONTAINER_DATA_DIR $IMAGE_NAME


#check the status of last command
if [ $? -eq 0 ]; then
        echo "Deployment completed successfully"
else
        echo " Command failed with exit code $?"
fi
