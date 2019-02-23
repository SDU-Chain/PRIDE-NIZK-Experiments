package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func main() {
	//获取的是工作目录，不一定是可执行文件所在目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, fi := range rd {
		name := fi.Name()
		if name[len(name)-4:] == ".txt" {

			file, err := os.Open(path.Join(dir, name))
			if err != nil {
				panic(err)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				var j struct {
					From              string
					To                string
					TimeCost          int64
					Timestamp         int64
					DataFirst4ByteHex string
				}
				err = json.Unmarshal([]byte(scanner.Text()), &j)
				if err != nil {

				} else {
					if j.DataFirst4ByteHex == "98d69d92" {
						fmt.Println(name[:len(name)-4], ",", j.TimeCost,",")
					}
				}
			}

			if err := scanner.Err(); err != nil {
				panic(err)
			}
 

		}
	}
}
