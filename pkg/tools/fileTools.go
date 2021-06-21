package tools

//import (
//	"archive/zip"
//	"io"
//	"io/ioutil"
//	"os"
//	"path/filepath"
//	"strings"
//	"golang-common-base/pkg/logger"
//)
//
//// 判断文件夹是否存在
//func PathExists(path string) (bool, error) {
//	_, err := os.Stat(path)
//	if err == nil {
//		return true, nil
//	}
//	if os.IsNotExist(err) {
//		return false, nil
//	}
//	return false, err
//}
//
//// 创建文件夹
//func CreateDictByPath(path string) (bool, error) {
//	exist, err := PathExists(path)
//	if err != nil {
//		logger.Error("get dir error![%v]\n", err)
//		return false, err
//	}
//
//	if exist {
//		return true, nil
//	} else {
//		// 创建文件夹
//		err := os.MkdirAll(path, os.ModePerm)
//		if err != nil {
//			logger.Error("mkdir failed![%v]\n", err)
//			return false, err
//		} else {
//			return true, nil
//		}
//	}
//}
//
//// 删除文件夹
//func RemoveDictPath(path string) error {
//	err := os.RemoveAll(path)
//	return err
//}
//
//// 文件拷贝操作
//func CopyFile(sourceFile string, destFile string) error {
//	input, err := ioutil.ReadFile(sourceFile)
//	if err != nil {
//		logger.Error(err)
//		return err
//	}
//
//	err = ioutil.WriteFile(destFile, input, 0644)
//	if err != nil {
//		logger.Error("Error creating", destFile, err)
//		return err
//	}
//
//	return nil
//}
//
//// ========================= 压缩与解压缩 ==================================
//// 压缩 srcFile could be a single file or a directory
//func Zip(srcFile string, destZip string) error {
//	// 预防：旧文件无法覆盖
//	os.RemoveAll(destZip)
//
//	zipfile, err := os.Create(destZip)
//	if err != nil {
//		return err
//	}
//	defer zipfile.Close()
//
//	// 打开：zip文件
//	archive := zip.NewWriter(zipfile)
//	defer archive.Close()
//
//	// 遍历路径信息
//	filepath.Walk(srcFile, func(path string, info os.FileInfo, err error) error {
//		// 如果是源路径，提前进行下一个遍历
//		if path == srcFile {
//			return nil
//		}
//
//		if err != nil {
//			return err
//		}
//
//		header, err := zip.FileInfoHeader(info)
//		if err != nil {
//			return err
//		}
//
//		header.Name = strings.TrimPrefix(path, filepath.Dir(srcFile)+"/")
//		// header.Name = path
//		if info.IsDir() {
//			header.Name += "/"
//		} else {
//			header.Method = zip.Deflate
//		}
//
//		writer, err := archive.CreateHeader(header)
//		if err != nil {
//			return err
//		}
//
//		if !info.IsDir() {
//			file, err := os.Open(path)
//			if err != nil {
//				return err
//			}
//			defer file.Close()
//			_, err = io.Copy(writer, file)
//		}
//		return err
//	})
//
//	return err
//}
//
//// 解压缩
//func Unzip(zipFile string, destDir string) error {
//	zipReader, err := zip.OpenReader(zipFile)
//	if err != nil {
//		return err
//	}
//	defer zipReader.Close()
//
//	for _, f := range zipReader.File {
//		fpath := filepath.Join(destDir, f.Name)
//		if f.FileInfo().IsDir() {
//			os.MkdirAll(fpath, os.ModePerm)
//		} else {
//			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
//				return err
//			}
//
//			inFile, err := f.Open()
//			if err != nil {
//				return err
//			}
//			defer inFile.Close()
//
//			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
//			if err != nil {
//				return err
//			}
//			defer outFile.Close()
//
//			_, err = io.Copy(outFile, inFile)
//			if err != nil {
//				return err
//			}
//		}
//	}
//	return nil
//}
