# Скрипт сборки Jotnal IDE для различных платформ (Windows версия)
# Использование: .\build.ps1 [windows|linux|darwin|all|clean]

param(
    [Parameter(Position=0)]
    [ValidateSet('windows', 'win', 'linux', 'darwin', 'mac', 'macos', 'all', 'clean')]
    [string]$Target = 'all'
)

$VERSION = if ($env:VERSION) { $env:VERSION } else { "1.0.0" }
$OUTPUT_DIR = ".\build"

# Цвета для вывода
function Write-Success { param($Message) Write-Host $Message -ForegroundColor Green }
function Write-Info { param($Message) Write-Host $Message -ForegroundColor Blue }
function Write-Step { param($Message) Write-Host $Message -ForegroundColor Cyan }

Write-Info "Сборка Jotnal IDE v$VERSION"

# Создаём директорию для сборки
if (-not (Test-Path $OUTPUT_DIR)) {
    New-Item -ItemType Directory -Path $OUTPUT_DIR | Out-Null
}

# Функция сборки
function Build-Binary {
    param(
        [string]$GOOS,
        [string]$GOARCH,
        [string]$OUTPUT,
        [string]$CC = ""
    )

    Write-Step "Сборка для $GOOS/$GOARCH..."

    $env:CGO_ENABLED = "1"
    $env:GOOS = $GOOS
    $env:GOARCH = $GOARCH

    if ($CC) {
        $env:CC = $CC
    }

    $OutputPath = Join-Path $OUTPUT_DIR $OUTPUT

    try {
        go build -ldflags="-s -w" -o $OutputPath .\cmd\ide

        if ($LASTEXITCODE -eq 0) {
            Write-Success "✓ Создан: $OutputPath"
            $fileInfo = Get-Item $OutputPath
            $sizeMB = [math]::Round($fileInfo.Length / 1MB, 2)
            Write-Host "  Размер: $sizeMB MB" -ForegroundColor Gray
        } else {
            throw "Ошибка при сборке"
        }
    } catch {
        Write-Host "✗ Ошибка сборки для $GOOS/$GOARCH" -ForegroundColor Red
        Write-Host $_.Exception.Message -ForegroundColor Red
        exit 1
    }
}

# Выбор платформы
switch ($Target) {
    { $_ -in 'windows', 'win' } {
        Write-Info "Сборка для Windows..."
        Build-Binary -GOOS "windows" -GOARCH "amd64" -OUTPUT "jotnal.exe"
    }

    'linux' {
        Write-Info "Сборка для Linux..."
        Write-Host "Внимание: для сборки Linux на Windows требуется WSL или Docker" -ForegroundColor Yellow
        Build-Binary -GOOS "linux" -GOARCH "amd64" -OUTPUT "jotnal"
    }

    { $_ -in 'darwin', 'mac', 'macos' } {
        Write-Info "Сборка для macOS..."
        Write-Host "Внимание: для сборки macOS на Windows требуется Docker" -ForegroundColor Yellow
        Build-Binary -GOOS "darwin" -GOARCH "amd64" -OUTPUT "jotnal-darwin"
    }

    'all' {
        Write-Info "Сборка для всех платформ..."
        Build-Binary -GOOS "windows" -GOARCH "amd64" -OUTPUT "jotnal.exe"
        Write-Host ""
        Write-Host "Для сборки под Linux и macOS используйте WSL или Docker" -ForegroundColor Yellow
    }

    'clean' {
        Write-Info "Очистка директории сборки..."
        if (Test-Path $OUTPUT_DIR) {
            Remove-Item -Recurse -Force $OUTPUT_DIR
        }
        Remove-Item -ErrorAction SilentlyContinue "jotnal.exe", "jotnal", "jotnal-darwin"
        Write-Success "✓ Очистка завершена"
        exit 0
    }
}

Write-Success "`n✓ Сборка завершена успешно!"
Write-Info "`nСобранные файлы находятся в директории: $OUTPUT_DIR"
