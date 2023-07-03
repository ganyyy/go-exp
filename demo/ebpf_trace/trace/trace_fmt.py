from bcc import BPF
from time import strftime

ELF_FILE_PATH = b"/home/gan/code/go/my_code/go-exp/demo/ebpf_trace/ebpf_trace"

b = BPF(src_file=b"trace_fmt.c")

b.attach_uprobe(name=ELF_FILE_PATH,
                sym=b"os.open", fn_name=b"go_open_file")

b.attach_uprobe(name=ELF_FILE_PATH,
                sym=b"os.(*File).Write", fn_name=b"go_write_data")


def print_event(cpu, data, size):
    event = b["events"].event(data)
    print("%-9s %-6d arg1:%-20s, arg2:%-8d, arg3:%-4d" %
          (strftime("%H:%M:%S"), event.uid, event.arg1[:event.ln], event.arg2, event.arg3))


print("%-9s %-6s %-26s %-16s %-16s" % ("TIME", "PID", "ARG1", "ARG2", "ARG3"))

b["events"].open_perf_buffer(print_event)

while True:
    try:
        b.perf_buffer_poll()
    except KeyboardInterrupt:
        exit()
