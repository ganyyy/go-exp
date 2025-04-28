# 项目根目录
PROTO_DIR := proto
PROTOGEN_DIR := protogen
MODULE_DIR := module
COMMON_PACKAGE := protoc-gen-go-handle/common

.PHONY: all-proto clean install proto


# 内部函数：生成proto文件
# 参数：$(1) 是proto文件列表或@文件名，例如 "file1.proto file2.proto" 或 "@filelist.txt"
define generate_proto
    mkdir -p $(PROTOGEN_DIR); \
    protoc \
        --proto_path=$(PROTO_DIR) \
        --go_out=$(PROTOGEN_DIR) \
        --go_opt=paths=source_relative \
        --go-handle_out=$(MODULE_DIR) \
        --go-handle_opt=module_base_path=$(MODULE_DIR) \
        --go-handle_opt=common_pkg=$(COMMON_PACKAGE) \
        $(1); \
    echo "Generated proto code for: $(1)"
endef

all-proto:
	@TMP=$$(mktemp); \
	find $(PROTO_DIR) -name '*.proto' > $$TMP; \
	$(call generate_proto,@$$TMP); \
	rm -f $$TMP; \
	echo "protoc-gen-go-handle: $(PROTO_DIR) -> $(PROTOGEN_DIR)"


proto:
	@if [ -z "$(PROTO)" ]; then \
		echo "Error: PROTO is required. Usage: make proto PROTO=path/to/your.proto"; \
		exit 1; \
	fi; \
	$(call generate_proto,$(addprefix $(PROTO_DIR)/,$(PROTO)))

clean:
	@rm -rf $(PROTOGEN_DIR)
	@echo "clean: $(PROTOGEN_DIR)"

install:
	go install -trimpath -ldflags="-s -w"