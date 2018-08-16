package jsonparse

import (
	"fmt"
	"regexp"
	"simplejson"
	"testing"

	"github.com/stretchr/testify/assert"
)

// var j = `{
// 	"a": "hello"
// }`
var j = "[1,2]"

func TestMust(t *testing.T) {
	sj, err := simplejson.NewJson([]byte(j))
	if err != nil {
		panic(err)
	}

	m := sj.MustMap()
	em := make(map[string]interface{})
	assert.Equal(t, m, em)
	fmt.Println(m, em)
}

func TestAssert(t *testing.T) {
	var a interface{} = "123"
	astr, ok := a.(string)
	fmt.Println(astr, ok)
}

func TestRegexp(t *testing.T) {
	reg := regexp.MustCompile(`(\w+?)[_]*(\d+)`)
	// reg := regexp.MustCompile(`\W+`)

	content := "bl_common_1"
	b := reg.MatchString(content)
	assert.Equal(t, b, true)

	ss := reg.FindStringSubmatch(content)
	fmt.Println(ss)

	assert.Equal(t, len(ss), 3)
	assert.Equal(t, ss[0], content)
	assert.Equal(t, ss[1], "bl_common")
	assert.Equal(t, ss[2], "1")
}

func TestTraceBack(t *testing.T) {
	Traceback(nil)
}

func TestLoadConfig(t *testing.T) {
	loadConfig()
	fmt.Printf("%+v\n", sysConfig)
}

func TestStr(t *testing.T) {
	var s = "123"
	var a = s
	fmt.Println(a, s)
}
