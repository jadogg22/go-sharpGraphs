# stage 1  Build vite frontend 
FROM node:18 AS frontend-builder

WORKDIR /app/frontend

# copy the frontend source cod to to working directory
COPY frontend/ package*.json ./
Copy frontend/ . 

# Install dependencies and build the frontend

RUN npm install
RUN npm run build

# Stage 2: Build Go server
FROM golang:1.20 AS go-builder

WORKDIR /app

# Copy go source code to working directory

Copy go.mod go.sum ./
RUN go mod download

COPY . .

# Copy the built frontend from the previous stage
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Build the go aplication
RUN go build -o /app/server .

# Stage 3: Final image
FROM golang:1.20

WORKDIR /app

# Copy the built go executable from the previous stage
COPY --from=go-builder /app/server .
COPY --from=go-builder /app/frontend/dist ./frontend/dist

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./app/server"]


