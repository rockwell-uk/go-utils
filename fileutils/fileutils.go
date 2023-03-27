package fileutils

//nolint:gosec
import (
	"crypto/md5"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

const (
	packageName = "fileutils"
)

func Find(folderPath, ext string) ([]string, error) {
	var funcName string = "Find"

	var files []string

	if !FolderExists(folderPath) {
		return []string{}, fmt.Errorf("%v.%v: target does not exist [%v]", packageName, funcName, folderPath)
	}

	if !IsFolder(folderPath) {
		return []string{}, fmt.Errorf("%v.%v: target is not a folder [%v]", packageName, funcName, folderPath)
	}

	sym, err := IsSymlink(folderPath)
	if err != nil {
		return []string{}, fmt.Errorf("%v.%v: error checking symlink [%v], [%v]", packageName, funcName, folderPath, err.Error())
	}
	if sym {
		return []string{}, fmt.Errorf("%v.%v: target is a symlink [%v]", packageName, funcName, folderPath)
	}

	err = filepath.WalkDir(folderPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			files = append(files, s)
		}
		return nil
	})

	if err != nil {
		return []string{}, fmt.Errorf("%v.%v: error walking target [%v], [%v]", packageName, funcName, folderPath, err.Error())
	}

	return files, nil
}

func Folders(folderPath string) ([]string, error) {
	var funcName string = "Folders"

	if !FolderExists(folderPath) {
		return []string{}, fmt.Errorf("%v.%v: target does not exist [%v]", packageName, funcName, folderPath)
	}

	if !IsFolder(folderPath) {
		return []string{}, fmt.Errorf("%v.%v: target is not a folder [%v]", packageName, funcName, folderPath)
	}

	var folders []string

	err := filepath.WalkDir(folderPath, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if d.IsDir() && s != folderPath {
			folders = append(folders, s)
		}

		return nil
	})

	if err != nil {
		return []string{}, fmt.Errorf("%v.%v: error walking target [%v], [%v]", packageName, funcName, folderPath, err.Error())
	}

	return folders, nil
}

func EmptyFolder(folderPath string) error {
	var funcName string = "EmptyFolder"

	if !FolderExists(folderPath) {
		return fmt.Errorf("%v.%v: target does not exist [%v]", packageName, funcName, folderPath)
	}

	dir, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return fmt.Errorf("%v.%v: error reading target [%v], [%v]", packageName, funcName, folderPath, err.Error())
	}

	for _, d := range dir {
		os.RemoveAll(path.Join([]string{folderPath, d.Name()}...))
	}

	return nil
}

func FolderIsWriteable(folderPath string) (bool, error) {
	var funcName string = "FolderIsWritable"

	if !FolderExists(folderPath) {
		return false, fmt.Errorf("%v.%v: target does not exist [%v]", packageName, funcName, folderPath)
	}

	if !IsFolder(folderPath) {
		return false, fmt.Errorf("%v.%v: target is not a folder [%v]", packageName, funcName, folderPath)
	}

	return FileIsWriteable(folderPath), nil
}

func IsFile(fileName string) bool {
	i, err := os.Stat(fileName)
	if err != nil {
		return false
	}
	return !i.IsDir()
}

func IsFolder(folderPath string) bool {
	i, err := os.Stat(folderPath)
	if err != nil {
		return false
	}
	return i.IsDir()
}

func FolderExists(folderPath string) bool {
	return FileExists(folderPath)
}

