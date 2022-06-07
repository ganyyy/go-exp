package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

type ParseParam struct {
	InputFile  string
	OutputPath string
	GoPackage  string
	UseNumber  bool
}

const (
	jsonExt = ".json"
)

func (param *ParseParam) InitOutput() error {
	if stat, err := os.Stat(param.OutputPath); err != nil {
		if fe, ok := err.(*fs.PathError); !ok || fe.Err != syscall.ENOENT {
			return err
		}
		err := os.Mkdir(param.OutputPath, fs.ModePerm)
		if err != nil {
			return err
		}
	} else {
		if !stat.IsDir() {
			return errors.New("must input valid dir path")
		}
	}
	return nil
}

func (param *ParseParam) Parse() error {
	stat, err := os.Stat(param.InputFile)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		// 输入的是一整个目录
		filepath.WalkDir(param.InputFile, func(path string, d fs.DirEntry, err error) error {
			// 目前只支持到子一级的目录, 如果有需要, 后期再添加
			if err != nil || d.IsDir() {
				return nil
			}
			//TODO 启用协程池处理
			return param.parseFile(path)
		})
	} else {
		// 单个文件
		return param.parseFile(param.InputFile)
	}

	return nil
}

func (param *ParseParam) parseFile(path string) error {
	var ext = filepath.Ext(path)
	if ext != jsonExt {
		return nil
	}

	var base = filepath.Base(path)
	if len(base) <= len(jsonExt) {
		return nil
	}
	var bs, err = os.ReadFile(path)
	if err != nil {
		return err
	}
	return parseInputData(base[:len(base)-len(ext)], bs, param)
}

func parseInputData(base string, data []byte, param *ParseParam) error {
	var decoder = json.NewDecoder(bytes.NewReader(data))
	if param.UseNumber {
		decoder.UseNumber()
	}
	var v naiveValue
	if err := decoder.Decode(&v); err != nil {
		return err
	}
	var obj, err = parseValue(v)
	if err != nil && err != ErrEmptySlice {
		return err
	}
	if err == ErrEmptySlice {
		// 顶层的空切片, 无视即可
		return nil
	}
	obj.KeyName = base
	obj.TypeName = title(obj.KeyName)

	var allType = ParseAllType(obj)

	var output = param.OutputPath + "/" + base + ".go"
	var tp TemplateParse
	tp.Root = obj
	tp.AllType = allType
	tp.PkgName = param.GoPackage
	return tp.Parse(output)
}
