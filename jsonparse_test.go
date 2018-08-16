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
			"friendlyName": "卧室灯",
			"isReachable": true,
			"description": "由BroadLink生产的灯",
			"manufacturerName": "Sample Manufacturer",
			"icon": "产品图片URL",
			"brand": "品牌",
			"displayCategories": [
				"LIGHT"
			],
			"roomName": "用户设置的房间名称",
			"cookie": {
				"familyId": "家庭id",
				"familyName": "用户设置的家庭名称",
				"extraDetail1": "某些设备可能会用到这个cookie，需要在控制时原样返回",
				"extraDetail2": "某些设备可能会用到这个cookie，需要在控制时原样返回",
				"extraDetail3": "某些设备可能会用到这个cookie，需要在控制时原样返回",
				"extraDetail4": "某些设备可能会用到这个cookie，需要在控制时原样返回"
			},
			"capabilities": [{
				"type": "DNAInterface",
				"interface": "DNA.PowerControl",
				"version": "2",
				"properties": {
					"supported": [{
						"name": "powerState"
					}],
					"proactivelyReported": true,
					"retrievable": true
				}
			}]
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
			"name": "bl_common3",
			"interfaceVersion": "bl_common2",
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

func TestJsonParse(t *testing.T) {
	jp, err := NewJsonParser(tmplInStr, tmplOutStr, inStr)
	assert.Equal(t, err, nil)

	out, ret := jp.GetResult()
	assert.Equal(t, ret, Success)
	fmt.Println(out)

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
