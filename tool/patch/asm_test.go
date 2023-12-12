package patch

import (
	"fmt"
	"strings"
	"testing"

	"github.com/chenzhuoyu/iasm/x86_64"
	"golang.org/x/arch/x86/x86asm"
)

func TestAsm(t *testing.T) {
	var mov = x86asm.Inst{Op: x86asm.MOVQ, Args: x86asm.Args{
		x86asm.RAX, x86asm.Imm(0x12345678),
	}}
	var jmp = x86asm.Inst{Op: x86asm.JMP, Args: x86asm.Args{
		x86asm.RAX,
	}}

	t.Log(x86asm.GoSyntax(mov, 0, nil), x86asm.GoSyntax(jmp, 0, nil))
	t.Log(x86asm.IntelSyntax(mov, 0, nil), x86asm.IntelSyntax(jmp, 0, nil))

	t.Logf("======分割线=======")

	d := func(bs []byte) {
		ins, err := x86asm.Decode(bs, 64)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf(ins.String())
	}

	d([]byte{0x48, 0xba, 0xde, 0xbc, 0x9a, 0x78, 0x56, 0x34, 0x12, 0x00})
	d([]byte{0xff, 0xe2})
}

func TestX8664(t *testing.T) {

	const REGISTER = x86_64.R12

	showASM := func(register x86_64.Register) {
		pp := func(ins []byte) {
			if len(ins) == 0 {
				t.Log("[]byte{}")
				return
			}
			var _buf [64]byte
			var bs strings.Builder

			const fmtHexInner = "0x%02X, "
			var fmtHexEnd = fmtHexInner[:len(fmtHexInner)-2]

			w := func(b byte, format string) { bs.Write(fmt.Appendf(_buf[:0], format, b)) }

			ln := len(ins)
			bs.WriteString("[]byte{")
			for _, b := range ins[:ln-1] {
				w(b, fmtHexInner)
			}
			w(ins[ln-1], fmtHexEnd)
			bs.WriteString("}")

			t.Log(bs.String())
		}
		// 这个第三方库, 采用类似于orm的形式, 将汇编逻辑转变成二进制字节码
		// 注意: 立即数得是64位的数(数值的字面量需要大于math.MaxInt32), 否则指令会被优化.
		// 这里建议只修改目标寄存器, MOVQ和JMPQ的目标寄存器都需要修改
		{
			var p x86_64.Program
			p.MOVQ(int64(0x123456789abcde), register)
			pp(p.AssembleAndFree(0))
		}

		{
			// 直接寻址
			var p x86_64.Program
			p.JMPQ(register)
			pp(p.AssembleAndFree(0))
		}

		{
			var p x86_64.Program
			// 间接寻址
			p.JMPQ(&x86_64.MemoryOperand{
				Size: 8,
				Addr: x86_64.Addressable{
					Type: x86_64.Memory,
					Memory: x86_64.MemoryAddress{
						Base: register,
					},
				},
			})
			pp(p.AssembleAndFree(0))
		}
	}

	showASM(REGISTER)
}
