package update

import (
	"log"
	"os"
	"os/signal"
	"patch"
	"plugin"
	"reflect"
	"strings"
	"syscall"
	"unsafe"
)

const versionPath = "./plugin/version"
const exchange = "Exchange"

type backSymbol struct {
	back []byte
	name string
}

var (
	symbolBackup = map[unsafe.Pointer]backSymbol{}
	loaded       bool
)

func loadPlugin() {
	// 缺少了类型检查..

	var handler *plugin.Plugin
	var mainHandler unsafe.Pointer
	var symbol plugin.Symbol
	var err error

	var version []byte
	var pluginPath string
	version, err = os.ReadFile(versionPath)
	if err != nil {
		log.Printf("read version file %v error %v", versionPath, err)
		return
	}
	pluginPath = strings.TrimSpace(string(version))
	if _, err = os.Stat(pluginPath); err != nil {
		log.Printf("cannot found plugin %v", pluginPath)
		return
	}
	// 加载动态库
	handler, err = plugin.Open(pluginPath)
	if err != nil {
		log.Printf("load %v error:%v", pluginPath, err)
		return
	}

	// 加载主函数表
	mainHandler, err = patch.PluginOpen("")
	if err != nil {
		log.Printf("load main symbol error:%v", err)
		return
	}
	defer func() {
		err = patch.PluginClose(mainHandler)
		if err != nil {
			log.Printf("close mainHandler error:%v", err)
		} else {
			log.Printf("close mainHandler success!")
		}
	}()

	// 加载替换函数列表
	symbol, err = handler.Lookup(exchange)
	if err != nil {
		log.Printf("lookup plugin %v error %v", exchange, err)
		return
	}

	symbolMapPtr, ok := symbol.(*map[string]string)
	if !ok {
		log.Printf("plugin symbol %v type %v error!", exchange, reflect.TypeOf(symbol))
		return
	}
	if symbolMapPtr == nil {
		log.Printf("plugin symbol %v is nil!", exchange)
		return
	}

	symbolMap := *symbolMapPtr

	type exchangeFunc struct {
		oldFunc unsafe.Pointer
		oldName string
		newFunc interface{}
	}

	var patchFunc = make([]exchangeFunc, 0, len(symbolMap))

	for newFunc, oldFunc := range symbolMap {
		symbol, err := handler.Lookup(newFunc)
		if err != nil {
			log.Printf("load plugin symbol %v error %v", newFunc, err)
			return
		}
		mainSymbol, err := patch.LookupSymbol(mainHandler, oldFunc)
		if err != nil {
			log.Printf("load main symbol %v error %v", oldFunc, err)
			return
		}
		patchFunc = append(patchFunc, exchangeFunc{
			oldFunc: mainSymbol,
			oldName: oldFunc,
			newFunc: symbol,
		})
	}

	if loaded {
		restorePlugin()
	}

	for _, p := range patchFunc {
		var back backSymbol
		back.name = p.oldName
		back.back = patch.Backup(p.oldFunc)
		symbolBackup[p.oldFunc] = back
		patch.Patch(p.oldFunc, patch.FuncAddr(p.newFunc))
	}

	loaded = true
	log.Printf("patch %+v success!", symbolMapPtr)
}

func restorePlugin() {
	var backName []string
	for addr, back := range symbolBackup {
		patch.Restore(addr, back.back)
		backName = append(backName, back.name)
	}
	for k := range symbolBackup {
		delete(symbolBackup, k)
	}
	log.Printf("restore patch %+v success!", backName)
	loaded = false
}

func RunUpdateMonitor() {
	var sigChan = make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGUSR1, syscall.SIGUSR2)

	for sig := range sigChan {
		switch sig {
		case syscall.SIGUSR1:
			loadPlugin()
		case syscall.SIGUSR2:
			restorePlugin()
		}
	}
}
