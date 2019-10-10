## Tomus

Tomus is a library that i use to setup monitoring and observability of the buffalo apps i work with. This will setup the app i pass to use `New Relic` and `Logentries`.

### Usage

```go
//Inside your buffalo app.go
tomus.Setup(app, r)
```

### Logger

```go
import (
    "github.com/paganotoni/tomus"
)
...
tomus.Logger.Info("This is the way to use tomus logger")

```

This will:

- Adds a health check endpoint at `/admin/info`
- Use `NEWRELIC_ENV`, `ENABLE_NEWRELIC`, `NEWRELIC_LICENSE_KEY` and `APP_NAME` to add a newrelic middleware.
- Use `LOGENTRIES_TAG` to build a Logentries buffalo logger.
- Add a cross service request tracking middleware that will add the `X-Request-ID` to the context so it can be passed to other HTTP calls as header.
