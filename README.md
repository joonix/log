# Log

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/joonix/log)

Formatter for logrus, allowing log entries to be recognized by the Google Cloud Platform.
The goal is to keep concerns separate from infrastructure, services should log to
stdout/stderr and have their log automatically forwarded by fluentd.

Fluentd is the default forwarder in GCP environments. Each service shouldn't have
to know about where to send logs or what to authenticate as.

Example:

```go
package main

import (
	"os"
	"fmt"
	"flag"

	log "github.com/sirupsen/logrus"
	joonix "github.com/joonix/log"
)

func main() {
	lvl := flag.String("level", log.DebugLevel.String(), "log level")
	flag.Parse()

	level, err := log.ParseLevel(*lvl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.SetLevel(level)
	log.SetFormatter(&joonix.FluentdFormatter{})

	log.Debug("hello world!")		
}
```

## Alternatives

- https://github.com/TV4/logrus-stackdriver-formatter (seems abandoned)
- https://github.com/knq/sdhook (implemented as a hook, doesn't require fluentd)
- https://github.com/joonix/log/issues/2 (you can map the format yourself)
