package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/sycdtk/gotools/config"
	"github.com/sycdtk/gotools/csv"
	"github.com/sycdtk/gotools/logger"
)

//重名默认补全长度
const DEFAULT_SIZE = 2

type Record struct {
	X, Y int
}

func main() {

	flag.Usage = usage

	configFile := flag.String("c", "config.conf", "配置文件名称")

	fileName := flag.String("f", "data.csv", "数据文件名称")

	opType := flag.String("o", "", `r: replace，值匹配r列后，写入w列，若无w列则直接写回r列；
	t: transfer，转换数据，正则匹配r列后，写入w列，若无w列则直接写回r列；
	d: duplicate removal，数据列去重，重复数据追加_01、_02；
	e: exchange，数据列交换，r列与w列交换；
	x: extract，正则提取数据，r列正则匹配分组，写入w列。仅适用数据集合s的key值；
	tr: trim，去除列前后空白字符；
	f: fill，对于r列数据相等的数据，补全其w列数据，不为空的数据将被补全；
	re: 正则测试`)

	dataSet := flag.String("s", "", "配置文件中的数据组")

	readCol := flag.Int("r", -1, "替换列，0为第一列")

	writeCol := flag.Int("w", -1, "匹配第一列后需要写入数据的列，0为第一列")

	re1 := flag.String("re1", "", "正则表达式")
	re2 := flag.String("re2", "", "正则匹配测试字符串")

	flag.Parse()

	config.Load(*configFile)

	if *opType != "r" && *opType != "t" && *opType != "d" && *opType != "e" && *opType != "x" && *opType != "f" && *opType != "tr" && *opType != "re" {
		fmt.Println("Error：o参数输入错误！")
		os.Exit(-1)
	}

	if (*opType == "r" || *opType == "t") && (len(*dataSet) == 0 || *readCol < 0) {
		fmt.Println("Error：操作类型为r或t时，参数s和r为必须参数！")
		os.Exit(-1)
	}

	if (*opType == "e" || *opType == "f") && (*readCol < 0 || *writeCol < 0) {
		fmt.Println("Error：操作类型为e或f时，参数r和w为必须参数！")
		os.Exit(-1)
	}

	if *opType == "x" && (len(*dataSet) == 0 || *readCol < 0 || *writeCol < 0) {
		fmt.Println("Error：操作类型为x时，参数s、r和w为必须参数！")
		os.Exit(-1)
	}

	if *opType == "tr" && *readCol < 0 {
		fmt.Println("Error：操作类型为tr时，参数r为必须参数！")
		os.Exit(-1)
	}

	if *opType == "re" && (len(*re1) == 0 || len(*re2) == 0) {
		fmt.Println("Error：操作类型为re时，参数re1和re2为必须参数！")
		os.Exit(-1)
	}

	if *opType == "re" {
		reTest(*re1, *re2)
	} else {
		dataFile := csv.NewCSV(*fileName)

		dataFile.Reader()

		if *opType == "r" {
			replace(dataFile, *dataSet, *readCol, *writeCol)
		} else if *opType == "t" {
			transfer(dataFile, *dataSet, *readCol, *writeCol)
		} else if *opType == "d" {
			duplicateRemoval(dataFile, *readCol)
		} else if *opType == "e" {
			exchange(dataFile, *readCol, *writeCol)
		} else if *opType == "x" {
			extract(dataFile, *dataSet, *readCol, *writeCol)
		} else if *opType == "f" {
			fill(dataFile, *readCol, *writeCol)
		} else if *opType == "tr" {
			trim(dataFile, *readCol)
		}

		dataFile.Writer(dataFile.Datas, false)

		fmt.Println("Done!")
	}

}

