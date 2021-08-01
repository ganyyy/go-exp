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
    if (sock < 0)
    {
        perror("create sock error");
    }
    return sock;
}