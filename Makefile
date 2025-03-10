# 定义变量
SWAG_CMD = swag init -g cmd/main.go

# 默认任务
all: swagger

# 生成 Swagger 文档
swagger:
	@echo "Generating Swagger docs..."
	$(SWAG_CMD)
	@echo "Swagger docs generated successfully!"

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all       - Default target, runs 'swagger'"
	@echo "  swagger   - Generate Swagger docs"
	@echo "  help      - Show this help message"

mock:
	@echo "start generate mock code"
	mockgen -source internal/repository/repo/book.go -destination internal/pkg/mocks/mock_repo.go -package mocks
	@echo  "generate mock code end"