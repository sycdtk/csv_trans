package main

import (
	"regexp"
	"strings"
	"testing"
)

func TestRe(t *testing.T) {

	ww := " 1.2.3.4 "
	t.Log("|" + strings.TrimSpace(ww) + "|")

	qq := "去玩儿\\阿斯蒂芬 asdf\\asdf"

	if strings.Contains(qq, "\\") {
		qq = strings.Replace(qq, "\\", " ", -1)
	}
	t.Log(qq)

	a := "aCentOS(2.6.18-194.el5PAE)(i686)asdfxx "
	re, _ := regexp.Compile(`CentOS\((.*?)\).*`)
	b := re.FindStringSubmatch(a)

	if re.MatchString(a) {
		for i, d := range b {
			t.Log(i, d)
		}
	}

	t.Log(re.MatchString("dell poweredge asd r520"))

	data := map[string][]*Record{}

	data["aa"] = append(data["aa"], &Record{1, 2})
	data["aa"] = append(data["aa"], &Record{3, 4})

	t.Log(data["aa"][1].X)

}
