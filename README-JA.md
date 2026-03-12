<!--
SPDX-FileCopyrightText: 2026 The jma-openapi contributors

SPDX-License-Identifier: MPL-2.0
-->

# jma-openapi

本ファイルは `README.md` の参考訳です。正本は英語版の [`README.md`](README.md) であり、仕様や運用判断が食い違う場合は `README.md` を優先してください。

気象庁（JMA）防災情報JSONデータのための非公式OpenAPIラッパーです。

## 概要

`jma-openapi` は、気象庁（JMA）の防災情報JSONエンドポイントに対して、よりシンプルで開発者が扱いやすいインターフェースを提供する非公式APIラッパープロジェクトです。

本プロジェクトは、JMAの気象・防災関連データを正規化し、OpenAPI仕様に基づいたREST APIとして提供することで、元のJSON構造の複雑さを意識せずに利用できるようにすることを目的としています。

本プロジェクトは気象庁とは一切関係がなく、承認・支援・保守も受けていません。

---

## 目的

- JMA BOSAI JSONエンドポイントのためのシンプルなラッパーを提供する
- 安定的で開発者にとって扱いやすいREST APIインターフェースを提供する
- 実装済みエンドポイントに対応するOpenAPI仕様を公開する
- 生のJMA JSON構造を利用する際の複雑さを軽減する
- 日本の気象・防災データを利用するアプリケーションの基盤となる

---

## アーキテクチャ

プロジェクトは **spec-first API設計** を採用します。

まず OpenAPI 仕様で API コントラクトを定義し、その仕様から Go コードを自動生成します。

その後、生成されたインターフェースを実装する形で API ハンドラを構築します。

---

## 技術スタック

- 言語: **Go**
- ルータ: **chi**
- OpenAPIコード生成: **oapi-codegen**
- HTTP基盤: **Go標準の `net/http`**
- テスト: **go test**
- Lint: **golangci-lint**
- APIドキュメントUI: **同梱**
- デプロイ: **Cloud Run**
- ロガー: **log/slog**

この構成は以下を重視しています。

- 高いパフォーマンス
- フレームワーク依存を最小化
- Goエコシステムとの高い互換性
- 長期的な保守性

---

## 開発方針

- 開発フローには **GitHub Flow** を採用する
- `main` ブランチは常に動作する状態を保つ
- **OpenAPI仕様を唯一のAPI定義として扱う**
- `oapi-codegen` を用いてサーバーコードを生成する
- 型定義、サーバーインターフェース、仕様埋め込みコードは生成ファイルへ分離する
- 生成コードと実装コードを分離する
- 人が参照しやすい OpenAPI UI を同梱する
- テストは `go test` で実行する
- Lint は `golangci-lint` で実行する
- 初期デプロイ先は Cloud Run とする
- 構造化ログには `log/slog` を利用する
- サービスのビルドとデプロイには `Dockerfile` を用いる
- 初期リリースの Cloud Run 最小インスタンス数は `0` とする
- シンプルで保守しやすいアーキテクチャを維持する

---

## スコープ

初期スコープには以下を含みます。

- ヘルスチェックエンドポイント
- 天気予報エンドポイント
- 地域メタデータエンドポイント
- 基本的なOpenAPI仕様

将来的には以下のエンドポイント追加も想定しています。

- 概況予報
- 警報・注意報
- 地震情報
- 津波情報
- その他のJMA BOSAI関連データ

---

## はじめに

### 前提条件

- Go `1.26.1`
- Docker
- `oapi-codegen`
- `golangci-lint`

`mise` を使う場合は、リポジトリ直下の `mise.toml` で Go バージョンを揃えられます。

### ローカル起動

依存関係を解決する。

```bash
go mod tidy
```

テストと lint を実行する。

```bash
go test ./...
golangci-lint run
```

GitHub Actions では、pull request と `main` への push を契機に `go test ./...`、`golangci-lint`、REUSE によるライセンスチェック、`govulncheck ./...`、`go build ./...` を実行します。

