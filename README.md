# Song API

API for managing a song library.

## Usage

1. Copy the example environment file and update it with your configuration:
    ```sh
    cp .env.example .env
    ```

2. Run the application:
    ```sh
    go run cmd/main.go
    ```

## Environment Variables

- `HTTP_PORT`: Port on which the server will run (default: `8080`)
- `SONG_DETAIL_API`: API endpoint for fetching song details
- `MODE`: Application mode (`development` or `production`)
- `DB_HOST`: Database host
- `DB_PORT`: Database port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name