# Log

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/joonix/log)

Formatter for logrus, allowing log entries to be recognized by the fluentd
Stackdriver agent on Google Cloud Platform.

Example:

```go
package main

import (
	"time"
	"net/http"

	log "github.com/sirupsen/logrus"
	joonix "github.com/joonix/log"
)

func main() {
	log.SetFormatter(joonix.NewFormatter())
	log.Info("hello world!")

	// log a HTTP request in your handler
	log.WithField("httpRequest", &joonix.HTTPRequest{
		Request: r,
		Status: http.StatusOK,
		ResponseSize: 31337,
		Latency: 123*time.Millisecond,
	}).Info("additional info")
}
```

## Alternatives

- https://github.com/TV4/logrus-stackdriver-formatter (seems abandoned)
- https://github.com/knq/sdhook (implemented as a hook, doesn't require fluentd)
- https://github.com/joonix/log/issues/2 (you can map the format yourself)
