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
		if (*outputFileName)[len(*outputFileName)-5:] != ".xlsx" {
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

	row = sheet.AddRow()
	cell = row.AddCell()
	cell.SetString("filename")
	cell = row.AddCell()
	cell.SetString("TransactionHash")
	cell = row.AddCell()
	cell.SetString("DataFirst4Byte")
	cell = row.AddCell()
	cell.SetString("TransactionBegin")
	cell = row.AddCell()
	cell.SetString("TransactionEnd")
	cell = row.AddCell()
	cell.SetString("TransactionTime")
	cell = row.AddCell()
	cell.SetString("NewTransaction")
	cell = row.AddCell()
	cell.SetString("BlockSeal")
	cell = row.AddCell()
	cell.SetString("TransactionLatency")

	for _, fi := range rd {
		json_list := make([]map[string]interface{}, 0)
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
					var obj map[string]interface{}
					err = json.Unmarshal([]byte(scanner.Text()), &obj)
					if err != nil {

					} else {
						json_list = append(json_list, obj)
						//if obj.DataFirst4ByteHex == "98d69d92" {
						//	row = sheet.AddRow()
						//	cell = row.AddCell()
						//	cell.SetString(name[:len(name)-4])
						//	cell = row.AddCell()
						//	cell.SetString(strconv.FormatInt(j.TimeCost, 10))
						//}
					}
				}

				if err := scanner.Err(); err != nil {
					panic(err)
				}

				for _, newTx := range json_list {
					if newTx["Type"] == "NewTransaction" {
						hash := newTx["TransactionHash"]
						row = sheet.AddRow()
						cell = row.AddCell()
						cell.SetString(name[:len(name)-4])
						cell = row.AddCell()
						cell.SetString(hash.(string))
						//find transaction begin
						var txBegin map[string]interface{}
						found := false
						for _, obj2 := range json_list {
							if obj2["Type"] == "TransactionBegin" {
								if obj2["TransactionHash"] == hash {
									//find
									found = true
									txBegin = obj2
									break
								}
							}
						}
						if !found {
							fmt.Println("Warning: unmatched transaction", hash)
							continue
						}

						//find transaction end
						var txEnd map[string]interface{}
						found = false
						for _, obj2 := range json_list {
							if obj2["Type"] == "TransactionEnd" {
								if obj2["TransactionHash"] == hash {
									//find
									found = true
									txEnd = obj2
									break
								}
							}
						}
						if !found {
							fmt.Println("Warning: unmatched transaction", hash)
							continue
						}

						cell = row.AddCell()
						cell.SetString(fmt.Sprint(txEnd["DataFirst4Byte"]))
						cell = row.AddCell()
						cell.SetFloat(txBegin["Timestamp"].(float64))
						cell = row.AddCell()
						cell.SetFloat(txEnd["Timestamp"].(float64))
						cell = row.AddCell()
						cell.SetFloat(txEnd["Timestamp"].(float64) - txBegin["Timestamp"].(float64))

						//find BlockSeal
						var BlockSeal map[string]interface{}
						found = false
						for _, obj2 := range json_list {
							if obj2["Type"] == "BlockSeal" {
								for _, obj3 := range obj2["TransactionHashs"].([]interface{}) {
									if obj3 == hash {
										found = true
										BlockSeal = obj2
										break
									}
								}
							}
						}
						cell = row.AddCell()
						cell.SetString(fmt.Sprint(newTx["Timestamp"]))
						cell = row.AddCell()
						cell.SetString(fmt.Sprint(BlockSeal["Timestamp"]))
						cell = row.AddCell()
						cell.SetFloat(BlockSeal["Timestamp"].(float64) - newTx["Timestamp"].(float64))
						//reserved for formula

						if !found {
							fmt.Println("Warning: unmatched transaction", hash)
							continue
						}

					}

				}

			}
		}
	}

	err = file.Save(*outputFileName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Output at: ", *outputFileName)
}
