// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

// These tests all do multiple operations while a user is unstaged.

package test

import (
	"testing"
	"time"
)

// bob renames a non-conflicting file into a new directory while unstaged
func TestCrUnmergedRenameIntoNewDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/c", "world"),
		),
		as(bob, noSync(),
			rename("a/b", "d/e"),
			reenableUpdates(),
			lsdir("a/", m{"c": "FILE"}),
			lsdir("d/", m{"e": "FILE"}),
			read("a/c", "world"),
			read("d/e", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c": "FILE"}),
			lsdir("d/", m{"e": "FILE"}),
			read("a/c", "world"),
			read("d/e", "hello"),
		),
	)
}

// alice renames a non-conflicting file into a new directory while
// bob is unstaged.
func TestCrMergedRenameIntoNewDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "d/e"),
		),
		as(bob, noSync(),
			write("a/c", "world"),
			reenableUpdates(),
			lsdir("a/", m{"c": "FILE"}),
			lsdir("d/", m{"e": "FILE"}),
			read("a/c", "world"),
			read("d/e", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c": "FILE"}),
			lsdir("d/", m{"e": "FILE"}),
			read("a/c", "world"),
			read("d/e", "hello"),
		),
	)
}

// bob causes a simple rename(cycle while unstaged),
func TestCrRenameCycle(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("b"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("b", "a/b"),
		),
		as(bob, noSync(),
			rename("a", "b/a"),
			reenableUpdates(),
			lsdir("a/", m{"b": "DIR"}),
			lsdir("a/b/", m{"a": "SYM"}),
			lsdir("a/b/a", m{"b": "DIR"}),
		),
		as(alice,
			lsdir("a/", m{"b": "DIR"}),
			lsdir("a/b/", m{"a": "SYM"}),
			lsdir("a/b/a", m{"b": "DIR"}),
			write("a/c", "hello"),
		),
		as(bob,
			read("a/b/a/c", "hello"),
		),
	)
}

// bob causes a complicated rename(cycle while unstaged),
func TestCrComplexRenameCycle(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("b"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("b", "a/b"),
		),
		as(bob, noSync(),
			mkdir("b/c"),
			rename("a", "b/c/a"),
			reenableUpdates(),
			lsdir("a/", m{"b": "DIR"}),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"a": "SYM"}),
			lsdir("a/b/c/a", m{"b": "DIR"}),
		),
		as(alice,
			lsdir("a/", m{"b": "DIR"}),
			lsdir("a/b/", m{"c": "DIR"}),
			lsdir("a/b/c", m{"a": "SYM"}),
			lsdir("a/b/c/a", m{"b": "DIR"}),
			write("a/d", "hello"),
		),
		as(bob,
			read("a/b/c/a/d", "hello"),
		),
	)
}

// bob causes a complicated and large rename(cycle while unstaged),
func TestCrComplexLargeRenameCycle(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b/c"),
			mkdir("d/e/f"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("d", "a/b/c/d"),
		),
		as(bob, noSync(),
			mkdir("d/e/f/g/h/i"),
			rename("a", "d/e/f/g/h/i/a"),
			reenableUpdates(),
			lsdir("a/b/c/d/e/f/g/h/i", m{"a": "SYM"}),
			lsdir("a/b/c/d/e/f/g/h/i/a", m{"b": "DIR"}),
		),
		as(alice,
			lsdir("a/b/c/d/e/f/g/h/i", m{"a": "SYM"}),
			lsdir("a/b/c/d/e/f/g/h/i/a", m{"b": "DIR"}),
			write("a/j", "hello"),
		),
		as(bob,
			read("a/b/c/d/e/f/g/h/i/a/j", "hello"),
		),
	)
}

