# go-sharpGraphs

This is a simple revenue reporting tool that takes in reports from mcloud and generates reports from those. backend is written in go and the frontend is react and then is made static and hosted with go.

## running the project

running the project is simple make sure you have dockerinstalled and then just run the bash script 'Deploy.sh' in the root of the project. This will build the docker images and run the containers.

## Dispatchers

The dispatchers are hardcoded to be over spacific areas, if they change the dispatchers will need to be updated. in the function getLocation in pgk/database/dispatcher.go.




