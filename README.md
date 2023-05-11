# LogX - Logrus Adapter

Adapter to wrap loggers from Logrus log package (https://github.com/sirupsen/logrus)

## Install

```shell
go get -u github.com/logx-go/logrus-adapter
```

## Usage

```golang
package main

import (
	"github.com/logx-go/logrus-adapter/pkg/logrusadapter"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/logx-go/contract/pkg/logx"
)

func main() {
	lr := logrus.New()
	lr.Out = os.Stdout

	logger := logrusadapter.New(lr)

	logSomething(logger)
}

func logSomething(logger logx.Logger) {
	logger.Info("This is log message")
}
```

see [examples/usage.go](examples%2Fusage.go)

## Development

### Requirement
- Golang >=1.20
- golangci-lint (https://golangci-lint.run/)

### Tests

```shell
go test ./... -race
```

### Lint

```shell
golangci-lint run
```

## License

MIT License (see [LICENSE](LICENSE) file)