// bob and alice do a lot of complex renames cycle while unstaged
func TestCrComplexRenameNoCycle(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b/c/d/e/f/g"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b/c/d/e/f", "f"),
			rename("a/b/c/d", "f/g/d"),
			rename("a/b", "f/g/d/e/b"),
		),
		as(bob, noSync(),
			rename("a/b/c/d/e/f/g", "g"),
			rename("a/b/c/d/e", "g/e"),
			rename("a/b/c", "g/e/f/c"),
			rename("a", "g/e/f/c/d/a"),
			reenableUpdates(),
			lsdir("f", m{"c": "DIR"}),
			lsdir("f/c", m{}),
			lsdir("g", m{"e": "DIR", "d": "DIR"}),
			lsdir("g/e", m{"b": "DIR"}),
			lsdir("g/e/b", m{}),
			lsdir("g/d", m{"a": "DIR"}),
		),
		as(alice,
			lsdir("f", m{"c": "DIR"}),
			lsdir("f/c", m{}),
			lsdir("g", m{"e": "DIR", "d": "DIR"}),
			lsdir("g/e", m{"b": "DIR"}),
			lsdir("g/e/b", m{}),
			lsdir("g/d", m{"a": "DIR"}),
		),
	)
}

// bob renames a file while unmerged, at the same time alice writes to it
func TestCrUnmergedRenameWithParallelWrite(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("b"),
			write("a/foo", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/foo", "goodbye"),
		),
		as(bob, noSync(),
			rename("a/foo", "b/bar"),
			reenableUpdates(),
			lsdir("a", m{}),
			lsdir("b", m{"bar": "FILE"}),
			read("b/bar", "goodbye"),
		),
		as(alice,
			lsdir("a", m{}),
			lsdir("b", m{"bar": "FILE"}),
			read("b/bar", "goodbye"),
		),
	)
}

// bob makes a non-conflicting file executable while alice writes to it
func TestCrUnmergedSetexParallelWrite(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "goodbye"),
		),
		as(bob, noSync(),
			setex("a/b", true),
			reenableUpdates(),
			lsdir("a/", m{"b": "EXEC"}),
			read("a/b", "goodbye"),
		),
		as(alice,
			lsdir("a/", m{"b": "EXEC"}),
			read("a/b", "goodbye"),
		),
	)
}

// alice makes a non-conflicting file executable while bob writes to it
func TestCrMergedSetexParallelWrite(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			setex("a/b", true),
		),
		as(bob, noSync(),
			write("a/b", "goodbye"),
			reenableUpdates(),
			lsdir("a/", m{"b": "EXEC"}),
			read("a/b", "goodbye"),
		),
		as(alice,
			lsdir("a/", m{"b": "EXEC"}),
			read("a/b", "goodbye"),
		),
	)
}

// bob writes to a file while alice removes it
func TestCrUnmergedWriteToRemovedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rm("a/b"),
		),
		as(bob, noSync(),
			write("a/b", "goodbye"),
			reenableUpdates(),
			lsdir("a/", m{"b": "FILE"}),
			read("a/b", "goodbye"),
		),
		as(alice,
			lsdir("a/", m{"b": "FILE"}),
			read("a/b", "goodbye"),
		),
	)
}

// bob removes a file while alice writes to it
func TestCrMergedWriteToRemovedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "goodbye"),
		),
		as(bob, noSync(),
			rm("a/b"),
			reenableUpdates(),
			lsdir("a/", m{"b": "FILE"}),
			read("a/b", "goodbye"),
		),
		as(alice,
			lsdir("a/", m{"b": "FILE"}),
			read("a/b", "goodbye"),
		),
	)
}

// bob writes to a file to a directory that alice removes
func TestCrUnmergedCreateInRemovedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b/c/d/e", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rm("a/b/c/d/e"),
			rmdir("a/b/c/d"),
			rmdir("a/b/c"),
			rmdir("a/b"),
		),
		as(bob, noSync(),
			write("a/b/c/d/f", "goodbye"),
			reenableUpdates(),
			lsdir("a/b/c/d", m{"f": "FILE"}),
			read("a/b/c/d/f", "goodbye"),
		),
		as(alice,
			lsdir("a/b/c/d", m{"f": "FILE"}),
			read("a/b/c/d/f", "goodbye"),
		),
	)
}

