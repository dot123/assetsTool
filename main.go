package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Project struct {
	PackageUrl        string
	RemoteManifestUrl string
	RemoteVersionUrl  string
	Version           string
	Assets            map[string]struct {
		Size int
		Md5  string
	}
}

func loadJson(path string, dist interface{}) (err error) {
	var content []byte
	if content, err = ioutil.ReadFile(path); err == nil {
		err = json.Unmarshal(content, dist)
	}
	return err
}

func main() {
	var (
		path1   string
		path2   string
		resPath string
	)

	flag.StringVar(&path1, "oldProject", "./project.manifest", "旧的project.manifest")
	flag.StringVar(&path2, "newProject", "./assets/project.manifest", "新的project.manifest")
	flag.StringVar(&resPath, "resPath", "./assets/", "资源路径")
	flag.Parse()

	var (
		project1 Project
		project2 Project
	)

	if err := loadJson(path1, &project1); err != nil {
		fmt.Println(err)
	}

	if err := loadJson(path2, &project2); err != nil {
		fmt.Println(err)
	}

	var fileList []string

	for key, b := range project2.Assets {
		a, ok := project1.Assets[key]
		if ok {
			if a.Md5 != b.Md5 {
				fileList = append(fileList, key)
			}
		} else {
			fileList = append(fileList, key)
		}
	}

	sort.Strings(fileList)

	destPath := resPath + "diff/"
	os.RemoveAll(destPath)
	createDir(filepath.Dir(destPath))

	for k, v := range fileList {
		bytes, _ := ioutil.ReadFile(resPath + v)
		if bytes != nil {
			createDir(filepath.Dir(destPath + v))
			err := ioutil.WriteFile(destPath+v, bytes, 0666)
			checkError(err)
		}
		fmt.Println(k, v)
	}

	err := ioutil.WriteFile(resPath+"fileList.txt", list2Str(fileList), 0666)
	checkError(err)

	time.Sleep(time.Millisecond * 1000)
}

//创建文件夹
func createDir(dir string) error {
	exist, err := pathExists(dir)
	checkError(err)

	if exist {

	} else {
		//创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		checkError(err)
	}
	return nil
}

//判断文件夹是否存在
func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func checkError(err error) {
	if err != nil {
		log.Fatalf("%v\n", err)
		return
	}
}

func list2Str(list []string) []byte {
	b, err := json.MarshalIndent(list, "", "")
	if err != nil {
		checkError(err)
		return nil
	}
	return b
}