/*
对于值相同的列，补全数据
*/
func fill(dataFile *csv.CSV, readCol, wirteCol int) {
	//数据集合
	dataMap := map[string][]*Record{}

	for r, data := range dataFile.Datas {
		//数据分组
		dataMap[data[readCol]] = append(dataMap[data[readCol]], &Record{r, readCol})
	}

	//数据补全
	for _, v := range dataMap {
		logger.Debug(len(v))
		//存在重复数据
		if len(v) > 1 {
			fillStr := ""
			//提取补全数据
			for _, d := range v {
				if len(dataFile.Datas[d.X][wirteCol]) > 0 {
					//超过2条记录的，以最后一条记录为准
					fillStr = dataFile.Datas[d.X][wirteCol]
				}
			}

			logger.Debug(fillStr)

			//补全数据
			for _, d := range v {
				dataFile.Datas[d.X][wirteCol] = fillStr
			}
		}
	}
}

/*
去除开头结尾空白符
*/
func trim(dataFile *csv.CSV, readCol int) {
	for _, data := range dataFile.Datas {
		data[readCol] = strings.TrimSpace(data[readCol])
	}
}

/*
正则测试
*/
func reTest(re1, re2 string) {
	re, _ := regexp.Compile(re1)
	fmt.Println(re1)
	fmt.Println(re2)
	fmt.Println(re.MatchString(re2))
	if re.MatchString(re2) {
		fmt.Println(re.FindStringSubmatch(re2)[1])
	}

}

/*
正则提取数据
*/
func extract(dataFile *csv.CSV, dataSet string, readCol, writeCol int) {
	for key, _ := range config.GetNode(dataSet) {
		logger.Debug(key)
		re, _ := regexp.Compile(key)
		for r, _ := range dataFile.Datas {
			data := dataFile.Datas[r][readCol]
			if re.MatchString(data) {
				if writeCol == -1 {
					dataFile.Datas[r][readCol] = re.FindStringSubmatch(data)[1]
				} else {
					dataFile.Datas[r][writeCol] = re.FindStringSubmatch(data)[1]
				}
			}
		}
	}
}

/*
数据列交换
*/
func exchange(dataFile *csv.CSV, readCol, writeCol int) {
	for _, data := range dataFile.Datas {
		data[readCol], data[writeCol] = data[writeCol], data[readCol]
	}
}

/*
数据列去重
*/
func duplicateRemoval(dataFile *csv.CSV, readCol int) {
	//数据集合
	dataMap := map[string][]*Record{}

	for r, data := range dataFile.Datas {
		//数据预处理，不能包含单个的"\\"
		if strings.Contains(data[readCol], "\\") {
			data[readCol] = strings.Replace(data[readCol], "\\", " ", -1)
		}
		//数据分组
		dataMap[data[readCol]] = append(dataMap[data[readCol]], &Record{r, readCol})
	}

	//重复数据替换
	for _, v := range dataMap {
		//存在重复数据
		if len(v) > 1 {
			for i, d := range v {
				dataFile.Datas[d.X][d.Y] = dataFile.Datas[d.X][d.Y] + numToStr(i+1, DEFAULT_SIZE)
			}
		}
	}

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
					if value == "nil" {
						dataFile.Datas[r][readCol] = ""
					} else {
						dataFile.Datas[r][readCol] = value
					}
				} else {
					if value == "nil" {
						dataFile.Datas[r][writeCol] = ""
					} else {
						dataFile.Datas[r][writeCol] = value
					}
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
				if data == "nil" {
					logger.Debug(dataFile.Datas[r][readCol], data)
					dataFile.Datas[r][readCol] = ""
				} else {
					dataFile.Datas[r][readCol] = data
				}
			} else {
				if data == "nil" {
					logger.Debug(dataFile.Datas[r][readCol], data)
					dataFile.Datas[r][writeCol] = ""
				} else {
					dataFile.Datas[r][writeCol] = data
				}
			}
		}
	}
}

func numToStr(num, size int) string {

	ns := strconv.Itoa(num)
	l := size - len(ns)
	for i := 0; i < l; i++ {
		ns = "0" + ns
	}

	return "_" + ns
}

func usage() {
	fmt.Println(" CSV Tool V1.0 ;)\nhttps://github.com/sycdtk/csv_trans\n")
	flag.PrintDefaults()
}
