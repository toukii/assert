package equal

import (
	"fmt"
	"github.com/kr/pretty"
	"io"
	"os"
	"reflect"
	"runtime"
	"unsafe"
)

func ToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func ToByte(v string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return *(*[]byte)(unsafe.Pointer(sh))
}

func assert(wr io.Writer, result bool, f func(), cd int) {
	if !result {
		_, file, line, _ := runtime.Caller(cd + 1)
		wr.Write(ToByte(fmt.Sprintf("%s:%d", file, line)))
		f()
		os.Exit(0)
	}
}

func equal(wr io.Writer, exp, got interface{}, cd int, args ...interface{}) {
	fn := func() {
		for _, desc := range pretty.Diff(exp, got) {
			wr.Write(ToByte(fmt.Sprintf("! %s", desc)))
		}
		if len(args) > 0 {
			wr.Write(ToByte(fmt.Sprintf("! %s %s", " -", fmt.Sprint(args...))))
		}
	}
	result := reflect.DeepEqual(exp, got)
	assert(wr, result, fn, cd+1)
}

func tt(wr io.Writer, result bool, cd int, args ...interface{}) {
	fn := func() {
		wr.Write(ToByte(fmt.Sprintf("!  Failure")))
		if len(args) > 0 {
			wr.Write(ToByte(fmt.Sprintf("! %s %s", " -", fmt.Sprint(args...))))
		}
	}
	assert(wr, result, fn, cd+1)
}

func T(wr io.Writer, result bool, args ...interface{}) {
	tt(wr, result, 1, args...)
}

func Tf(wr io.Writer, result bool, format string, args ...interface{}) {
	tt(wr, result, 1, fmt.Sprintf(format, args...))
}

func Equal(wr io.Writer, args ...interface{}) {
	if wr == nil {
		wr = os.Stdout
	}
	length := len(args)
	if length <= 1 {
		return
	}
	for i := 0; i < length/2; i += 1 {
		if length%2 == 1 {
			equal(wr, args[2*i], args[2*i+1], 1, args[length-1])
		} else {
			equal(wr, args[2*i], args[2*i+1], 1)
		}
	}
}

func Equalf(wr io.Writer, format string, args ...interface{}) {
	length := len(args)
	if length <= 1 {
		return
	}
	for i := 0; i < length/2; i += 1 {
		if length%2 == 1 {
			equal(wr, args[2*i], args[2*i+1], 1, fmt.Sprintf(format, args[length-1]))
		} else {
			equal(wr, args[2*i], args[2*i+1], 1)
		}
	}
}

func NotEqual(wr io.Writer, args ...interface{}) {
	length := len(args)
	if length <= 1 {
		return
	}
	for i := 0; i < length/2; i += 1 {
		fn := func() {
			wr.Write(ToByte(fmt.Sprintf("!  Unexpected: <%#v>", args[2*i])))
			if length%2 == 1 {
				wr.Write(ToByte(fmt.Sprintf("! %s %s", " -", fmt.Sprint(args[length-1]))))
			}
		}
		result := !reflect.DeepEqual(args[2*i], args[2*i+1])
		assert(wr, result, fn, 1)
	}

}

func Panic(wr io.Writer, err interface{}, fn func()) {
	defer func() {
		equal(wr, err, recover(), 3)
	}()
	fn()
}
