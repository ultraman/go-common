package logger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

var (
	TimeFormatLayout = "2006_01_02"
	DefaultOutputer  = newConsoleOutput()
)

func NewOutputer(name, path string) Outputer {
	return newFileOutput(name, path)
}

type Outputer interface {
	io.WriteCloser
}

type consoleOutput struct {
	w io.WriteCloser
}

func (c *consoleOutput) Write(p []byte) (n int, err error) {
	return c.w.Write(p)
}

func (c *consoleOutput) Close() error {
	return c.w.Close()
}

var _ Outputer = &consoleOutput{}

func newConsoleOutput() Outputer {
	return &consoleOutput{
		w: os.Stdout,
	}
}

type fileOutput struct {
	name string
	path string
	file *os.File
	lock sync.Mutex
}

var _ Outputer = &fileOutput{}

func newFileOutput(name string, path string) Outputer {
	fo := &fileOutput{name: name, path: path}
	f := fo.createFile()
	if f == nil {
		return nil
	}
	fo.file = f
	go fo.check()
	return fo
}

func (f *fileOutput) Write(p []byte) (n int, err error) {
	lfile := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&f.file)))
	ff := (*os.File)(lfile)
	if ff != nil {
		f.lock.Lock()
		if runtime.GOOS == "windows" {
			n, err = ff.Write([]byte("\r\n"))
		}
		n, err = ff.Write(p)
		f.lock.Unlock()
		return n, err
	}
	return 0, err
}

func (f *fileOutput) Close() error {
	lfile := atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&f.file)), nil)
	ff := (*os.File)(lfile)
	if ff != nil {
		_ = ff.Sync()
		_ = ff.Close()
	}
	return nil
}

func (f *fileOutput) createFile() *os.File {
	prefix := f.name
	now := time.Now()
	format := now.Format(TimeFormatLayout)
	name := fmt.Sprintf("%s_%s.log", prefix, format)
	name = f.path + name
	var err error
	var file *os.File

	file, err = os.OpenFile(name, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil
	}
	return file
}

func (f *fileOutput) fileIsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

func (f *fileOutput) check() {
	for {
		nextDay := time.Now().Add(time.Hour * 24)
		nextDay = time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), 0, 0, 0, 0, nextDay.Location())
		dayTimer := time.NewTimer(nextDay.Sub(time.Now()))
		secTimer := time.NewTimer(time.Second * 2)
		select {
		case <-dayTimer.C:
			newFile := f.createFile()
			if newFile != nil {
				oldFile := f.file
				atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&f.file)), unsafe.Pointer(newFile))
				_ = oldFile.Sync()
				_ = oldFile.Close()
			}
		case <-secTimer.C:

		}
	}
}
