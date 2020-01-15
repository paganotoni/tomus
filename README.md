## Tomus

Tomus is a library that i use to setup monitoring and observability of the buffalo apps i work with. This will setup the app i pass to use `New Relic` and `Logentries`.

### Usage

```go
//Inside your buffalo app.go
config := tomus.Config{
    App: app,
}

if envy.Get("DATADOG_APM_ENABLED", "false") == "true" {
    config.APMMonitor = datadog.NewMonitor(envy.Get("APP_NAME", "app/no-name")) //you can use NewRelic here if needed.
}

Tools = tomus.New(config)
err := Tools.Start() //Woul
if err != nil {
    // Do something
}
```


This will:

- Adds a health check endpoint at `/admin/info`
- Use `LOGENTRIES_TAG` to build a Logentries buffalo logger.
- Add a cross service request tracking middleware that will add the `X-Request-ID` to the context so it can be passed to other HTTP calls as header.


### Use Logger

If you want to log something in another part of your project, you can use the tomus Wrapper Logger property. tomus.Wrapper should be initialized before using the Logger.

```go
import (
    "github.com/paganotoni/tomus"
)

var monitoringTools tomus.Wrapper

...

func Setup() {
    ...
    monitoringTools = tomus.New(config)
    ...
}


func something() {
    monitoringTools.logger.Info("This is the way to use tomus logger")
}

```
