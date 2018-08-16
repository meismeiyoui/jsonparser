package jsonparse

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"simplejson"
)

type templateParserStu struct {
	strTmpl  string
	sjTmpl   *simplejson.Json
	tagPaths map[string][]string
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
		tagPaths: make(map[string][]string),
	}

	tp.traverse()
	return tp, nil
}

func (tp *templateParserStu) traverse() {
	traverseJson(tp.sjTmpl, []string{}, tp.tagPaths, tp.prefixes)
	return
}

func (tp *templateParserStu) GetTag(inStr, tag string) (val interface{}, err error) {
	inJson, err := simplejson.NewJson([]byte(inStr))
	if err != nil {
		// Traceback(err)
		msg := fmt.Sprintf("Cannot parse to simplejson.json %v.", err)
		log.Fatal(msg)
		return nil, err
	}

	if tagPaths, ok := tp.tagPaths[tag]; ok {
		return inJson.GetPath(tagPaths...), nil
	}

	msg := fmt.Sprintf("Tag:%s Not Found.", tag)
	log.Fatal(msg)
	// Traceback(err)
	return nil, errors.New(msg)
}

func (tp *templateParserStu) GetTagPath(tag string) (path []string, err error) {
	if path, ok := tp.tagPaths[tag]; ok {
		return path, nil
	}

	msg := fmt.Sprintf("Tag:%s Not Found.", tag)
	return nil, errors.New(msg)
}

func (tp *templateParserStu) SetTag(inStr *string, tag string, val interface{}) (err error) {
	inJson, err := parseSimpleJson(*inStr)
	if err != nil {
		return err
	}

	path, err := tp.GetTagPath(tag)
	if err != nil {
		return err
	}

	inJson.SetPath(path, val)
	bys, err := inJson.Encode()
	if err != nil {
		Traceback(err)
		return err
	}

	*inStr = string(bys)
	return err
}

func (tp *templateParserStu) GetAllTag(inStr string) (m map[string]interface{}, err error) {
	inJson, err := parseSimpleJson(inStr)
	if err != nil {
		return nil, err
	}

	m = make(map[string]interface{})
	for k, path := range tp.tagPaths {
		val := inJson.GetPath(path...)
		if val != inJson {
			m[k] = val.Interface()
		} else {
			msg := fmt.Sprintf("GetPath fail: %s", path)
			log.Fatal(msg)
		}
	}
	return
}

func traverseJson(sj *simplejson.Json, path []string, tagPaths map[string][]string, prefixes map[string]bool) (err error) {
	sjm, err := sj.Map()
	if err != nil {
		return err
	}

	for k, tag := range sjm {
		path = append(path, k)
		switch tag.(type) {
		case string:
			tagStr := tag.(string)
			if IsTag(tagStr, prefixes) {
				// errors: path is unstable(temporary)
				tagPaths[tagStr] = make([]string, len(path))
				copy(tagPaths[tagStr], path)
			}
		case map[string]interface{}:
			traverseJson(sj.Get(k), path, tagPaths, prefixes)
		default:

		}
		path = path[:len(path)-1]
	}

	return nil
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
