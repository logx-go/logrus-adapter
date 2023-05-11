package main

import (
	"fmt"
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
	logger.Info("This is a log message", "hola", "foo")
	logger.Warningf("This is a %s", "warning")
	logger.Error("This is aa error message",
		"hola", "foo",
		"error", fmt.Errorf(`foo %s`, "bar"),
	)
}
