//nolint:goconst
package fileutils

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestMain(m *testing.M) {
	// we need a read only file for testing this package
	// git wont allow readonly file permissions
	err := os.Chmod("testdata/0444.perm", 0444)
	if err != nil {
		log.Fatal(err)
	}
}

func TestIsSymlink(t *testing.T) {
	tests := map[string]struct {
		path        string
		expected    bool
		shouldError bool
	}{
		"symlink folder": {
			path:     "testdata/testsymlinkfolder",
			expected: true,
		},
		"regular folder": {
			path:     "testdata/testfolder",
			expected: false,
		},
		"symlink file": {
			path:     "testdata/testsymlinkfile.txt",
			expected: true,
		},
		"regular file": {
			path:     "testdata/testfile.txt",
			expected: false,
		},
		"perms 644": {
			path:     "testdata/0644.perm",
			expected: false,
		},
		"perms 444": {
			path:     "testdata/0444.perm",
			expected: false,
		},
		"missing path": {
			path:        "testdata/nofile.txt",
			expected:    false,
			shouldError: true,
		},
	}

	for name, tt := range tests {
		actual, err := IsSymlink(tt.path)

		if err == nil && tt.shouldError {
			t.Errorf("%s: expected error, got nil", name)
		}

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestFileExists(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected bool
	}{
		"symlink folder": {
			path:     "testdata/testsymlinkfolder",
			expected: true,
		},
		"regular folder": {
			path:     "testdata/testfolder",
			expected: true,
		},
		"symlink file": {
			path:     "testdata/testsymlinkfile.txt",
			expected: true,
		},
		"regular file": {
			path:     "testdata/testfile.txt",
			expected: true,
		},
		"perms 644": {
			path:     "testdata/0644.perm",
			expected: true,
		},
		"perms 444": {
			path:     "testdata/0444.perm",
			expected: true,
		},
		"missing path": {
			path:     "testdata/nofile.txt",
			expected: false,
		},
	}

	for name, tt := range tests {
		actual := FileExists(tt.path)

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestMkDir(t *testing.T) {
	tests := map[string]struct {
		path        string
		shouldError bool
		cleanupFn   func(f string) error
	}{
		"symlink folder": {
			path:        "testdata/testsymlinkfolder",
			shouldError: false,
			cleanupFn: func(f string) error {
				return nil
			},
		},
		"regular folder": {
			path:        "testdata/testfolder",
			shouldError: false,
			cleanupFn: func(f string) error {
				return nil
			},
		},
		"symlink file": {
			path:        "testdata/testsymlinkfile.txt",
			shouldError: false,
			cleanupFn: func(f string) error {
				return nil
			},
		},
		"regular file": {
			path:        "testdata/testfile.txt",
			shouldError: true,
			cleanupFn: func(f string) error {
				return nil
			},
		},
		"missing path": {
			path:        "testdata/nofile.txt",
			shouldError: false,
			cleanupFn:   os.Remove,
		},
	}

	for name, tt := range tests {
		mkDirErr := MkDir(tt.path)
		exists := FileExists(tt.path)
		cleanupErr := tt.cleanupFn(tt.path)
		if cleanupErr != nil {
			t.Fatalf("cleanup error %v", cleanupErr)
		}

		if mkDirErr == nil && tt.shouldError {
			t.Errorf("%s: expected error, got nil", name)
		}

		if !exists {
			t.Errorf("%s: folder %v was not created", name, tt.path)
		}
	}
}

func TestFind(t *testing.T) {
	tests := map[string]struct {
		path      string
		extension string
		expected  []string
	}{
		"testdata": {
			path:      "testdata",
			extension: ".txt",
			expected: []string{
				"testdata/testfile.txt",
				"testdata/testsymlinkfile.txt",
			},
		},
	}

	for name, tt := range tests {
		actual, err := Find(tt.path, tt.extension)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestFolders(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected []string
	}{
		"testdata": {
			path: "testdata",
			expected: []string{
				"testdata/testfolder",
			},
		},
	}

	for name, tt := range tests {
		actual, err := Folders(tt.path)
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(tt.expected, actual) {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestEmptyFolder(t *testing.T) {
	targetFolder := "testdata/testfolder"
	extension := ".txt"

	for i := 0; i <= 3; i++ {
		err := MkFile(fmt.Sprintf("%v/%v%v", targetFolder, strconv.Itoa(i), extension))
		if err != nil {
			t.Fatal(err)
		}
	}

	err := EmptyFolder(targetFolder)
	if err != nil {
		t.Fatal(err)
	}

	files, err := Find(targetFolder, extension)
	if err != nil {
		t.Fatal(err)
	}

	if len(files) > 0 {
		t.Errorf("expected %v to be empty", targetFolder)
	}
}

func TestMkFile(t *testing.T) {
	targetFile := "testdata/testfolder/testfile.txt"

	err := MkFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(targetFile)

	if !FileExists(targetFile) {
		t.Errorf("expected %v to exist", targetFile)
	}
}

func TestFolderIsWriteable(t *testing.T) {
	tests := map[string]struct {
		path        string
		expected    bool
		shouldError bool
	}{
		"symlink folder": {
			path:     "testdata/testsymlinkfolder",
			expected: true,
		},
		"regular folder": {
			path:     "testdata/testfolder",
			expected: true,
		},
		"symlink file": {
			path:     "testdata/testsymlinkfile.txt",
			expected: false,
		},
		"regular file": {
			path:     "testdata/testfile.txt",
			expected: false,
		},
		"perms 644": {
			path:     "testdata/0644.perm",
			expected: false,
		},
		"perms 444": {
			path:     "testdata/0444.perm",
			expected: false,
		},
		"missing": {
			path:        "testdata/nofolder",
			expected:    false,
			shouldError: true,
		},
	}

	for name, tt := range tests {
		actual, err := FolderIsWriteable(tt.path)

		if err == nil && tt.shouldError {
			t.Errorf("%s: expected error, got nil", name)
		}

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestFileIsWriteable(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected bool
	}{
		"symlink folder": {
			path:     "testdata/testsymlinkfolder",
			expected: true,
		},
		"regular folder": {
			path:     "testdata/testfolder",
			expected: true,
		},
		"symlink file": {
			path:     "testdata/testsymlinkfile.txt",
			expected: true,
		},
		"regular file": {
			path:     "testdata/testfile.txt",
			expected: true,
		},
		"perms 644": {
			path:     "testdata/0644.perm",
			expected: true,
		},
		"perms 444": {
			path:     "testdata/0444.perm",
			expected: false,
		},
		"missing": {
			path:     "testdata/nofolder",
			expected: false,
		},
	}

	for name, tt := range tests {
		actual := FileIsWriteable(tt.path)

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestIsFile(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected bool
	}{
		"symlink folder": {
			path:     "testdata/testsymlinkfolder",
			expected: false,
		},
		"regular folder": {
			path:     "testdata/testfolder",
			expected: false,
		},
		"symlink file": {
			path:     "testdata/testsymlinkfile.txt",
			expected: true,
		},
		"regular file": {
			path:     "testdata/testfile.txt",
			expected: true,
		},
		"perms 644": {
			path:     "testdata/0644.perm",
			expected: true,
		},
		"perms 444": {
			path:     "testdata/0444.perm",
			expected: true,
		},
		"missing": {
			path:     "testdata/nofile",
			expected: false,
		},
	}

	for name, tt := range tests {
		actual := IsFile(tt.path)

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestIsFolder(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected bool
	}{
		"symlink folder": {
			path:     "testdata/testsymlinkfolder",
			expected: true,
		},
		"regular folder": {
			path:     "testdata/testfolder",
			expected: true,
		},
		"symlink file": {
			path:     "testdata/testsymlinkfile.txt",
			expected: false,
		},
		"regular file": {
			path:     "testdata/testfile.txt",
			expected: false,
		},
		"perms 644": {
			path:     "testdata/0644.perm",
			expected: false,
		},
		"perms 444": {
			path:     "testdata/0444.perm",
			expected: false,
		},
		"missing": {
			path:     "testdata/nofolder",
			expected: false,
		},
	}

	for name, tt := range tests {
		actual := IsFolder(tt.path)

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestFileNameWithoutExtension(t *testing.T) {
	tests := map[string]struct {
		path     string
		expected string
	}{
		"regular file": {
			path:     "testfile.txt",
			expected: "testfile",
		},
		"trailing slash": {
			path:     "/testfile.txt",
			expected: "testfile",
		},
		"inside folder": {
			path:     "testdata/testfile.txt",
			expected: "testfile",
		},
	}

	for name, tt := range tests {
		actual := FileNameWithoutExtension(tt.path)

		if tt.expected != actual {
			t.Errorf("%s: expected %v, got %v", name, tt.expected, actual)
		}
	}
}

func TestGetFile(t *testing.T) {
	targetFile := "testdata/testfolder/testfile.txt"
	expectType := "*os.File"

	err := MkFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(targetFile)

	f, err := GetFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	if fmt.Sprintf("%T", f) != expectType {
		t.Errorf("expected %v to be %v", targetFile, expectType)
	}
}

func TestWriteLine(t *testing.T) {
	targetFile := "testdata/testfolder/testfile.txt"
	expected := "testline"

	err := MkFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(targetFile)

	f, err := GetFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = WriteLine(f, expected)
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expected {
		t.Errorf("%v: expected content to equal %v [%v]", targetFile, expected, string(content))
	}
}

func TestGetMD5Hash(t *testing.T) {
	targetFile := "testdata/testfile.txt"
	expected := "d41d8cd98f00b204e9800998ecf8427e"

	hash, err := GetMD5Hash(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if hash != expected {
		t.Errorf("%v: expected hash to equal %v [%v]", targetFile, expected, hash)
	}
}

func TestFileSizeBytes(t *testing.T) {
	targetFile := "testdata/testfolder/testfile.txt"
	testline := "testline"
	expected := int64(8)

	err := MkFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(targetFile)

	f, err := GetFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	err = WriteLine(f, testline)
	if err != nil {
		t.Fatal(err)
	}

	size, err := FileSizeBytes(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if size != expected {
		t.Errorf("%v: expected hash to equal %v [%v]", targetFile, expected, size)
	}
}

func TestFileHash(t *testing.T) {
	targetFile := "testdata/testfile.txt"
	expected := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	hash, err := FileHash(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if hash != expected {
		t.Errorf("%v: expected hash to equal %v [%v]", targetFile, expected, hash)
	}
}

func TestWriteFile(t *testing.T) {
	targetFile := "testdata/test_writefile.txt"
	testContent := "test content"

	err := WriteFile(targetFile, testContent)
	if err != nil {
		t.Fatal(err)
	}

	if !FileExists(targetFile) {
		t.Errorf("%v: was not written", targetFile)
	}

	contents, err := os.ReadFile(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Remove(targetFile)
	if err != nil {
		t.Fatal(err)
	}

	if string(contents) != testContent {
		t.Errorf("%v: expected contents to equal %v [%v]", targetFile, testContent, string(contents))
	}
}
