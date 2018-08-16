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
			"name": "Discover",
			"interfaceVersion": "bl_common2",
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},
		"payload": {
			"scope": {
				"type": "bl_common1",
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
			"namespace": "Alexa.Discovery",
			"name": "Discover",
			"interfaceVersion": "bl_common2",
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},
		"payload": {
			"scope": {
				"type": "bl_common1",
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
			"name": "Discover",
			"interfaceVersion": 2,
			"messageId": "1bd5d003-31b9-476f-ad03-71d471922820"
		},
		"payload": {
			"scope": {
				"type": "===BearerToken======",
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
	jp, err := NewJsonParser(tmplInStr, tmplOutStr, inStr)
	assert.Equal(t, err, nil)

	out, err := jp.GetResult()
	assert.Equal(t, err, nil)
	fmt.Println(out)

}

func TestGetTag(t *testing.T) {
	jp, err := NewTemplateParser([]byte(tmplInStr), []string{"bl_common"})
	if err != nil {
		t.Error("NewJsonParser failed.", err)
	}

	tags, err := jp.GetAllTag(inStr)
	fmt.Println(jp)

	assert.Equal(t, err, nil)
	fmt.Println(tags["bl_common1"], "===BearerToken======")
	fmt.Println(tags["bl_common2"], 2)
}

func TestSetTag(t *testing.T) {

	jp, err := NewTemplateParser([]byte(tmplInStr), []string{"bl_common"})
	if err != nil {
		t.Error("NewJsonParser failed.", err)
	}

	err = jp.SetTag(&inStr, "bl_common1", "new_bl_common_1")
	contains := strings.Contains(inStr, "new_bl_common_1")
	assert.Equal(t, contains, true)

	err = jp.SetTag(&inStr, "bl_common1", []int{1, 2, 3})
	contains = strings.Contains(inStr, "new_bl_common_1")
	assert.Equal(t, contains, false)

	fmt.Println(inStr)

	assert.Equal(t, err, nil)
}