サーバーを起動する。

```bash
go run ./cmd/server
```

確認先:

- `http://localhost:8080/healthz`
- `http://localhost:8080/v1/areas`
- `http://localhost:8080/v1/areas?name=東京都&officeName=気象庁&child=130010`
- `http://localhost:8080/v1/areas?name=東京&nameMatchMode=prefix`
- `http://localhost:8080/v1/areas/130000`
- `http://localhost:8080/v1/forecasts/130000`
- `http://localhost:8080/v1/forecasts/130000/areas/130010`
- `http://localhost:8080/v1/forecasts/130000/areas`
- `http://localhost:8080/v1/forecasts/130000/weather-areas`
- `http://localhost:8080/v1/forecasts/130000/temperature-areas`
- `http://localhost:8080/v1/forecasts/130000/areas:resolve?q=東京&matchMode=suggested`
- `http://localhost:8080/v1/forecasts/130000/areas/44132/latest`
- `http://localhost:8080/v1/forecasts/130000/areas/130010/timeseries`
- `http://localhost:8080/openapi.yaml`
- `http://localhost:8080/openapi.json`
- `http://localhost:8080/docs`

代表的な API 利用例:

```bash
curl 'http://localhost:8080/v1/areas?name=東京都&officeName=気象庁&child=130010'
curl 'http://localhost:8080/v1/areas?name=東京&nameMatchMode=prefix'
curl 'http://localhost:8080/v1/forecasts/130000/areas/130010'
curl 'http://localhost:8080/v1/forecasts/130000/areas/44132'
curl 'http://localhost:8080/v1/forecasts/130000/areas'
curl 'http://localhost:8080/v1/forecasts/130000/weather-areas'
curl 'http://localhost:8080/v1/forecasts/130000/temperature-areas'
curl 'http://localhost:8080/v1/forecasts/130000/areas:resolve?q=東京&matchMode=suggested'
curl 'http://localhost:8080/v1/forecasts/130000/areas/44132/latest'
curl 'http://localhost:8080/v1/forecasts/130000/areas/130010/timeseries'
```

### 環境変数

- `PORT`
  - 既定値: `8080`
- `SERVICE_NAME`
  - 既定値: `jma-openapi`
- `APP_VERSION`
  - 既定値: `dev`
- `JMA_BASE_URL`
  - 既定値: `https://www.jma.go.jp/bosai`
- `LOG_FORMAT`
  - 既定値: `text`
  - 構造化ログが必要な場合は `json`

### OpenAPI 生成コードの再生成

```bash
(cd openapi && oapi-codegen -config oapi-codegen-types.yaml openapi.yaml)
(cd openapi && oapi-codegen -config oapi-codegen-server.yaml openapi.yaml)
(cd openapi && oapi-codegen -config oapi-codegen-spec.yaml openapi.yaml)
```

### ライセンスヘッダー

SPDX ヘッダーと `.license` サイドカーファイルを追加または更新する。

```bash
mise run license-annotate
```

REUSE 準拠を検証する。

```bash
mise run license-lint
```

### Docker で起動

Docker Compose でサービスを起動する。

```bash
docker compose up --build
```

アプリは `http://localhost:8080` で利用できます。

`docker compose build` と `docker compose up --build` は、イメージビルド中に `go test ./...` と `golangci-lint run` を実行します。両方に成功した場合のみコンテナが起動します。

### Cloud Run

イメージをビルドして push する。

```bash
docker build -t gcr.io/<PROJECT_ID>/jma-openapi:latest .
docker push gcr.io/<PROJECT_ID>/jma-openapi:latest
```

サービス定義を適用する。

```bash
gcloud run services replace deploy/cloudrun/service.yaml --region <REGION>
```

---

## 免責事項

これは非公式プロジェクトです。上流のJMA BOSAI JSONの形式やエンドポイントの挙動は、いつでも変更される可能性があります。互換性は保証されません。

---

## ライセンス

本プロジェクトは Mozilla Public License 2.0（MPL-2.0）のもとで提供されます。
