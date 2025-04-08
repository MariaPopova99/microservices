# Установка зависимостей через .bat файл
install-deps:
	@call install-deps.bat

# Получение зависимостей через Go
get-deps:
	@powershell -Command "go get -u google.golang.org/protobuf/cmd/protoc-gen-go"
	@powershell -Command "go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc"

# Генерация файлов
generate:
	@call generate.bat

generate-note-api:
	@echo "Generating the note API..."
