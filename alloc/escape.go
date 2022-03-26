package main

/*

相关配置在
cmd/compile/internal/ir/cfg.go


单个对象的分配大小
MaxStackVarSize = 10 * 1024 * 1024 = 10M

使用new/&/make/[]byte("") 创建的对象大小
MaxImplicitStackVarSize = 64 * 1024 = 64KB



*/

func EscapeToHeap1() {
	var arr1 [10 * 1024 * 1024]byte // 上限是这个值, 栈上超过了这个大小会直接分配到堆上
	var arr2 [10*1024*1024 + 1]byte // arr2逃逸到了堆上

	_ = arr1
	_ = arr2
}

func EscapeToHeap2() {
	var p1 = new([64 * 1024]byte)
	var p2 = new([64*1024 + 1]byte) // 逃逸到了堆上

	_ = p1
	_ = p2
}
