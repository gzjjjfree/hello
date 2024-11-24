package core
import (
	"fmt"
	"runtime"
)

var (
	version  = "1.0.0.0"
	build    = "Custom"
	codename = "一个学习代理转发的软件"
	intro    = "从 v2ray 借鉴学习"
)

func VersionStatement() string {
	return fmt.Sprintf("Hello %s (%s) %s (%s %s/%s)\n%s", Version(), codename, build, runtime.Version(), runtime.GOOS, runtime.GOARCH, intro)	
}

func Version() string {
	return version
}