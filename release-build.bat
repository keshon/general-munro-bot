go build -ldflags "-s -w" -o app.exe src/main.go

COPY README.md .\release