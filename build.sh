#!/bin/bash

# Скрипт сборки Jotnal IDE для различных платформ

set -e

VERSION=${VERSION:-"1.0.0"}
OUTPUT_DIR="./build"

# Цвета для вывода
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Сборка Jotnal IDE v${VERSION}${NC}"

# Создаём директорию для сборки
mkdir -p "$OUTPUT_DIR"

# Функция сборки
build() {
    local GOOS=$1
    local GOARCH=$2
    local OUTPUT=$3
    local CC=$4

    echo -e "${GREEN}Сборка для ${GOOS}/${GOARCH}...${NC}"

    if [ -n "$CC" ]; then
        CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH CC=$CC go build -ldflags="-s -w" -o "$OUTPUT_DIR/$OUTPUT" ./cmd/ide
    else
        CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="-s -w" -o "$OUTPUT_DIR/$OUTPUT" ./cmd/ide
    fi

    echo -e "${GREEN}✓ Создан: $OUTPUT_DIR/$OUTPUT${NC}"
    ls -lh "$OUTPUT_DIR/$OUTPUT"
}

# Выбор платформы
case "${1:-all}" in
    windows|win)
        build windows amd64 jotnal.exe x86_64-w64-mingw32-gcc
        ;;
    linux)
        build linux amd64 jotnal ""
        ;;
    darwin|mac|macos)
        build darwin amd64 jotnal-darwin ""
        ;;
    all)
        echo -e "${BLUE}Сборка для всех платформ...${NC}"
        build linux amd64 jotnal ""
        build windows amd64 jotnal.exe x86_64-w64-mingw32-gcc
        echo -e "${GREEN}Сборка для macOS требует запуска на macOS системе${NC}"
        ;;
    clean)
        echo -e "${BLUE}Очистка директории сборки...${NC}"
        rm -rf "$OUTPUT_DIR"
        rm -f jotnal jotnal.exe jotnal-darwin
        echo -e "${GREEN}✓ Очистка завершена${NC}"
        ;;
    *)
        echo "Использование: $0 {windows|linux|darwin|all|clean}"
        echo ""
        echo "Примеры:"
        echo "  $0 windows  - собрать только для Windows"
        echo "  $0 linux    - собрать только для Linux"
        echo "  $0 all      - собрать для Linux и Windows"
        echo "  $0 clean    - очистить сборки"
        exit 1
        ;;
esac

echo -e "${GREEN}✓ Сборка завершена успешно!${NC}"
