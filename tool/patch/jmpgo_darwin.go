package patch

var (
	jmpGo = [...]byte{
		0x48, 0xBA,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // mov ptr rdx
		0xFF, 0x22, // jmp QWORD PTR [rdx]
	}
)
