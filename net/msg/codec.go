package msg

/*

Header: 4 bytes
Body: 	bytes

*/

type Message struct {
	Message string // 这里假定使用json进行编码/解码
}
