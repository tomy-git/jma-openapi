<!--
SPDX-FileCopyrightText: 2026 The jma-openapi contributors

SPDX-License-Identifier: MPL-2.0
-->

# ADR 003: Deployment and Logging

- 日付: 2026-03-12
- 状態: Accepted

## Context
初期デプロイ先として Cloud Run を選び、Dockerfile でビルドしたコンテナを載せる方針が README と設計書で共有されている。構造化ログには Go 標準の `log/slog` を使い、ローカルと本番で handler を切り替える必要がある。

## Decision
- Cloud Run を実行基盤とし、Dockerfile で multi-stage ビルドを行ってバイナリを作成、`deploy/cloudrun/service.yaml` で `min-instances: 0` を含む設定を管理する。
- ロギングは `log/slog` を採用し、ローカルでは text handler 本番では JSON handler を切り替え、共通のフィールドセット (`request_id`、`path`、`method`、`status`、`latency_ms` など) を維持する。

## Consequences
- Cloud Run 固有の動作に合わせてアプリはステートレスで起動時間を短く保つ必要があり、Dockerfile は `CGO_ENABLED=0` や `-trimpath` など最小化設定を含める。
- `log/slog` による構造化ログは運用側でフィールドを揃える意味で効果的だが、ローカルと本番で handler を切り替える制御を明文化しなければならない。
- `deploy/cloudrun/service.yaml` に `min-instances: 0` を明示することでスケールダウンとコスト削減を明文化できるが、Cold Start を考慮した監視も必要。

## Alternatives Considered
1. Cloud Run ではなく Compute Engine や GKE
   - 理由: 運用負荷が高く、README が示す「シンプルで保守しやすい」方針と乖離するため却下。
2. logrus などのサードパーティロガー
   - 理由: Go 標準の `log/slog` で十分な構造化ログが得られる上、依存を増やさずビルド・デプロイを簡素化できるため採用しなかった。
