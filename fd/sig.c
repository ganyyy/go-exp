#include <sys/signalfd.h>
#include <signal.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdio.h>


#define handle_err(msg) \
    do {perror(msg); exit(EXIT_FAILURE);} while(0);

int main(int argc, char const *argv[])
{
    sigset_t mask;
    sigemptyset(&mask);
    sigaddset(&mask, SIGINT);

    if (sigprocmask(SIG_BLOCK, &mask, NULL) == -1)
        handle_error("sigprocmask");

    int sfd = signalfd(-1, &mask, 0);
    if(sfd < 0)
    {
        perror("create signal fd error!");
        exit -1;
    }
    struct signalfd_siginfo sfd_info;
    for (;;) {
        size_t s = read(sfd, &sfd_info, sizeof(struct signalfd_siginfo));
        if (s != sizeof(struct signalfd_siginfo))
            handle_error("read");
 
        if (sfd_info.ssi_signo == SIGINT) {
            printf("Got SIGINT\n");
        } else if (sfd_info.ssi_signo == SIGQUIT) {
            printf("Got SIGQUIT\n");
            exit(EXIT_SUCCESS);
        } else {
            printf("Read unexpected signal\n");
        }
    }

    return 0;
}
