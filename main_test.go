package main

import (
	"regexp"
	"testing"
)

func TestRe(t *testing.T) {
	re, _ := regexp.Compile(`dell \S+ r520`)
	t.Log(re.ReplaceAllString("dell poweredge asd r520", "$1"))
	t.Log(re.MatchString("dell poweredge asd r520"))
}