func FileExists(fileName string) bool {
	if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func IsSymlink(fileName string) (bool, error) {
	var funcName string = "IsSymlink"

	fi, err := os.Lstat(fileName)
	if err != nil {
		return false, fmt.Errorf("%v.%v: error checking file info [%v], [%v]", packageName, funcName, fileName, err.Error())
	}

	return fi.Mode()&os.ModeSymlink == os.ModeSymlink, nil
}

func MkDir(dir string) error {
	var funcName string = "MkDir"

	var sym bool
	var err error

	if FileExists(dir) {
		sym, err = IsSymlink(dir)
		if err != nil {
			return fmt.Errorf("%v.%v: error checking if symlink [%v], [%v]", packageName, funcName, dir, err.Error())
		}
	}

	if !sym {
		return os.MkdirAll(dir, os.ModePerm)
	}

	return nil
}

func MkFile(fileName string) error {
	var funcName string = "MkFile"

	if FileExists(fileName) {
		return fmt.Errorf("%v.%v: file already exists [%v]", packageName, funcName, fileName)
	}

	emptyFile, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("%v.%v: error creating file [%v]", packageName, funcName, fileName)
	}
	emptyFile.Close()

	return nil
}

func FileNameWithoutExtension(fileName string) string {
	return strings.TrimSuffix(filepath.Base(fileName), filepath.Ext(fileName))
}

func FileIsWriteable(fileName string) bool {
	return syscall.Access(fileName, syscall.O_RDWR) == nil
}

func GetFile(fileName string) (*os.File, error) {
	var funcName string = "GetFile"

	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0655)
	if err != nil {
		return nil, fmt.Errorf("%v.%v: error opening file [%v], [%v]", packageName, funcName, fileName, err.Error())
	}

	return f, nil
}

func WriteFile(fileName string, fileContent string) error {
	var funcName string = "WriteFile"

	f, err := GetFile(fileName)
	if err != nil {
		return fmt.Errorf("%v.%v: error preparing to write file [%v], [%v]", packageName, funcName, fileName, err.Error())
	}
	defer f.Close()

	_, err = f.Write([]byte(fileContent))
	if err != nil {
		return fmt.Errorf("%v.%v: error writing file [%v], [%v]", packageName, funcName, fileName, err.Error())
	}

	return nil
}

func WriteLine(f *os.File, line string) error {
	var funcName string = "WriteLine"

	if _, err := f.Write([]byte(line)); err != nil {
		return fmt.Errorf("%v.%v: error writing line [%v], [%v]", packageName, funcName, f.Name(), err.Error())
	}

	return nil
}

func GetMD5Hash(filePath string) (string, error) {
	var funcName string = "GetMD5Hash"

	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("%v.%v: error opening file [%v], [%v]", packageName, funcName, filePath, err.Error())
	}
	defer file.Close()

	hash := md5.New() //nolint:gosec
	_, err = io.Copy(hash, file)
	if err != nil {
		return "", fmt.Errorf("%v.%v: error hash file [%v], [%v]", packageName, funcName, filePath, err.Error())
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func FileSizeBytes(filePath string) (int64, error) {
	var funcName string = "FileSizeBytes"

	file, err := os.Open(filePath)
	if err != nil {
		return 0, fmt.Errorf("%v.%v: error opening file [%v], [%v]", packageName, funcName, filePath, err.Error())
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return 0, fmt.Errorf("%v.%v: error getting file info [%v], [%v]", packageName, funcName, filePath, err.Error())
	}

	return stat.Size(), nil
}

func FileSize(fileName, units string) (string, error) {
	fileBytes, err := FileSizeBytes(fileName)
	if err != nil {
		return "", err
	}
	var bytes float64 = float64(fileBytes)

	switch units {
	case "kb":
	case "mb":
	case "gb":
	case "tb":
	case "pb":
	case "xb":
	case "zb":
		return fmt.Sprintf("%.2f%v", ByteSizeConvert(fileBytes, units), units), nil
	}

	return fmt.Sprintf("%.0fb", bytes), nil
}

func FileHash(fileName string) (string, error) {
	var funcName string = "FileHash"

	f, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("%v.%v: error opening file [%v], [%v]", packageName, funcName, fileName, err.Error())
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("%v.%v: error hashing file [%v], [%v]", packageName, funcName, fileName, err.Error())
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
