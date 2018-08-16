package jsonparse

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"simplejson"
)

type Kind uint8

const (
	ElemKind Kind = iota
	MapKind
	ArrayKind
	LastArrayElemKind
)

type nodeStu struct {
	NodeName string
	kind     Kind
}

type pathMapStu map[string][]nodeStu

type templateParserStu struct {
	strTmpl  string
	sjTmpl   *simplejson.Json
	tagPaths map[string][]nodeStu
	prefixes map[string]bool // {prefix: true}
}

func NewTemplateParser(tmpl []byte, prefixes []string) (*templateParserStu, error) {
	j, err := simplejson.NewJson(tmpl)
	if err != nil {
		Traceback(err)
		log.Println("Error json format.")
		return nil, err
	}

	prefixesMap := make(map[string]bool)
	for _, p := range prefixes {
		prefixesMap[p] = true
	}

	tp := &templateParserStu{
		strTmpl:  string(tmpl),
		sjTmpl:   j,
		prefixes: prefixesMap,
		tagPaths: make(map[string][]nodeStu),
	}

	err = tp.traverse()
	if err != nil {
		Error("traverse fail:", err)
	}

	return tp, err
}

func (tp *templateParserStu) traverse() (err error) {
	if err := traverseJson(tp.sjTmpl, []nodeStu{}, tp.tagPaths, tp.prefixes); err != nil {
		Error("traversonJson error", err)
		return err
	}
	tp.markLastArrayElem()
	return
}

func (tp *templateParserStu) GetTag(inStr, tag string, indexes ...int) (val []interface{}, ret int) {
	inJson, err := simplejson.NewJson([]byte(inStr))
	if err != nil {
		// Traceback(err)
		msg := fmt.Sprintf("Cannot parse to simplejson.json %v.", err)
		log.Fatal(msg)
		return nil, ErrJsonFormat
	}

	if path, ok := tp.tagPaths[tag]; ok {
		return getPath(inJson, path, indexes...)
	}

	msg := fmt.Sprintf("Tag:%s Not Found.", tag)
	log.Fatal(msg)
	// Traceback(err)
	return nil, ErrTag
}

func getPath(j *simplejson.Json, path []nodeStu, indexes ...int) (val []interface{}, ret int) {
	if len(path) == 0 {
		return nil, ErrPathIsNone
	}

	_, err := j.Array()
	if err == nil {
		return getInArray(j, path, indexes...)
	}

	_, err = j.Map()
	if err == nil {
		return getInMap(j, path, indexes...)
	}

	Error("j is not map nor array.")
	return nil, ErrData
}

func getInMap(j *simplejson.Json, path []nodeStu, indexes ...int) (val []interface{}, ret int) {
	if v, ret := tryGet(j, path); ret == Success {
		return []interface{}{v}, ret
	}

	nodename := path[0].NodeName
	child := j.Get(nodename)

	return getPath(child, path[1:], indexes...)
}

func tryGet(j *simplejson.Json, path []nodeStu) (val interface{}, ret int) {
	node := path[0].NodeName
	if len(path) == 1 {
		return j.Get(node).Interface(), Success
	}

	return nil, ErrNotReachEnd
}

func getInArray(j *simplejson.Json, path []nodeStu, indexes ...int) (val []interface{}, ret int) {

	a := j.MustArray()

	for i, _ := range a {
		child := j.GetIndex(i)
		v, ret := getPath(child, path, indexes...)
		if ret != Success {
			Error("getPath fail: ", ret)
			return nil, ret
		}

		val = append(val, v...)
	}

	return val, ret
}

func (tp *templateParserStu) GetTagPath(tag string) (path []nodeStu, err error) {
	if path, ok := tp.tagPaths[tag]; ok {
		return path, nil
	}

	msg := fmt.Sprintf("Tag:%s Not Found.", tag)
	return nil, errors.New(msg)
}