// alice writes to a file to a directory that bob removes
func TestCrMergedCreateInRemovedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkfile("a/b/c/d/e", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b/c/d/f", "goodbye"),
		),
		as(bob, noSync(),
			rm("a/b/c/d/e"),
			rmdir("a/b/c/d"),
			rmdir("a/b/c"),
			rmdir("a/b"),
			reenableUpdates(),
			lsdir("a/b/c/d", m{"f": "FILE"}),
			read("a/b/c/d/f", "goodbye"),
		),
		as(alice,
			lsdir("a/b/c/d", m{"f": "FILE"}),
			read("a/b/c/d/f", "goodbye"),
		),
	)
}

// bob writes a file while unmerged, at the same time alice renames it
func TestCrMergedRenameWithParallelWrite(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			mkdir("b"),
			write("a/foo", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/foo", "b/bar"),
		),
		as(bob, noSync(),
			write("a/foo", "goodbye"),
			reenableUpdates(),
			lsdir("a", m{}),
			lsdir("b", m{"bar": "FILE"}),
			read("b/bar", "goodbye"),
		),
		as(alice,
			lsdir("a", m{}),
			lsdir("b", m{"bar": "FILE"}),
			read("b/bar", "goodbye"),
		),
	)
}

// bob has two back-to-back resolutions
func TestCrDoubleResolution(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b/c", "hello"),
		),
		as(bob, noSync(),
			write("a/b/d", "goodbye"),
			reenableUpdates(),
			lsdir("a", m{"b": "DIR"}),
			lsdir("a/b", m{"c": "FILE", "d": "FILE"}),
			read("a/b/c", "hello"),
			read("a/b/d", "goodbye"),
		),
		as(alice,
			lsdir("a", m{"b": "DIR"}),
			lsdir("a/b", m{"c": "FILE", "d": "FILE"}),
			read("a/b/c", "hello"),
			read("a/b/d", "goodbye"),
			// Make a few more revisions
			write("a/b/e", "hello"),
			write("a/b/f", "goodbye"),
		),
		as(bob,
			read("a/b/e", "hello"),
			read("a/b/f", "goodbye"),
			disableUpdates(),
		),
		as(alice,
			rm("a/b/f"),
		),
		as(bob, noSync(),
			rm("a/b/e"),
			reenableUpdates(),
			lsdir("a", m{"b": "DIR"}),
			lsdir("a/b", m{"c": "FILE", "d": "FILE"}),
			read("a/b/c", "hello"),
			read("a/b/d", "goodbye"),
		),
		as(alice,
			lsdir("a", m{"b": "DIR"}),
			lsdir("a/b", m{"c": "FILE", "d": "FILE"}),
			read("a/b/c", "hello"),
			read("a/b/d", "goodbye"),
		),
	)
}

// bob makes files in a directory renamed by alice
func TestCrUnmergedMakeFilesInRenamedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rename("a/b", "b"),
		),
		as(bob, noSync(),
			write("a/b/c", "hello"),
			write("a/b/d", "goodbye"),
			reenableUpdates(),
			lsdir("a", m{}),
			lsdir("b", m{"c": "FILE", "d": "FILE"}),
			read("b/c", "hello"),
			read("b/d", "goodbye"),
		),
		as(alice,
			lsdir("a", m{}),
			lsdir("b", m{"c": "FILE", "d": "FILE"}),
			read("b/c", "hello"),
			read("b/d", "goodbye"),
		),
	)
}

// bob makes files in a directory renamed by alice
func TestCrMergedMakeFilesInRenamedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b/c", "hello"),
			write("a/b/d", "goodbye"),
		),
		as(bob, noSync(),
			rename("a/b", "b"),
			reenableUpdates(),
			lsdir("a", m{}),
			lsdir("b", m{"c": "FILE", "d": "FILE"}),
			read("b/c", "hello"),
			read("b/d", "goodbye"),
		),
		as(alice,
			lsdir("a", m{}),
			lsdir("b", m{"c": "FILE", "d": "FILE"}),
			read("b/c", "hello"),
			read("b/d", "goodbye"),
		),
	)
}

