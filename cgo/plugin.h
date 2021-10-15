#include <stdio.h>
#include <netinet/in.h>
#include <stdlib.h>

int add(int a, int b)
{
    return a + b;
}

int tcpSocket()
{
    int sock = socket(AF_INET, SOCK_STREAM, 0);
//    struct GoString gs;
//    gs.p = "12345";
//    gs.n = 5;
//    extern GoInt GoF(GoInt, GoInt, GoString);
//    long long val = GoF(100, 200, gs);
    if (sock < 0)
    {
        perror("create sock error");
    }
    return sock;
}