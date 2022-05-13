#include <stdio.h>
#include <stdlib.h>
#include "_cgo_export.h"


int helloFromC() {
    printf("Hi from C\n");
    //call Go function
    HelloFromGo();
    return 0;
}


int addFromC(int a, int b)
{
    printf("hello world!\n");
    return Add(a, b);
}


void print_str(const char *s)
{
    printf("%s\n", s?s:"nil");
    // C语言调用 Go语言函数
    GoString name;
    name.p = "123";
    name.n = 3;
    char *greeting = HelloByGo(name);
    printf("call go: %s\n", greeting);
    free(greeting);
}