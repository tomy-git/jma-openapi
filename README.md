# jma-openapi

Unofficial OpenAPI wrapper for the Japan Meteorological Agency (JMA) BOSAI JSON data.

気象庁（JMA）防災情報JSONデータのための非公式OpenAPIラッパーです。

## Overview

`jma-openapi` is an unofficial API wrapper project that provides a simpler and more developer-friendly interface for the Japan Meteorological Agency (JMA) BOSAI JSON endpoints.

The project normalizes JMA weather and disaster-related data and exposes it through a consistent REST API defined by an OpenAPI specification. This makes it easier for application developers to consume JMA data without dealing with the complexity of the original JSON structures.

`jma-openapi` は、気象庁（JMA）の防災情報JSONエンドポイントに対して、よりシンプルで開発者が扱いやすいインターフェースを提供する非公式APIラッパープロジェクトです。

本プロジェクトは、JMAの気象・防災関連データを正規化し、OpenAPI仕様に基づいたREST APIとして提供することで、元のJSON構造の複雑さを意識せずに利用できるようにすることを目的としています。

This project is not affiliated with, endorsed by, or maintained by the Japan Meteorological Agency.

本プロジェクトは気象庁とは一切関係がなく、承認・支援・保守も受けていません。

---

## Goals

- Provide a simple wrapper for JMA BOSAI JSON endpoints  
- JMA BOSAI JSONエンドポイントのためのシンプルなラッパーを提供する

- Offer a stable and developer-friendly REST API interface  
- 安定的で開発者にとって扱いやすいREST APIインターフェースを提供する

- Publish an OpenAPI specification for implemented endpoints  
- 実装済みエンドポイントに対応するOpenAPI仕様を公開する

- Reduce the complexity of consuming raw JMA JSON structures  
- 生のJMA JSON構造を利用する際の複雑さを軽減する

- Serve as a foundation for applications that use Japanese weather and disaster data  
- 日本の気象・防災データを利用するアプリケーションの基盤となる

---

## Architecture

The project follows a **spec-first API design** approach.

OpenAPI specifications define the API contract first, and server code is generated from the specification using Go tooling.

API handlers then implement the generated interfaces.

プロジェクトは **spec-first API設計** を採用します。

まず OpenAPI 仕様で API コントラクトを定義し、その仕様から Go コードを自動生成します。

その後、生成されたインターフェースを実装する形で API ハンドラを構築します。

---

## Technology Stack

- Language: **Go**
- Router: **chi**
- OpenAPI: **oapi-codegen**
- HTTP layer: **net/http (Go standard library)**
- Testing: **go test**
- Lint: **golangci-lint**
- API docs UI: **bundled**
- Deployment: **Cloud Run**
- Logging: **log/slog**

技術スタック

- 言語: **Go**
- ルータ: **chi**
- OpenAPIコード生成: **oapi-codegen**
- HTTP基盤: **Go標準の net/http**
- テスト: **go test**
- Lint: **golangci-lint**
- APIドキュメントUI: **同梱**
- デプロイ: **Cloud Run**
- ロガー: **log/slog**

This combination provides:

- High performance
- Minimal framework dependency
- Strong compatibility with Go ecosystem
- Long-term maintainability

この構成は以下を重視しています。

- 高いパフォーマンス
- フレームワーク依存を最小化
- Goエコシステムとの高い互換性
- 長期的な保守性

---

## Development Policy

- Follow **GitHub Flow** for development
- Keep `main` always in a working state
- Use **OpenAPI specification as the source of truth**
- Generate server interfaces using `oapi-codegen`
- Generate types, server interfaces, and embedded spec bindings into dedicated generated files
- Implement business logic in handlers separate from generated code
- Bundle an OpenAPI UI for human-friendly API reference
- Run tests with `go test`
- Run lint with `golangci-lint`
- Target Cloud Run as the initial deployment platform
- Use `log/slog` for structured logging
- Build and deploy the service with a `Dockerfile`
- Configure Cloud Run minimum instances as `0` for the initial release
- Keep the architecture simple and modular

開発方針

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

## Scope

The initial scope of the project includes the following:

初期スコープには以下を含みます。

- Health check endpoint  
- ヘルスチェックエンドポイント

- Forecast endpoint  
- 天気予報エンドポイント

- Area metadata endpoint  
- 地域メタデータエンドポイント

- Basic OpenAPI specification  
- 基本的なOpenAPI仕様

Additional endpoints may be added later for:

- Overview forecasts
- Weather warnings
- Earthquake information
- Tsunami alerts
- Other JMA BOSAI data resources

将来的には以下のエンドポイント追加も想定しています。

- 概況予報
- 警報・注意報
- 地震情報
- 津波情報
- その他のJMA BOSAI関連データ

---

## Disclaimer

This is an unofficial project. The upstream JMA BOSAI JSON format and endpoint behavior may change at any time. Compatibility is not guaranteed.

これは非公式プロジェクトです。上流のJMA BOSAI JSONの形式やエンドポイントの挙動は、いつでも変更される可能性があります。互換性は保証されません。

---

## License

This project is licensed under the Mozilla Public License 2.0 (MPL-2.0).

本プロジェクトは Mozilla Public License 2.0（MPL-2.0）のもとで提供されます。
