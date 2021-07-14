// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package fsnotify

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestEventStringWithValue(t *testing.T) {
	formatEvent := func(name, oldName string, op Op, id uint64) string {
		return fmt.Sprintf("newpath=%s, oldpath=%s, op=%s, eventID=%d", name, oldName, op.String(), id)
	}
	for opMask, expectedString := range map[Op]string{
		Chmod | Create: formatEvent("/usr/someFile", "", Create|Chmod, 0),
		Rename:         formatEvent("/usr/someFile", "", Rename, 0),
		Remove:         formatEvent("/usr/someFile", "", Remove, 0),
		Write | Chmod:  formatEvent("/usr/someFile", "", Write|Chmod, 0),
	} {
		event := Event{Name: "/usr/someFile", Op: opMask, OldName: "", ID: 0}
		if event.String() != expectedString {
			t.Fatalf("Expected %s, got: %v", expectedString, event.String())
		}
	}
	renameEvent := Event{Name: "/usr/someFile_Rename", Op: Rename, OldName: "/usr/someFile", ID: 123}
	renameEventExpectedString := formatEvent("/usr/someFile_Rename", "/usr/someFile", Rename, 123)
	if renameEvent.String() != renameEventExpectedString {
		t.Fatalf("Expected %s, got: %v", renameEventExpectedString, renameEvent.String())
	}
}

func TestEventOpStringWithValue(t *testing.T) {
	expectedOpString := "WRITE|CHMOD"
	event := Event{Name: "someFile", Op: Write | Chmod}
	if event.Op.String() != expectedOpString {
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	}
}

func TestEventOpStringWithNoValue(t *testing.T) {
	expectedOpString := ""
	event := Event{Name: "testFile", Op: 0}
	if event.Op.String() != expectedOpString {
		t.Fatalf("Expected %s, got: %v", expectedOpString, event.Op.String())
	}
}

// TestWatcherClose tests that the goroutine started by creating the watcher can be
// signalled to return at any time, even if there is no goroutine listening on the events
// or errors channels.
func TestWatcherClose(t *testing.T) {
	t.Parallel()

	name := tempMkFile(t, "")
	w := newWatcher(t)
	err := w.Add(name)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(name)
	if err != nil {
		t.Fatal(err)
	}
	// Allow the watcher to receive the event.
	time.Sleep(time.Millisecond * 100)

	err = w.Close()
	if err != nil {
		t.Fatal(err)
	}
}
