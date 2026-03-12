<!--
SPDX-FileCopyrightText: 2026 The jma-openapi contributors

SPDX-License-Identifier: MPL-2.0
-->

# ADR 001: Language and Router Selection

- 日付: 2026-03-12
- 状態: Accepted

## Context
`README.md` および設計書で「シンプルさ」「保守性」「OpenAPI spec-first」を重視した構成が宣言されており、安定した型システムと生成ツールのある言語に寄せる必要がある。JMA BOSAI JSON の多様なスキーマとオープンAPIとの整合を保つには、実装言語とルータが生成コードと馴染むことが前提となる。

## Decision
- 言語は Go に決定し、標準 `net/http` を基盤とした実装とする。
- ルータは `chi` を採用し、`oapi-codegen` の出力であるハンドラインターフェースと自然に接続できる形とする。

### 理由
1. Go は `README.md` で採用することが明示されており、静的型と軽量実行、Docker/Cloud Run に合わせたビルドを容易にする。
2. `chi` は `net/http` に近く、生成されたサーバインターフェースにマウントしやすいため、シンプルで保守しやすいルーティング層を実現できる。
3. Go + `chi` の組み合わせは Cloud Run におけるスケーラブルな HTTP サービスにマッチし、長期的な運用負荷を低く保てる。

## Consequences
- 長期的には Go の型厳格性と `chi` の軽量ルーティングによって仕様変更の影響範囲を予測しやすくなる。
- JavaScript や Python を選ばなかったことで、既存の Node/Python エコシステムに慣れている開発者には学習コストが発生する。
- `chi` の軽量性により、将来的により豊富なルーティング機能が必要になった場合はミドルウェア自作などの追加作業が発生する可能性がある。

## Alternatives Considered
1. Node.js + Express
   - 特徴: 豊富なライブラリ、動的言語で構成の自由度が高い。
   - 理由: 静的仕様（OpenAPI）との整合やバイナリ展開を考えると型のない実装が障壁となり、Cloud Run での起動コストも増えると判断。
2. Python + FastAPI
   - 特徴: OpenAPI 自動生成能力が高く、非同期実装も容易。
   - 理由: ランタイム依存とパッケージ管理・スピンアップ時間を短く保てないこと、および README/設計書との整合性を欠くため採用を見送った。
