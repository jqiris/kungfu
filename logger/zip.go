/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package logger

import (
	"archive/zip"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

//压缩文件
//files 文件数组，可以是不同dir下的文件或者文件夹
//dest 压缩文件存放地址
func Compress(files []*os.File, dest string) error {
	d, _ := os.Create(dest)
	defer d.Close()
	w := zip.NewWriter(d)
	defer w.Close()
	for _, file := range files {
		err := compress(file, "", w)
		if err != nil {
			return err
		}
	}
	return nil
}

func compress(file *os.File, prefix string, zw *zip.Writer) error {
	info, err := file.Stat()
	if err != nil {
		return err
	}
	if info.IsDir() {
		prefix = prefix + "/" + info.Name()
		fileInfos, err := file.Readdir(-1)
		if err != nil {
			return err
		}
		for _, fi := range fileInfos {
			f, err := os.Open(file.Name() + "/" + fi.Name())
			if err != nil {
				return err
			}
			err = compress(f, prefix, zw)
			if err != nil {
				return err
			}
		}
	} else {
		header, err := zip.FileInfoHeader(info)
		header.Name = prefix + "/" + header.Name
		if err != nil {
			return err
		}
		writer, err := zw.CreateHeader(header)
		if err != nil {
			return err
		}
		_, err = io.Copy(writer, file)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

//解压
func DeCompress(zipFile, dest string) (err error) {
	//目标文件夹不存在则创建
	if _, err = os.Stat(dest); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(dest, 0755)
		}
	}

	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}

	defer reader.Close()

	for _, file := range reader.File {
		if file.FileInfo().IsDir() {

			err := os.MkdirAll(dest+"/"+file.Name, 0755)
			if err != nil {
				Error(err)
			}
			continue
		} else {

			err = os.MkdirAll(getDir(dest+"/"+file.Name), 0755)
			if err != nil {
				return err
			}
		}

		if err := DeCompressFile(file, dest); err != nil {
			return err
		}
	}
	return
}

func DeCompressFile(file *zip.File, dest string) error {
	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	filename := dest + "/" + file.Name

	w, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer w.Close()

	_, err = io.Copy(w, rc)
	if err != nil {
		return err
	}
	return nil
}

func getDir(path string) string {
	return subString(path, 0, strings.LastIndex(path, "/"))
}

func subString(str string, start, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		Error("start is wrong")
		return ""
	}

	if end < start || end > length {
		Error("end is wrong")
		return ""
	}

	return string(rs[start:end])
}

func CompressZip(src string, dest string) (err error) {

	f, err := ioutil.ReadDir(src)
	if err != nil {
		Error(err)
	}

	fZip, _ := os.Create(dest)
	w := zip.NewWriter(fZip)

	defer fZip.Close()
	defer w.Close()
	for _, file := range f {
		fw, _ := w.Create(file.Name())
		fileContent, err := ioutil.ReadFile(src + file.Name())
		if err != nil {
			Error(err)
		}
		_, err = fw.Write(fileContent)

		if err != nil {
			Error(err)
		}
	}
	return
}
