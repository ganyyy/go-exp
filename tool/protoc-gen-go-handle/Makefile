# 项目根目录
PROTO_DIR := proto
PROTOGEN_DIR := protogen
MODULE_DIR := module
COMMON_PACKAGE := protoc-gen-go-handle/common

.PHONY: proto clean install

proto:
	@TMP=$$(mktemp); \
	find $(PROTO_DIR) -name '*.proto' > $$TMP; \
	mkdir -p $(PROTOGEN_DIR); \
	protoc \
		--proto_path=$(PROTO_DIR) \
		--go_out=$(PROTOGEN_DIR) \
		--go_opt=paths=source_relative \
		--go-handle_out=$(MODULE_DIR) \
		--go-handle_opt=module_base_path=$(MODULE_DIR) \
		--go-handle_opt=common_pkg=$(COMMON_PACKAGE) \
		@$$TMP; \
	rm -f $$TMP; \
	echo "protoc-gen-go-handle: $(PROTO_DIR) -> $(PROTOGEN_DIR)"

clean:
	@rm -rf $(PROTOGEN_DIR)
	@echo "clean: $(PROTOGEN_DIR)"

install:
	go install -trimpath -ldflags="-s -w"