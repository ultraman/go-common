# go-common

```bash
go get github.com/ultraman/go-common@v0.1.1
```

## restclient
主要借鉴/抄袭/复制于client-go rest实现

```go
package main

import (
	"github.com/ultraman/go-common/rest"
)

func main() {
	config := &rest.Config{
		Host:      "https://localhost",
		UserAgent: "go-user-agent",
		ContentConfig: rest.ContentConfig{
			Codec:       rest.NewJsonMarshaler(),
			ContentType: "application/json;charset=UTF-8",
			Timeout:     10,
		},
	}

	client, err := rest.NewRESTClientFor(config)
	if err != nil {
		return nil
	}
	client.Get().Path("/api/v1/user/info").Timeout(10).Do(ctx).Raw()
}
```

logger

```go

package main

import (
	"github.com/ultraman/go-common/logger"
	"github.com/ultraman/go-common/logger/zerolog"
	"time"
)

func main() {
	logger.DefaultLogger = zerolog.NewLogger(
		logger.WithLevel(logger.InfoLevel),
		logger.WithOutputer(logger.NewOutputer("hello", "")),
		logger.WithCallerSkipCount(4),
		zerolog.WithTimeFormat(time.RFC3339))
	logger.Info("logger init success")
}
```