func (tp *templateParserStu) SetTag(inStr *string, tag string, val ...interface{}) (err error) {
	inJson, err := parseSimpleJson(*inStr)
	if err != nil {
		return err
	}

	path, err := tp.GetTagPath(tag)
	if err != nil {
		return err
	}

	setPath(inJson, path, val...)
	// inJson.SetPath(path, val)
	bys, err := inJson.Encode()
	if err != nil {
		Traceback(err)
		return err
	}

	*inStr = string(bys)
	return err
}

func setPath(j *simplejson.Json, path []nodeStu, val ...interface{}) (ret int) {
	if len(path) == 0 {
		Error("path is none.")
		return ErrPathIsNone
	}

	if len(val) == 0 {
		Error("val in none.")
		return ErrData
	}

	_, err := j.Array()
	if err == nil {
		ret = setInArray(j, path, val...)
		return
	}

	_, err = j.Map()
	if err == nil {
		ret = setInMap(j, path, val...)
		if ret != Success {
			Error("setInMap fail", ret)
		}
		return
	}

	Error("j is not a map nor array. what?")
	return ErrData
}

func setInArray(j *simplejson.Json, path []nodeStu, val ...interface{}) (ret int) {
	a := j.MustArray()

	cn := path[0]
	for i, _ := range a {
		ret = setPath(j.GetIndex(i), path, val...)

		Roadmap("LastArrayKind====", cn.NodeName, cn.kind)
		if cn.kind == LastArrayElemKind {
			val = val[1:]
		}

		if ret != Success {
			Error("setPath fail:", ret)
			return ret
		}
	}

	return
}

func setInMap(j *simplejson.Json, path []nodeStu, val ...interface{}) (ret int) {
	ret = trySet(j, path, val[0])
	if ret == Success {
		return
	}

	node := path[0].NodeName
	child := j.Get(node)
	Roadmap("mappath", path, child.Interface())
	setPath(child, path[1:], val...)

	return Success
}

func trySet(j *simplejson.Json, path []nodeStu, val interface{}) (ret int) {
	node := path[0].NodeName
	if len(path) == 1 {
		j.Set(node, val)
		return Success
	}

	return ErrNotReachEnd
}

func traverseJson(sj *simplejson.Json, path []nodeStu, tagPaths pathMapStu, prefixes map[string]bool) (err error) {
	sjm, err := sj.Map()
	if err != nil {
		return err
	}

	for k, tag := range sjm {
		var node = nodeStu{}
		node.NodeName = k
		path = append(path, node)
		switch tag.(type) {
		case string:
			path[len(path)-1].kind = ElemKind
			tagStr := tag.(string)
			if IsTag(tagStr, prefixes) {
				// !! errors: path is unstable(temporary)
				tagPaths[tagStr] = make([]nodeStu, len(path))
				copy(tagPaths[tagStr], path)
			}

		case map[string]interface{}:
			path[len(path)-1].kind = MapKind
			traverseJson(sj.Get(k), path, tagPaths, prefixes)

		case []interface{}:
			path[len(path)-1].kind = ArrayKind
			traverseJson(sj.Get(k).GetIndex(0), path, tagPaths, prefixes)

		default:

		}
		path = path[:len(path)-1]
	}

	return nil
}

func (tp *templateParserStu) markLastArrayElem() {
	for _, path := range tp.tagPaths {
		for i := len(path) - 1; i >= 0; i-- {
			if path[i].kind == ArrayKind {
				if i < len(path)-1 {
					path[i+1].kind = LastArrayElemKind
				}
				break
			}
		}
	}
}

func IsTag(v string, prefixes map[string]bool) bool {
	regTag := regexp.MustCompile(`(\w+?)[_]*(\d+)`)

	match := regTag.MatchString(v)
	if !match {
		return false
	}

	ss := regTag.FindStringSubmatch(v)
	if _, ok := prefixes[ss[1]]; ok {
		return true
	}

	return false
}

func parseSimpleJson(inStr string) (*simplejson.Json, error) {
	inJson, err := simplejson.NewJson([]byte(inStr))
	if err != nil {
		msg := fmt.Sprintf("Cannot parse to simplejson.json %v.", err)
		log.Fatal(msg)
		return nil, err
	}

	return inJson, nil
}
