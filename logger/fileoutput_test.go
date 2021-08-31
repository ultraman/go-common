package logger

import "testing"

func TestFileOutput_Write(t *testing.T) {
	out := newFileOutput("x", "./")
	out.Write([]byte("hello world"))
}
