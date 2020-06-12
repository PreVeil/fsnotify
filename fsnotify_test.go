// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !plan9

package fsnotify

import (
	"os"
	"testing"
	"time"
)

func TestEventStringWithValue(t *testing.T) {
	for opMask, expectedString := range map[Op]string{
		Chmod | Create: `newpath="/usr/someFile", oldpath="", op=CREATE|CHMOD, eventID=0`,
		Rename:         `newpath="/usr/someFile", oldpath="", op= RENAME, eventID=0`,
		Remove:         `newpath="/usr/someFile", oldpath="", op=REMOVE, eventID=0`,
		Write | Chmod:  `newpath="/usr/someFile", oldpath="", op=WRITE|CHMOD, eventID=0`,
	} {
		event := Event{Name: "/usr/someFile", Op: opMask, OldName: "", ID: 0}
		if event.String() != expectedString {
			t.Fatalf("Expected %s, got: %v", expectedString, event.String())
		}

	}
	renameEvent := Event{Name: "/usr/someFile_Rename", Op: Rename, OldName: "/usr/someFile", ID: 0}
	renameEventExpectedString := `newpath="/usr/someFile_Rename", oldpath="/usr/someFile", op= RENAME, eventID=0`
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
