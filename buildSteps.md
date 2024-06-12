# building 

docker build -t my-go-server .
docker run -p 5000:5000 -v goServer/Data:/app/Data my-go-server
