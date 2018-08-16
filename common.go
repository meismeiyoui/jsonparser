package jsonparse

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
)

func Traceback(err error) {
	if EnableTraceBack {
		log.Println("[TRACEBACK]<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<")
		log.Println(err)
		debug.PrintStack()
		log.Println("[TRACEBACK]>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")
	}
}

func PrettyPrintJson(srcdata interface{}) (ret string, err error) {
	var src []byte
	switch srcdata.(type) {
	case []byte:
		src = srcdata.([]byte)
	case string:
		src = []byte(srcdata.(string))
	}

	var resultBytes []byte
	var jsonFormat map[string]interface{}
	if err := json.Unmarshal(src, &jsonFormat); err != nil {
		log.Fatal(err)
	}
	if resultBytes, err = json.MarshalIndent(jsonFormat, "", "	"); err != nil {
		return ret, err
	}

	// fmt.Printf("%s\n", resultBytes)
	ret = string(resultBytes)
	return
}

// GetCurFilename returns the filename in short at the stack whitch is skipDepth above cur stack.
// Arguments skipDepth is usually 1 in normal use which is the caller of the GetCurFilename.
func GetCurFilename(skipDepth int) (short, path string, err error) {
	_, filename, _, ok := runtime.Caller(skipDepth)
	if !ok {
		err := "runtime.Caller error "
		panic(err)
	}
	pathArr := strings.Split(filename, string(os.PathSeparator))
	short = pathArr[len(pathArr)-1]
	path = strings.Join(pathArr[:len(pathArr)-1], string(os.PathSeparator))
	return
}

// GetLogInfo gets the runtime info for the calling function in form of "|short(filename)|fn|lineno".
func GetLogInfo(skipDepth int) (info string) {
	pc, _, lineno, _ := runtime.Caller(skipDepth)
	short, _, _ := GetCurFilename(skipDepth + 1)
	f := runtime.FuncForPC(pc)
	fn := f.Name()
	fnArray := strings.Split(fn, ".")
	var fnShort string
	if fnArray != nil {
		fnShort = fnArray[len(fnArray)-1]
	}
	info = fmt.Sprintf("|%v|%v|%v", short, fnShort, lineno)
	return
}

func Roadmap(mark ...interface{}) {
	log.Printf("\n[ROADMAP] =============== %v ================== [ROADMAP]\n", mark)

}
