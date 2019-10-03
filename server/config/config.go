package config

import (
	"flag"
	"os"
)

const (
	// InternalError  服务器内部错误
	InternalError = iota + 50000
	// IllegalArgument 参数错误
	IllegalArgument
	// DatabaseError 数据库错误
	DatabaseError
	// HTTPError HTTP请求错误
	HTTPError
	// ExternalError 接口请求错误
	ExternalError
	// UnmarshalTypeFail 结构解析错误
	UnmarshalTypeFail
	// MarshalTypeFail 结构转换错误
	MarshalTypeFail
	// TypeTranslateFail 类型转化错误
	TypeTranslateFail
	// ValueNotExist 值不存在
	ValueNotExist
)

// Config is a configuration interface
type Config interface {
	IsSet(name string) bool
	Bool(name string) bool
	Int(name string) int
	IntSlice(name string) []int
	Int64(name string) int64
	Int64Slice(name string) []int64
	String(name string) string
	StringSlice(name string) []string
	Uint(name string) uint
	Uint64(name string) uint64
	Set(name, value string) error
}

var (
	initializers []func(Config)
	config       Config
	Auth         = os.Getenv("AUTH") == "true"
	Online       = flag.Bool("online", false, "online flag")
	SkipConsul   = flag.Bool("skip_consul", false, "Skip to connect to consul")
)

func IsSet(name string) bool           { return config.IsSet(name) }
func Bool(name string) bool            { return config.Bool(name) }
func Int(name string) int              { return config.Int(name) }
func IntSlice(name string) []int       { return config.IntSlice(name) }
func Int64(name string) int64          { return config.Int64(name) }
func Int64Slice(name string) []int64   { return config.Int64Slice(name) }
func String(name string) string        { return config.String(name) }
func StringSlice(name string) []string { return config.StringSlice(name) }
func Uint(name string) uint            { return config.Uint(name) }
func Uint64(name string) uint64        { return config.Uint64(name) }
func Set(name, value string) error     { return config.Set(name, value) }

// AddInitializer Add a initializer, call on initialized
func AddInitializer(fc func(Config)) {
	initializers = append(initializers, fc)
}

// Initialize initialize process configure
func Initialize(c Config) {
	config = c
	InitDatabase()
	for _, initFunc := range initializers {
		initFunc(c)
	}
}
