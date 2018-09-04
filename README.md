# GOV.UK PaaS LogCacheAdapter

Application for collecting metrics from GOV.UK PaaS via log-cache in the following formats:

- Prometheus text format

## How to run

This application requires the `log-cache` API URL to be provided.
It can be passed in as `LOG_CACHE_API` environment variable or as
`-log-cache-api` / `-a` flag.

```
make run
# or
go run main.go handlers.go middleware.go server.go utils.go # make sure to list all files
```

### Available flags

| Environment Variable | Flag | Description | Default |
|---|---|---|---|
| `LOG_CACHE_API` | `-log-cache-api`, `-a` | Required. The log-cache API URL. | `N/A` |
| `PORT` | `-port`, `-p` | The port server should be running on. | `"8080"` |
| `DEBUG` | `-verbose`, `-v` | Run the server in debugging mode. | `"false"`

## Testing

To run test, simply execute the following command:

```sh
make test
```

If you run into problems with go test caching, set the enviornment variable
`GOCACHE=off`.

## API

The application exposes single endpoint, `/metrics` - root. It requires two
headers to be provided, in order for it to work as expected. These are:

### `Accept` header

Which tells the application at what format would you like to receive the logs.
Currently only one formats is recognised:

- Prometheus: `text/plain`

### `Authorization` header

In format of:

```
bearer ${JWT_WEB_TOKEN}
```

Which tells the application, what metrics should you have access to.

It can work to your advantage, if you create an user with restricted access, to
target only the applications and instances you want.

For instance, you can obtain that token with the use of CLI command:

```sh
cf oauth-token
```

### Example

If you put everything together and execute the following curl command:

```sh
curl -H "Accept: text/plain" -H "Authorization: $(cf oauth-token)" http://localhost:8080/
```

You should retreive the Prometheus text of different metrics your user has
access to.

