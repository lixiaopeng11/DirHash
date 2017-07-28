package main

import (
	"fmt"
	"log"
	"os"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"flag"
	"time"
	"runtime"
	"github.com/ryanuber/go-glob"

)

func listAll(path string, logFile *os.File, filter string) {
	// 遍历目录
	files, _ := ioutil.ReadDir(path)
	for _, fi := range files {
		
		//过滤文件或目录
		if (filter != ""){
			//ok :=strings.Contains(path + "/" + fi.Name(),filter)
			ok :=glob.Glob(filter, path + "/" + fi.Name())
			fmt.Println(filter, path + "/" + fi.Name(),ok)
			if (ok){
				continue
			}
		}


		//fs := fi.NewFileSearch()
		if fi.IsDir() {
			//若是目录，遍历此目录的文件
			go listAll(path + "/" + fi.Name(), logFile, filter)

		} else {
			//println(path + "/" + fi.Name())
			go getInfoFile( path + "/" + fi.Name(), fi, logFile )
		}
	}
}

func getInfoFile(path string, info os.FileInfo, logFile *os.File )  {
	//获取文件信息
	if (!info.IsDir()){
		h := sha1.New()
		file, err := os.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		_, error := io.Copy(h, file)
		if error != nil {
			log.Fatal(err)
		}
		//fmt.Printf("%s,%d,%x\n",path,info.Size(), h.Sum(nil))
		result :=fmt.Sprintf("%s,%x,%d\n",path, h.Sum(nil), info.Size())
		logFile.WriteString(result)
	}

}



func main() {
	
	//dir    := flag.String("root", ".", "遍历目录")
	var dir string
	flag.StringVar(&dir, "root", ".", "遍历目录")
	var filter string
	flag.StringVar(&filter, "filter", "", "需要过滤的目录或文件")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())  // 读取CPU数量 最大限度利用上多核性能
	
	
	fileName := "file.log"
	// 判断文件是否存在， 若文件存在，则删除
	_,  err   := os.Stat(fileName)

	if os.IsExist(err) {
		err = os.Remove(fileName)
		if err != nil {
			//如果删除失败则输出 file remove Error!
			fmt.Println("file remove Error!")
			//输出错误详细信息
			fmt.Printf("%s", err)
		}
	}

	
	logFile,err  :=  os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY, 0600)
	defer logFile.Close()
	
	if err != nil {
		log.Fatalln("open file error !")
	}
	log.SetFlags(log.Lshortfile)
	listAll( dir, logFile, filter)
	time.Sleep(1e9)

	
}