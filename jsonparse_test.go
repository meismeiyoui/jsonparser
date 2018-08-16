package jsonparse

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tmplInStr = `{
	"directive": {
		"header": {
			"namespace": "DNA.Discovery",
			"name": "bl_common3",
			"interfaceVersion": "bl_common2",
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},


		"endpoints": [{
			"endpointId": "bl_common1",
			"friendlyName": "卧室灯"

		}],
		
		"payload": {
			"scope": {
				"token": "some-access-token"
			},
			"options": {
				"enableIntent": false,
				"additionals": {}
			}
		}
	}
}`

var tmplOutStr = `{
	"directive": {
		"header": {
			"namespace": "DNA.Discovery",
			"name": "bl_common3",
			"interfaceVersion": "bl_common2",
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},

		"hello": {
			"endpoints": [{
				"endpointId": "bl_common1",
				"friendlyName": "卧室灯"

			}]
		},

		"payload": {
			"scope": {
				"token": "some-access-token"
			},
			"options": {
				"enableIntent": false,
				"additionals": {}
			}
		}
	}
}`

var inStr = `{
	"directive": {
		"header": {
			"namespace": "DNA.Discovery",
			"name": "inStr",
			"interfaceVersion": "2",
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},


		"endpoints": [{
			"endpointId": "abcd"
		},
		{
			"endpointId": "eeeeee"
		}
		],
		
		"payload": {
			"scope": {
				"type": "ccccccc",
				"token": "some-access-token"
			},
			"options": {
				"enableIntent": false,
				"additionals": {}
			}
		}
	}
}`

var outStr = `{
	"directive": {
		"header": {
			"namespace": "Alex.Discovery",
			"name": "bl_common3",
			"interfaceVersion": "bl_common2",
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},

		"hello": {
			"endpoints": [{
				"endpointId": "abcd"
			},
			{
				"endpointId": "eeeeee"
			}
			]
		},
		
		"payload": {
			"scope": {
				"type": "ccccccc",
				"token": "some-access-token"
			},
			"options": {
				"enableIntent": false,
				"additionals": {}
			}
		}
	}
}`

func TestJsonParse(t *testing.T) {
	jp, err := NewJsonParser(tmplInStr, tmplOutStr)
	assert.Equal(t, err, nil)

	out, ret := jp.GetResult(inStr, outStr)
	assert.Equal(t, ret, Success)
	fmt.Println(PrettyPrintJson(out))

}

func TestGetTag(t *testing.T) {
	fmt.Println(MapKind)
	jp, err := NewTemplateParser([]byte(tmplInStr), []string{"bl_common"})
	if err != nil {
		t.Error("NewJsonParser failed.", err)
	}

	val, ret := jp.GetTag(inStr, "bl_common1")

	assert.Equal(t, ret, Success)
	fmt.Println(val, ret)

}

func TestSetTag(t *testing.T) {

	jp, err := NewTemplateParser([]byte(tmplInStr), []string{"bl_common"})
	if err != nil {
		t.Error("NewJsonParser failed.", err)
	}

	err = jp.SetTag(&inStr, "bl_common1", "new_bl_common_1", "new_bl_common_2")
	contains := strings.Contains(inStr, "new_bl_common_1")
	assert.Equal(t, contains, true)

	fmt.Println(PrettyPrintJson(inStr))
	assert.Equal(t, err, nil)
}
