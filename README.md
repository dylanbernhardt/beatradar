# Beat Radar

Beat Radar is a Go-based application that fetches and manages music tracks from Beatstats. It provides a RESTful API for retrieving song information based on genre and release date.

## Features

- Fetch songs from Beatstats based on genre and date
- Caching layer using Redis for improved performance
- RESTful API with OpenAPI specification
- Kubernetes deployment ready

## Prerequisites

- Go 1.21 or later
- Redis server
- Docker (for containerization)
- Kubernetes cluster (for deployment)

## Getting Started

### Local Development

1. Clone the repository:
   ```
   git clone https://github.com/dylanbernhardt/beatradar.git
   cd beatradar
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Set up environment variables:
   ```
   export PORT=8080
   export BEATSTATS_URL=https://www.beatstats.com
   export REDIS_URL=redis://localhost:6379
   export CACHE_TTL=86400
   export SCRAPER_TIMEOUT=30
   ```

4. Run the application:
   ```
   go run cmd/beatradar/main.go
   ```

The server will start on `http://localhost:8080`.

### Running with Docker

1. Build the Docker image:
   ```
   docker build -t beatradar:latest .
   ```

2. Run the container:
   ```
   docker run -p 8080:8080 -e REDIS_URL=redis://host.docker.internal:6379 beatradar:latest
   ```

### Deploying to Kubernetes

1. Update the Docker image in `deployments/kubernetes/deployment.yaml` to point to your Docker registry.

2. Create a Kubernetes secret for the Redis URL:
   ```
   kubectl create secret generic beatradar-secrets --from-literal=redis-url=redis://your-redis-host:6379
   ```

3. Apply the Kubernetes configurations:
   ```
   kubectl apply -f deployments/kubernetes/deployment.yaml
   kubectl apply -f deployments/kubernetes/service.yaml
   ```

## API Usage

The API is documented using OpenAPI. You can find the full specification in `api/openapi.yaml`.

Example API call:

```
GET /songs?genre=House&date=2024-06-26
```

## Redis Setup

This application uses Redis for caching and requires a Redis server to be running. Follow these steps to set up Redis for use with the application:

### 1. Install Redis

If you haven't already installed Redis, you can download and install it from the [official Redis website](https://redis.io/download).

- For macOS users with Homebrew:
  ```
  brew install redis
  ```
- For Ubuntu/Debian users:
  ```
  sudo apt-get install redis-server
  ```

### 2. Start Redis Server

Once installed, start the Redis server:

- On most systems:
  ```
  redis-server
  ```
- On macOS with Homebrew:
  ```
  brew services start redis
  ```

### 3. Configure the Application

The application uses two environment variables for Redis configuration:

- `REDIS_URL`: The URL of your Redis server (default: "localhost:6379")
- `REDIS_PASSWORD`: The password for your Redis server (if required)

Set these environment variables before running the application:

```sh
export REDIS_URL=localhost:6379
export REDIS_PASSWORD=your_redis_password  # If your Redis server requires authentication
```

### 4. Verify Connection

You can verify the Redis connection by running the application. If configured correctly, you should see a log message indicating a successful connection to Redis.

### 5. Redis Configuration (Optional)

For production environments, it's recommended to configure Redis with:

- Password protection
- Proper network security settings
- Persistence options

Refer to the [Redis Security](https://redis.io/topics/security) documentation for best practices.

### Troubleshooting

If you encounter issues connecting to Redis:

1. Ensure the Redis server is running (`redis-cli ping` should return "PONG").
2. Check that the `REDIS_URL` is correct.
3. Verify that the `REDIS_PASSWORD` is correct (if using authentication).
4. Check Redis server logs for any error messages.

For more detailed Redis setup and configuration options, refer to the [official Redis documentation](https://redis.io/documentation).


## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.