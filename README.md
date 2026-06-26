# Сервис Shortlink

- [Результаты pprof diff](profiles/diff.md)

### Сборка
```bash
go build -ldflags="-X main.buildVersion=1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(git rev-parse HEAD)" -o cmd/shortener/shortener ./cmd/shortener
```

### Запуск
```bash
go run -ldflags="-X main.buildVersion=1.0.1 -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X main.buildCommit=$(git rev-parse HEAD)" ./cmd/shortener
```