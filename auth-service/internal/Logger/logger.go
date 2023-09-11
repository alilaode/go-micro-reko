package logger

import (
	"github.com/sirupsen/logrus"
	"os"
)

func Init() {
	logrus.SetOutput(os.Stdout)

	logrus.SetLevel(logrus.DebugLevel)
}
