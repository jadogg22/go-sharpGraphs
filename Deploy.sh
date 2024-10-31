
#!/bin/bash

# THis is just a simple script that will grab the latest changes from githb
# #After it has grabbe that latest data It stops the containers and then 
# it rebuilds the system.
# Define variables
COMPOSE_FILE="docker-compose.yml"

# Stop and remove existing containers and networks
echo "Stopping and removing existing containers and networks..."
docker-compose down

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

# Build and start the Docker containers
echo "Building and starting Docker containers..."
docker-compose up --build -d

# Check the status of the last command
if [ $? -eq 0 ]; then
    echo "Deployment completed successfully"
else
    echo "Command failed with exit code $?"
fi

