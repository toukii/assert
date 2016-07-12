package equal

import (
	"testing"
)

func TestLineNumbers(t *testing.T) {
	Equal(nil)
	Equal(nil, "foo")
	Equal(nil, "foo", "foo")
	Equal(nil, "foo", "foo", "msg!")
	Equal(nil, "foo", "foo", "msg!", "msg")
	Equal(nil, "foo", "foo", "msg!", "msg", "haha")
	//Equal(nil, "foo", "bar", "this should blow up")
}

func TestNotEqual(t *testing.T) {
	NotEqual(nil, "foo", "bar", "msg!")
	//NotEqual(nil, "foo", "foo", "this should blow up")
}
