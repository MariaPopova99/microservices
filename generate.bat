@echo off
REM Проверим, существует ли директория pkg/note_v1, если нет — создадим
if not exist pkg\note_v1 (
    mkdir pkg\note_v1
)

REM Генерация через protoc
protoc --proto_path=api\note_v1 --go_out=pkg\note_v1 --go_opt=paths=source_relative --plugin=protoc-gen-go=C:/Users/hallo/go/bin/protoc-gen-go.exe --go-grpc_out=pkg\note_v1 --go-grpc_opt=paths=source_relative --plugin=protoc-gen-go-grpc=C:/Users/hallo/go/bin/protoc-gen-go-grpc.exe api\note_v1\note.proto

REM Проверяем завершение процесса
if %ERRORLEVEL% neq 0 (
    echo Ошибка при генерации
    exit /b %ERRORLEVEL%
)

echo Генерация завершена успешно.
