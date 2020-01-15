## Tomus

Tomus is a library that i use to setup monitoring and observability of the buffalo apps i work with. This will setup the app i pass to use `New Relic` and `Logentries`.

### Usage

```go
//Inside your buffalo app.go
tomus.Setup(tomus.Config{
    App:          app,
    RenderEngine: r,

    APMKind:     tomus.APMKindDatadog,  //The type of APM to use
    ServiceName: "service",         //Service name to use in APM 
    Environment: "production",          //Environment to use in APM
    EnableAPM:   true,                  //If APM is enabled or not
})
```


This will:

- Adds a health check endpoint at `/admin/info`
- Use `LOGENTRIES_TAG` to build a Logentries buffalo logger.
- Add a cross service request tracking middleware that will add the `X-Request-ID` to the context so it can be passed to other HTTP calls as header.


### Use Logger

If you want to register something in another part of your project, you can use the tomus logger:

```go
import (
    "github.com/paganotoni/tomus"
)

...

tomus.logger.Info("This is the way to use tomus logger")
```
