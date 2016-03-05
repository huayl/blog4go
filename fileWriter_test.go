// Copyright (c) 2015, huangjunwei <huangjunwei@youmi.net>. All rights reserved.

package blog4go

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"testing"
)

// test if log lose in multi goroutine mode
func TestFileWriterMultiGoroutine(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer blog.Close()
	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}

	var wg sync.WaitGroup
	var beginWg sync.WaitGroup

	f := func() {
		defer wg.Done()
		beginWg.Wait()
		for i := 0; i < 100; i++ {
			blog.Infof("haha %s. en\\en, always %d and %f, %t, %+v", "eddie", 18, 3.1415, true, t)
		}
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		beginWg.Add(1)
	}

	// write 100 * 100 lines to /tmp/info.log
	for i := 0; i < 100; i++ {
		go f()
		beginWg.Done()
	}

	wg.Wait()
	blog.flush()

	out, err := exec.Command("/bin/sh", "-c", "/usr/bin/wc -l /tmp/info.log").Output()
	if nil != err {
		t.Errorf("count file lines failed. err: %s", err.Error())
	}

	arr := strings.Split(string(out), " ")
	intStr := arr[len(arr)-2]
	lines, err := strconv.Atoi(intStr)
	if nil != err {
		t.Errorf("line string convert to int failed. err: %s", err.Error())
	}

	if 100*100 != lines {
		t.Errorf("it loses %d lines.", 100*100-lines)
	}

	// clean logs
	_, err = exec.Command("/bin/sh", "-c", "/bin/rm /tmp/info.log").Output()
	if nil != err {
		t.Errorf("clean files failed. err: %s", err.Error())
	}
}

func TestSizeAndLinesBaseLogrotate(t *testing.T) {
	err := NewFileWriter("/tmp", false)
	defer blog.Close()
	if nil != err {
		t.Errorf("initialize file writer faied. err: %s", err.Error())
	}
	blog.SetRotateSize(1)

}

func BenchmarkFileWriters(b *testing.B) {
	err := NewFileWriter("/tmp", false)
	defer blog.Close()
	if nil != err {
		fmt.Println(err.Error())
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		blog.Debugf("haha %s. en\\en, always %d and %f", "eddie", 18, 3.1415)
	}
}
