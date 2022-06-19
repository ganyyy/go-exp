package main

import (
	"fmt"
	"os"

	"ganyyy.com/go-exp/demo/golang-lua/scripts"
	"github.com/arnodel/golua/lib/base"
	"github.com/arnodel/golua/runtime"
)

type User = scripts.User

// GoFunctionFunc func(*Thread, *GoCont) (Cont, error)

func UserToString(t *runtime.Thread, cont *runtime.GoCont) (runtime.Cont, error) {
	var data *runtime.UserData
	err := cont.Check1Arg()
	if err != nil {
		return nil, err
	}

	data, err = cont.UserDataArg(0)
	if err != nil {
		return nil, err
	}
	next := cont.Next()
	if user, ok := data.Value().(*User); !ok {
		return nil, err
	} else {
		t.Push1(next, runtime.StringValue(user.String()))
		return next, nil
	}
}

func UserSetName(t *runtime.Thread, cont *runtime.GoCont) (runtime.Cont, error) {
	var data *runtime.UserData
	var newName string
	err := cont.CheckNArgs(2)
	if err != nil {
		return nil, err
	}
	data, err = cont.UserDataArg(0)
	if err != nil {
		return nil, err
	}
	newName, err = cont.StringArg(1)
	if err != nil {
		return nil, err
	}
	if user, ok := data.Value().(*User); ok {
		user.SetName(newName)
		return cont.Next(), nil
	}
	return cont.Next(), nil
}

// func ParseLuaToGo(gt reflect.Type, lt runtime.Value) reflect.Value {
// 	var val = reflect.Zero(gt)
// 	switch gt.Kind() {
// 	case reflect.Bool:
// 		val.SetBool(lt.AsBool())
// 	case reflect.String:
// 		val.SetString(lt.AsString())
// 	case reflect.Float32, reflect.Float64:
// 	case reflect.Int,
// 		reflect.Int8,
// 		reflect.Int16,
// 		reflect.Int32,
// 		reflect.Int64:
// 		val.SetInt(lt.AsInt())
// 	case reflect.Uint,
// 		reflect.Uint8,
// 		reflect.Uint16,
// 		reflect.Uint32,
// 		reflect.Uint64:
// 		val.SetUint(uint64(lt.AsInt()))
// 	case reflect.Map:
// 		// 不支持了, 艹, 本来就转不了. 日了狗了
// 		lt.AsTable().Next() // 费尽心思写这个干啥. 艹
// 	case reflect.Array, reflect.Slice:
// 	case reflect.Struct:

// 	case reflect.Ptr:
// 		// 递归搞
// 	}
// 	return val
// }

// //TODO 用反射包一层?
// // 目前只支持基础类型,
// func RegisterGoFunc(name string, f interface{}, isMethod bool) *runtime.GoFunction {
// 	var ft = reflect.TypeOf(f)
// 	if ft.Kind() != reflect.Func {
// 		return nil
// 	}
// 	var vt = reflect.ValueOf(f)
// 	var argsLen = ft.NumIn()

// 	if isMethod && (argsLen < 1 || ft.In(0).Kind() != reflect.Ptr) {
// 		return nil
// 	}

// 	// 如果为了性能考虑, 这样写是肯定不行的
// 	// 运行时使用反射调用方法性能贼差

// 	var run = func(t *runtime.Thread, cont *runtime.GoCont) (runtime.Cont, error) {
// 		err := cont.CheckNArgs(argsLen)
// 		if err != nil {
// 			return nil, err
// 		}
// 		var callArgs = make([]reflect.Value, 0, argsLen)
// 		for i, arg := range cont.Args() {

// 			callArgs = append(callArgs, reflect.ValueOf(arg.Interface()))
// 		}
// 		results := vt.Call(callArgs)
// 		var retValue = make([]runtime.Value, 0, len(results))
// 		for _, ret := range results {
// 			retValue = append(retValue, runtime.AsValue(ret.Interface()))
// 		}
// 		t.Push(cont, retValue...)
// 		return cont.Next(), nil
// 	}

// 	return runtime.NewGoFunction(run, name, ft.NumIn(), false)
// }

func main() {

	var rt = runtime.New(os.Stdout)
	var userMetaTable = runtime.NewTable()

	base.Load(rt)

	var methodTable = runtime.NewTable()
	methodTable.Set(runtime.StringValue("SetName"), runtime.FunctionValue(runtime.NewGoFunction(UserSetName, "SetName", 2, false)))
	methodTable.Set(runtime.StringValue("String"), runtime.FunctionValue(runtime.NewGoFunction(UserToString, "String", 1, false)))

	// var fieldTable = runtime.NewTable()

	userMetaTable.Set(runtime.StringValue("__index"), runtime.TableValue(methodTable))

	var user = &User{
		Name: "gan",
		Age:  100,
	}
	rt.GlobalEnv().Set(runtime.StringValue("user"),
		runtime.UserDataValue(runtime.NewUserData(user, userMetaTable)))

	var source, _ = os.ReadFile(scripts.CounterPath)

	// 编译lua
	chunk, _ := rt.CompileAndLoadLuaChunk("counter", source, runtime.TableValue(rt.GlobalEnv()))
	term := runtime.NewTerminationWith(rt.MainThread().CurrentCont(), 0, false)
	// 执行lua脚本
	runtime.Call(rt.MainThread(), runtime.FunctionValue(chunk), nil, term)

	// 获取全局的对象
	v := rt.MainThread().GlobalEnv().Get(runtime.StringValue("NXT"))
	// 这是个函数对象, 可以执行
	ret, _ := runtime.Call1(rt.MainThread(), v)
	fmt.Println(ret.AsInt())

	fmt.Println(user)
}
