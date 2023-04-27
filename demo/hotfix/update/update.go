package update

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"plugin"
	"reflect"
	"strings"
	"syscall"
	"unsafe"

	"ganyyy.com/go-exp/patch"
)

const (
	PluginBase  = "./plugin/"
	VersionPath = PluginBase + "version"
)

const (
	VersionSymbol  = "Version"
	ExchangeSymbol = "Exchange"
)

type Plugin struct {
	pluginpath string
	err        string        // set if plugin failed to load
	loaded     chan struct{} // closed when loaded
	syms       map[string]any
}

// ShowPlugin Plugin data
func ShowPlugin(pp *plugin.Plugin) {
	var p = (*Plugin)(unsafe.Pointer(pp))
	log.Printf("plugin %v, error %v", p.pluginpath, p.err)
	for k, v := range p.syms {
		log.Printf("plugin symbol %v, value %v", k, v)
	}
}

type backSymbol struct {
	back []byte
	name string
}

var (
	symbolBackup = map[unsafe.Pointer]backSymbol{}
	loaded       bool
)

func loadSymbol[T any](name string, pluginHandler *plugin.Plugin) (T, error) {
	var symbol, err = pluginHandler.Lookup(name)
	var empty T
	if err != nil {
		return empty, err
	}
	var wantType, realType = reflect.TypeOf(empty), reflect.TypeOf(symbol)
	var realValue = reflect.ValueOf(symbol)
	if !realType.AssignableTo(wantType) {
		if realType.Kind() == reflect.Pointer {
			realType = realType.Elem()
			realValue = realValue.Elem()
		}
		if !realType.AssignableTo(wantType) {
			return empty, fmt.Errorf("cannot exchange symbol %v type %v to %v", name, realType, wantType)
		}
	}
	return realValue.Interface().(T), nil
}

func loadPlugin() {
	// 缺少了类型检查..

	var handler *plugin.Plugin
	var mainHandler unsafe.Pointer
	var pluginVersion string
	var err error
	var version []byte
	var pluginFileVersion string
	var pluginName string

	version, err = os.ReadFile(VersionPath)
	if err != nil {
		log.Printf("read version file %v error %v", VersionPath, err)
		return
	}
	pluginFileVersion = strings.TrimSpace(string(version))
	pluginName = PluginBase + pluginFileVersion
	if _, err = os.Stat(pluginName); err != nil {
		log.Printf("cannot found plugin %v", pluginName)
		return
	}
	// 加载动态库
	handler, err = plugin.Open(pluginName)
	if err != nil {
		log.Printf("load %v error:%v", pluginName, err)
		return
	}

	// 输出插件相关的信息
	// ShowPlugin(handler)

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

	var symbolMap map[string]string

	symbolMap, err = loadSymbol[map[string]string](ExchangeSymbol, handler)
	if err != nil {
		log.Printf("loadSymbol Exchange error:%v", err)
		return
	}

	// 对比版本信息
	pluginVersion, err = loadSymbol[string](VersionSymbol, handler)
	if err != nil {
		log.Printf("loadSymbol Version error:%v", err)
		return
	}

	if pluginVersion != pluginFileVersion {
		log.Printf("plugin file version:%v, plugin version %v not match!", pluginFileVersion, pluginVersion)
		return
	}

	//do patch

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
	log.Printf("patch %+v success!", symbolMap)
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
