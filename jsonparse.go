package jsonparse

type jsonParserStu struct {
	tmplIn  *templateParserStu
	tmplOut *templateParserStu
	// inStr   string
	rules *[]RuleStu
}

func init() {
	loadConfig()
}

func NewJsonParser(tmplIn, tmplOut string) (*jsonParserStu, error) {
	var inPrefixes, outPrefixes []string

	inPrefixes = append(inPrefixes, sysConfig.Prefixes.Common...)
	inPrefixes = append(inPrefixes, sysConfig.Prefixes.In...)
	tmplInParser, err := NewTemplateParser([]byte(tmplIn), inPrefixes)
	if err != nil {
		return nil, err
	}

	outPrefixes = append(outPrefixes, sysConfig.Prefixes.Common...)
	outPrefixes = append(outPrefixes, sysConfig.Prefixes.Out...)
	tmplOutParser, err := NewTemplateParser([]byte(tmplOut), outPrefixes)
	if err != nil {
		return nil, err
	}

	jp := &jsonParserStu{
		tmplIn:  tmplInParser,
		tmplOut: tmplOutParser,
		// inStr:   inStr,
	}

	return jp, nil
}

type RuleStu struct {
	ruleStr string
}

func (r *RuleStu) applyRule(inVal string, outVal *string) {
	return
}

func (jp *jsonParserStu) GetResult(in string, out string) (s string, ret int) {
	s, ret = jp.parseJson(in, out)
	return
}

func (jp *jsonParserStu) parseJson(in, out string) (s string, ret int) {
	// 1. common rule

	var tmp = out
	for tag, _ := range jp.tmplOut.tagPaths {
		valIn, ret := jp.tmplIn.GetTag(in, tag)
		if ret != Success {
			Error(ret)
			continue
		}
		jp.tmplOut.SetTag(&tmp, tag, valIn...)
	}

	// 2. rules

	s = tmp
	return
}