// bob moves and setexes a file that was written by alice
func TestCrConflictMoveAndSetexWrittenFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			rename("a/b", "a/c"),
			setex("a/c", true),
			reenableUpdates(),
			lsdir("a/", m{"c$": "EXEC"}),
			read("a/c", "world"),
		),
		as(alice,
			lsdir("a/", m{"c$": "EXEC"}),
			read("a/c", "world"),
		),
	)
}

// bob moves and setexes a file that was removed by alice
func TestCrConflictMoveAndSetexRemovedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rm("a/b"),
		),
		as(bob, noSync(),
			rename("a/b", "a/c"),
			setex("a/c", true),
			reenableUpdates(),
			lsdir("a/", m{"c$": "EXEC"}),
			read("a/c", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c$": "EXEC"}),
			read("a/c", "hello"),
		),
	)
}

// bob creates a directory with the same name that alice used for a
// file that used to exist at that location, but bob first moved it
func TestCrMergedRecreatedAndUnmergedMovedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/b", "world"),
		),
		as(bob, noSync(),
			rename("a/b", "a/d/b"),
			rm("a/d/b"),
			write("a/d/b/c", "uh oh"),
			reenableUpdates(),
			lsdir("a/", m{"d$": "DIR", "b$": "FILE"}),
			lsdir("a/d", m{"b$": "DIR"}),
			read("a/b", "world"),
			read("a/d/b/c", "uh oh"),
		),
		as(alice,
			lsdir("a/", m{"d$": "DIR", "b$": "FILE"}),
			lsdir("a/d", m{"b$": "DIR"}),
			read("a/b", "world"),
			read("a/d/b/c", "uh oh"),
		),
	)
}

func TestCrUnmergedCreateFileInRenamedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			write("a/c", "touch"),
		),
		as(bob, noSync(),
			mkdir("a/d"),
			rename("a/b", "a/d/e"),
			write("a/d/e/f", "hello"),
			reenableUpdates(),
			lsdir("a/", m{"c": "FILE", "d": "DIR"}),
			lsdir("a/d/", m{"e": "DIR"}),
			lsdir("a/d/e", m{"f": "FILE"}),
			read("a/d/e/f", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c": "FILE", "d": "DIR"}),
			lsdir("a/d/", m{"e": "DIR"}),
			lsdir("a/d/e", m{"f": "FILE"}),
			read("a/d/e/f", "hello"),
		),
	)
}

// bob moves a file that was removed by alice
func TestCrUnmergedMoveOfRemovedFile(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a"),
			write("a/b", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rm("a/b"),
		),
		as(bob, noSync(),
			rename("a/b", "a/c"),
			reenableUpdates(),
			lsdir("a/", m{"c$": "FILE"}),
			read("a/c", "hello"),
		),
		as(alice,
			lsdir("a/", m{"c$": "FILE"}),
			read("a/c", "hello"),
		),
	)
}

// bob makes, sets the mtime, and remove a file.  Regression test for
// KBFS-1163.
func TestCrUnmergedSetAttrOfAddedAndRemovedDir(t *testing.T) {
	test(t,
		users("alice", "bob"),
		as(alice,
			mkdir("a/b/c"),
			mkfile("a/b/c/d", "hello"),
		),
		as(bob,
			disableUpdates(),
		),
		as(alice,
			rm("a/b/c/d"),
			rm("a/b/c"),
			rm("a/b"),
			rm("a"),
		),
		as(bob, noSync(),
			setmtime("a/b/c", time.Now().Add(1*time.Minute)),
			mkfile("e", "world"),
			reenableUpdates(),
			lsdir("", m{"e$": "FILE"}),
			read("e", "world"),
		),
		as(alice,
			lsdir("", m{"e$": "FILE"}),
			read("e", "world"),
		),
	)
}
