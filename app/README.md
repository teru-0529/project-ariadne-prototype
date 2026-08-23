# Project ARIADNE Prototype - Application

Project ARIADNE Prototype のWindowsアプリケーション。

Wails v3 + Svelte 5 + TypeScript + Go で構成する。

## Structure

- `backend/services/` : Go backend services
- `frontend/src/` : Svelte frontend
- `frontend/bindings/` : Wails generated bindings
- `build/` : Wails build configuration
- `bin/` : build output (Git管理対象外)

## Development

```bash
./tools/dev.sh
```

Wails development mode でアプリケーションを起動する。

Frontend / Backend の変更はhot reloadされる。

## Build

```bash
./tools/build.sh
```

Windows実行ファイルを以下に生成する。

```text
bin/ariadne-prototype.exe
```

## Toolchain

- Go 1.27.0
- Node.js 24.19.0
- npm 11.17.0

Node.js / npm は Volta でバージョンを固定する。
