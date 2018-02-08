package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/sycdtk/gotools/config"
	"github.com/sycdtk/gotools/csv"
	"github.com/sycdtk/gotools/logger"
)

func main() {

	flag.Usage = usage

	fileName := flag.String("f", "data.csv", "文件名称")

	opType := flag.String("o", "", `r:[replace] ，值匹配r列后，写入w列，若无w列则直接写回r列；
	t:[transfer] 转换数据，正则匹配r列后，写入w列，若无w列则直接写回r列`)

	dataSet := flag.String("s", "", "配置文件中的数据组")

	readCol := flag.Int("r", -1, "替换列，0为第一列")

	writeCol := flag.Int("w", -1, "匹配第一列后需要写入数据的列，0为第一列")

	flag.Parse()

	if *opType != "r" && *opType != "t" {
		logger.Info("Error：o参数输入错误！")
		os.Exit(-1)
	}

	if len(*dataSet) == 0 {
		logger.Info("Error：s参数为必须参数！")
		os.Exit(-1)
	}

	if *readCol < 0 {
		logger.Info("Error：r参数为必须参数！")
		os.Exit(-1)
	}

	dataFile := csv.NewCSV(*fileName)

	dataFile.Reader()

	if *opType == "r" {
		replace(dataFile, *dataSet, *readCol, *writeCol)
	} else if *opType == "t" {
		transfer(dataFile, *dataSet, *readCol, *writeCol)
	}

	dataFile.Writer(dataFile.Datas, false)

	fmt.Println("Done!")
}

/*
数据转换
*/
func transfer(dataFile *csv.CSV, dataSet string, readCol, writeCol int) {
	//正则匹配
	for key, value := range config.GetNode(dataSet) {
		re, _ := regexp.Compile(key)
		for r, _ := range dataFile.Datas {
			if re.MatchString(dataFile.Datas[r][readCol]) {
				if writeCol == -1 {
					dataFile.Datas[r][readCol] = value
				} else {
					dataFile.Datas[r][writeCol] = value
				}
			}
		}
	}
}

/*
数据替换
*/
func replace(dataFile *csv.CSV, dataSet string, readCol, writeCol int) {
	for r, _ := range dataFile.Datas {
		if data := config.Read(dataSet, dataFile.Datas[r][readCol]); len(data) > 0 {
			if writeCol == -1 {
				dataFile.Datas[r][readCol] = data
			} else {
				dataFile.Datas[r][writeCol] = data
			}
		}
	}
}

func usage() {
	fmt.Println(" CSV Tool V1.0 ;)\n")
	flag.PrintDefaults()
}
