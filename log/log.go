package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func Info(pkgName string, args ...interface{}) {
	logrus.Info(generateInfoArgs(pkgName, args)...)
}

func Infoln(pkgName string, args ...interface{}) {
	logrus.Infoln(generateInfoArgs(pkgName, args)...)
}

func Warn(pkgName string, args ...interface{}) {
	logrus.Warn(generateInfoArgs(pkgName, args)...)
}

func Warnln(pkgName string, args ...interface{}) {
	logrus.Warnln(generateInfoArgs(pkgName, args)...)
}

func Error(pkgName string, err error, args ...interface{}) {
	logrus.Error(generateErrorArgs(pkgName, err, args)...)
}

func Errorln(pkgName string, err error, args ...interface{}) {
	logrus.Errorln(generateErrorArgs(pkgName, err, args)...)
}

func generateInfoArgs(pkgName string, args []interface{}) []interface{} {
	var logArgs []interface{}
	logArgs = append(logArgs, fmt.Sprintf("[%s]", pkgName))
	logArgs = append(logArgs, args...)
	return logArgs
}

func generateErrorArgs(pkgName string, err error, args []interface{}) []interface{} {
	var logArgs []interface{}
	logArgs = append(logArgs, fmt.Sprintf("[%s]", pkgName))
	logArgs = append(logArgs, args...)
	logArgs = append(logArgs, "->")
	logArgs = append(logArgs, err)
	return logArgs
}
