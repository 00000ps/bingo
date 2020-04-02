package utils

import (
	"bingo/pkg/cmd"
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// func GetRemoteImg(url string) []byte {
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		log.Error("------------  get remote img  : %#v %#v", err.Error(), url)
// 		return nil
// 	}
// 	defer resp.Body.Close()
// 	pix, err := Read(resp.Body)
// 	return pix
// }

// GetLinkTarget used to replace the softlink by real object
func GetLinkTarget(cont string) string {
	if s := strings.Split(cont, "->"); len(s) == 2 {
		return s[1]
	}
	return ""
}

// GetSysLinkTarget used to replace the softlink by real object
func GetSysLinkTarget(file string) string {
	tarfile, err := os.Readlink(file)
	if err != nil {
		return ""
	}
	return tarfile
}

func GetParentBasename(file string) string { return path.Base(path.Dir(file)) }

// GetBackFile return the backup file name
func GetBackFile(file string) string { return file + ".bingo." + GetTimestamp() }

func IsExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return os.IsExist(err)
	}
	return true
}
func IsDirExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return fi.IsDir()
}
func IsFileExist(path string) bool {
	fi, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return !fi.IsDir()
}

// ListFiles generates a list of directory
func ListFiles(dirPath string) (list []string) {
	dir, err := ioutil.ReadDir(dirPath)
	if err == nil {
		for _, fi := range dir {
			local := dirPath + "/" + fi.Name()
			if fi.IsDir() {
				list = append(list, ListFiles(local)...)
			} else {
				if IsFileExist(local) {
					list = append(list, local)
				}
			}
		}
	}
	return
}

func Rename(oldName, newName string) error {
	if d := path.Dir(newName); !IsDirExist(d) {
		CreateDir(d)
	}
	return os.Rename(oldName, newName)
}
func Remove(name string) error {
	if IsFileExist(name) {
		return os.Remove(name)
	}
	return nil
}

func CopyFile(src, dst string) (w int64, err error) {
	dir := path.Dir(dst)
	if !IsDirExist(dir) {
		CreateDir(dir)
	}
	srcF, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcF.Close()
	dstF, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dstF.Close()
	return io.Copy(dstF, srcF)
}

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}
func GetParentDirectory(dirctory string) string {
	return substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		// log.Err(err)
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func WriteNewFileString(file string, data string) error { return WriteNewFile(file, []byte(data)) }
func WriteNewFile(file string, data []byte) error {
	dir := path.Dir(file)
	if !IsDirExist(dir) {
		CreateDir(dir)
	}
	if IsFileExist(file) {
		Remove(file)
	}
	fi, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		// log.Warning("WriteNewFile error,err msg is", err)
		return err
	}
	defer fi.Close()
	_, err = fi.Write(data)
	return err
}
func WriteBase64(file string, b64 string) error {
	d, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return err
	}
	return WriteNewFile(file, d)
}
func WriteNewJSON(file string, v interface{}) error {
	d, err := JSONBeauty(v)
	if err != nil {
		// log.Warning("WriteNewJson error,err msg is %s", err)
		return err
	}
	return WriteNewFile(file, d)
}
func WriteFile(file string, data []byte) error {
	dir := path.Dir(file)
	if !IsDirExist(dir) {
		CreateDir(dir)
	}
	fi, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		return err
	}
	defer fi.Close()
	_, err = fi.Write(data)
	return err
}
func CreateFile(file string) error { return WriteFile(file, []byte{}) }
func AppendString(file, str string, a ...interface{}) error {
	return WriteFile(file, []byte(fmt.Sprintf(str, a...)))
}
func AppendLine(file, str string) error {
	if IsWindows() {
		return AppendString(file, str+"\r\n")
	}
	return AppendString(file, str+"\n")
}
func AppendBytesLine(file string, data []byte) error {
	if IsWindows() {
		return WriteFile(file, append(data, []byte("\r\n")...))
	}
	return WriteFile(file, append(data, []byte("\n")...))
}

func CopyDir(src, dst string) error {
	if c, str := cmd.RunCmd("cp -r " + src + " " + dst); !c {
		return errors.New(str)
	}
	return nil
}
func CreateDir(dname string) error {
	if err := os.MkdirAll(dname, 0777); err != nil {
		return err
	}
	return nil
}
func CreateTsDir(dname string) (string, error) {
	path := dname + "_" + GetTimestamp()
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return "", err
	}
	return path, nil
}

func ReadFileOnce(file string) ([]byte, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}
func ReadFileAll(file string) ([]byte, error) { return ReadBlock(file) }
func ReadFile(file string) ([]byte, error)    { return ReadBlock(file) }

//  ReadJSON used to read from file and marshal it to json format
func ReadJSON(file string, v interface{}) (data []byte, err error) {
	data, err = ReadFile(file)
	if err != nil {
		// log.Notice("----read file err %v", err)
		return
	}
	err = json.Unmarshal([]byte(data), v)
	if err != nil {
		// log.Notice("----json err %v", err)
		return
	}
	return
}

//func processBlock(line []byte) {
//	os.Stdout.Write(line)
//}

// FileInfo returns file basic info
func FileInfo(filepath string) (os.FileInfo, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return f.Stat()
}

const defaultBlockSize = 10 * 1024

// ReadBlock is the simple call of ReadFileByBlock
func ReadBlock(file string) ([]byte, error) {
	return ReadFileByBlock(file, defaultBlockSize, nil)
}

// ReadFileByBlock will use less mem
func ReadFileByBlock(filePth string, bufSize int, hookfn func([]byte)) ([]byte, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return BlockReader(f, bufSize, hookfn)
}

// Read is the simple call of BlockReader
func Read(r io.Reader) ([]byte, error) { return BlockReader(r, defaultBlockSize, nil) }

// BlockReader will use less mem
func BlockReader(r io.Reader, bufSize int, hookfn func([]byte)) ([]byte, error) {
	var data []byte
	// defer r.Close()
	if bufSize == 0 || bufSize == -1 {
		bufSize = defaultBlockSize
	}
	buf := make([]byte, bufSize) //一次读取多少个字节
	bfRd := bufio.NewReader(r)
	for {
		n, err := bfRd.Read(buf)
		if hookfn == nil {
			data = append(data, buf[:n]...)
		} else {
			hookfn(buf[:n]) //n是成功读取字节数
		}
		if err != nil { //遇到任何错误立即返回，并忽略EOF错误信息
			if err == io.EOF {
				return data, nil
			}
			return data, err
		}
	}
	return data, nil
}

// GetFilesAndDirs 获取指定目录下的所有文件和目录
func GetFilesAndDirs(dirPth string) (files []string, dirs []string, err error) {
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, nil, err
	}
	pthSep := string(os.PathSeparator)
	//suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写
	for _, fi := range dir {
		// 过滤指定格式
		if strings.HasPrefix(fi.Name(), "~") ||
			// strings.HasPrefix(fi.Name(), "~")  ||
			strings.HasPrefix(fi.Name(), ".") {
			// logger.Warning("skip file: %s", path.Clean(dirPth+pthSep+fi.Name()))
			continue
		}
		// fmt.Println(dirPth + pthSep + fi.Name())
		if fi.IsDir() { // 目录, 递归遍历
			dirs = append(dirs, path.Clean(dirPth+pthSep+fi.Name()))
			fs, ds, e := GetFilesAndDirs(dirPth + pthSep + fi.Name())
			files = append(files, fs...)
			dirs = append(dirs, ds...)
			err = e
		} else {
			files = append(files, path.Clean(dirPth+pthSep+fi.Name()))
		}
	}
	return files, dirs, nil
}
