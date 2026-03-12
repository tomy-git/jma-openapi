# ADR 002: OpenAPI Generation Strategy

- 日付: 2026-03-12
- 状態: Accepted

## Context
設計書と README で「spec-first」「oapi-codegen による型とサーバー生成」「OpenAPI UI 同梱」が明示されており、API 契約と実装の間に人手の差異が入らないよう統一された生成戦略が必要。

## Decision
- `openapi/openapi.yaml` を正本とし、`oapi-codegen` で `types`、`server`、`spec` の3種生成物を `internal/gen` に置く。
- `openapi/openapi.json` も併記し、 `openapi.yaml` との内容差分が生じないよう CI などで同期を維持する。
- API ドキュメントは Scalar を使って `/docs` で同梱し、UI は `openapi.yaml` を参照する形にする。

## Consequences
- `spec` と実装が生成コードを通じて厳密に同期するため、契約変更時に `oapi-codegen` の再実行が必須となる。
- JSON/YAML 両方を公開することで人間・機械双方の参照性が向上するが、差分検出と整合性維持を運用で確保する必要がある。
- Scalar で UI を提供するために生成済み仕様のパス `/docs` を露出する運用ルールが必要となる。

## Alternatives Considered
1. 手書きの構造体とハンドラ
   - 理由: 契約と実装の乖離が生まれやすく、将来的なエンドポイント追加時に検証コストが膨らむため却下。
2. 仕様を YAML のみで公開
   - 理由: 外部システムが JSON での利用を想定しており、JSON を生成物として同梱することで利便性と自動化の両立を図るため不採用。
