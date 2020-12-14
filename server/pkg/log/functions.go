package log

func Debug(args ...interface{}) {
	sugaredLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	sugaredLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	sugaredLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	sugaredLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	sugaredLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	sugaredLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	sugaredLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	sugaredLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	sugaredLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	sugaredLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	sugaredLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	sugaredLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	sugaredLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	sugaredLogger.Fatalf(template, args...)
}