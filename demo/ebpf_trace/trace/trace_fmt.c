#include <uapi/linux/ptrace.h>

struct data_t
{
    u32 uid;
    char arg1[20];
    int ln;
    int arg2;
    int arg3;
};

BPF_PERF_OUTPUT(events);

int go_open_file(struct pt_regs *ctx)
{
    // arg1: {RAX, RBX}, string
    // arg2: RCX, int
    // arg3: RDX, int
    struct data_t data = {};
    data.uid = bpf_get_current_uid_gid();

    // 读取用户空间的字符串到data
    u32 size = (u32)(ctx->bx);
    bpf_probe_read_user(&data.arg1, sizeof(data.arg1), (void *)ctx->ax);
    data.ln = size;
    data.arg2 = ctx->cx;
    data.arg3 = ctx->di;
    events.perf_submit(ctx, &data, sizeof(data));
    return 0;
}

int go_write_data(struct pt_regs *ctx)
{
    // receiver: RAX
    // arg1: {RBX, RCX, RDI}, []byte
    struct data_t data = {};
    data.uid = bpf_get_current_uid_gid();

    // 读取用户空间的字符串到data
    u32 size = (u32)(ctx->cx);
    bpf_probe_read_user(&data.arg1, sizeof(data.arg1), (void *)ctx->bx);
    data.ln = size;
    events.perf_submit(ctx, &data, sizeof(data));
    return 0;
}