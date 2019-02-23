package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tealeg/xlsx"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

func main() {
	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	outputFileName := flag.String("output", "", "The output file name.")
	flag.Parse()
	if *outputFileName == "" {
		flag.Usage()
		os.Exit(1)
	}

	if len(*outputFileName) >= 5 {
		if (*outputFileName)[len(*outputFileName)-5:] == ".xlsx" {
		} else {
			*outputFileName += ".xlsx"
		}
	} else {
		*outputFileName += ".xlsx"
	}

	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	//获取的是工作目录，不一定是可执行文件所在目录
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, fi := range rd {
		name := fi.Name()
		if len(name) >= 4 {
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
							row = sheet.AddRow()
							cell = row.AddCell()
							cell.SetString(name[:len(name)-4])
							cell = row.AddCell()
							cell.SetString(strconv.FormatInt(j.TimeCost, 10))
						}
					}
				}

				if err := scanner.Err(); err != nil {
					panic(err)
				}

			}
		}
	}

	err = file.Save(*outputFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Output at: ",*outputFileName)
}
