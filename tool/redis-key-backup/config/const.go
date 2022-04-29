package config

const (
	Host       = "host"      // 指定redis地址 Addr:Port
	HostFull   = "host, hh"  // 指定redis地址 Addr:Port
	DB         = "db"        // 指定数据库
	DBFull     = "db, n"     // 指定数据库
	Key        = "key"       // 指定key
	KeyFull    = "key, k"    // 指定key
	Auth       = "auth"      // 指定数据库密码
	AuthFull   = "auth, a"   // 指定数据库密码
	File       = "file"      // 指定文件输入
	FileFull   = "file, f"   // 指定文件输入
	Input      = "input"     // 指定控制台输入
	InputFull  = "input, i"  // 指定控制台输入
	Output     = "output"    //  指定控制台输出
	OutputFull = "output, o" //  指定控制台输出
)

const (
	TypeString = "string"
	TypeHMap   = "hash"
	TypeZSet   = "zset"
	TypeList   = "list"
	// ... etc
)
