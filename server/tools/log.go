package tools

// ILog Log 接口
type ILog interface {
	Debug(v ...interface{})
	Info(v ...interface{})
	Error(v ...interface{}) error
	Debugf(format string, params ...interface{})
	Infof(format string, params ...interface{})
	Errorf(format string, params ...interface{}) error
}
