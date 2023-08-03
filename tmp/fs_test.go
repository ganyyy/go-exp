package main

import (
	"io/fs"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFSFunction(t *testing.T) {
	// 是否是有效的路径
	t.Log(fs.ValidPath("path"))

	// 打开一个文件
	var of, err = os.Open("./fs_test.go")
	require.NoError(t, err)
	var f fs.File = of
	defer f.Close()

	var entries []fs.DirEntry
	var fStat fs.FileInfo

	// 打开一个目录
	var dir = os.DirFS("./")
	//	// 读取文件
	_, err = fs.ReadFile(dir, "fs_test.go")
	require.NoError(t, err)
	// 	// 读取目录
	entries, err = fs.ReadDir(dir, ".")
	require.NoError(t, err)
	if len(entries) > 0 {
		t.Logf("entrie:%v", fs.FormatDirEntry(entries[0]))
	}
	// 	// 读取文件信息
	fStat, err = fs.Stat(dir, "fs_test.go")
	require.NoError(t, err)
	t.Logf("file stat:%v", fs.FormatFileInfo(fStat))

	// 匹配文件名
	var matches []string
	matches, err = fs.Glob(dir, "*.go")
	require.NoError(t, err)
	t.Logf("matches:%v", strings.Join(matches, ","))

	// 读取子目录
	var subDir fs.FS
	subDir, err = fs.Sub(dir, "benchmark")
	require.NoError(t, err)
	_ = subDir

	// // 迭代目录
	err = fs.WalkDir(subDir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		t.Logf("path:%v, info:%v", path, fs.FormatDirEntry(d))
		return nil
	})
	require.NoError(t, err)

}
