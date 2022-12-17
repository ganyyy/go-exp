#include "textflag.h"
#include "funcdata.h"
#include "go_asm.h"

TEXT ·myadd(SB), $0-24 
    MOVQ a+0(FP), AX
    MOVQ b+8(FP), BX
    ADDQ BX, AX
    MOVQ AX, ret+16(FP)
    RET

TEXT ·mtest(SB), $0-0
    LEAQ (AX), AX
    LEAQ 100(CX), CX
    MOVQ 200(DX), DX
    RET
