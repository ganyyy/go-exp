CC=go build
CFLAGS=-x
PLUGIN_FLAG="-buildmode=plugin"


plugin: plugin.go
	$(CC) $(PLUGIN_FLAG) -o $@.so $(CFLAGS) $^

main: main.go
	$(CC) -o $@ $(CFLAGS) $^

