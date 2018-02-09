package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/sycdtk/gotools/config"
	"github.com/sycdtk/gotools/csv"
)

//重名默认补全长度
const DEFAULT_SIZE = 2

type Record struct {
	X, Y int
}

func main() {

	flag.Usage = usage

	fileName := flag.String("f", "data.csv", "文件名称")

	opType := flag.String("o", "", `r: replace，值匹配r列后，写入w列，若无w列则直接写回r列；
	t: transfer，转换数据，正则匹配r列后，写入w列，若无w列则直接写回r列；
	d: duplicate removal，数据列去重，重复数据追加_01、_02；
	e: exchange，数据列交换，列1与列3交换；
	re: 正则测试`)

	dataSet := flag.String("s", "", "配置文件中的数据组")

	readCol := flag.Int("r", -1, "替换列，0为第一列")

	writeCol := flag.Int("w", -1, "匹配第一列后需要写入数据的列，0为第一列")

	re1 := flag.String("re1", "", "正则表达式")
	re2 := flag.String("re2", "", "正则匹配测试字符串")

	flag.Parse()

	if *opType != "r" && *opType != "t" && *opType != "d" && *opType != "e" && *opType != "re" {
		fmt.Println("Error：o参数输入错误！")
		os.Exit(-1)
	}

	if (*opType == "r" || *opType == "t") && (len(*dataSet) == 0 || *readCol < 0) {
		fmt.Println("Error：操作类型为r或t时，参数s和r为必须参数！")
		os.Exit(-1)
	}

	if *opType == "e" && (*readCol < 0 || *writeCol < 0) {
		fmt.Println("Error：操作类型为e时，参数r和w为必须参数！")
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
		}

		dataFile.Writer(dataFile.Datas, false)

		fmt.Println("Done!")
	}

}

/*
正则测试
*/
func reTest(re1, re2 string) {
	re, _ := regexp.Compile(re1)
	fmt.Println(re1, re2, re.MatchString(re2))
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

	//数据分组
	for r, data := range dataFile.Datas {
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

func numToStr(num, size int) string {

	ns := strconv.Itoa(num)
	l := size - len(ns)
	for i := 0; i < l; i++ {
		ns = "0" + ns
	}

	return "_" + ns
}

func usage() {
	fmt.Println(" CSV Tool V1.0 ;)\n")
	flag.PrintDefaults()
}
