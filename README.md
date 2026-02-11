## Docker Development

TL;DR: Run the following to have docker-compose also start up the Go webserver.

```sh
docker compose --profile backend up --build
```

After any changes you'll have to restart the docker-compose.

## Local Development

[mise](https://mise.jdx.dev/) can be used to get the local development tools.

If you don't use mise, check the `mise.toml` file for the required tool versions and environment variables.

1. Generate sqlc bindings (needs to be rerun any time `sql/query.sql` is updated):
   ```sh
   sqlc generate
   ```
2. Download deps
   ```sh
   go mod download
   ```
3. Run the postgres server
   ```sh
   docker compose up --build
   ```
4. Run the webserver in a separate terminal
   ```sh
   cd server
   go run .
   ```

## Debug Data

There's vibed mock-data in `scripts/mock.sql`.
By executing `scripts/mock.sh` it will **wipe all tables** and insert the mock data.
