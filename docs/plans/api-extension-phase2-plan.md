<!--
SPDX-FileCopyrightText: 2026 The jma-openapi contributors

SPDX-License-Identifier: MPL-2.0
-->

# API 拡張第2フェーズ 開発計画書

## 1. 概要

### 1.1 目的

第1フェーズで追加した area 検索と forecast area 個別取得を、探索用途とクライアント実装用途の両方で使いやすい形へ拡張する。具体的には、forecast area の専用一覧・専用個別・reverse-lookup・latest・timeseries API を追加し、`GET /v1/areas` には地名検索モードを導入する。`openapi/openapi.yaml` を正本とする spec-first フローは維持し、段階的に実装できる計画へ落とし込む。

### 1.2 スコープ

- 対象:
  - `GET /v1/areas` への `nameMatchMode` の追加
  - `GET /v1/forecasts/{officeCode}/areas` の追加
  - `GET /v1/forecasts/{officeCode}/weather-areas` の追加
  - `GET /v1/forecasts/{officeCode}/weather-areas/{areaCode}` の追加
  - `GET /v1/forecasts/{officeCode}/temperature-areas` の追加
  - `GET /v1/forecasts/{officeCode}/temperature-areas/{areaCode}` の追加
  - `GET /v1/forecasts/{officeCode}/areas:resolve` の追加
  - `GET /v1/forecasts/{officeCode}/areas/{areaCode}/latest` の追加
  - `GET /v1/forecasts/{officeCode}/areas/{areaCode}/timeseries` の追加
  - 追加 API に伴う OpenAPI、生成コード、handler / usecase / mapper / client / test / README 更新
- 非対象:
  - JMA 上流 endpoint 自体の追加
  - DB、永続キャッシュ、認証認可、監視基盤の新設
  - かな変換、外部辞書、N-gram などを使う高度な曖昧検索
  - `GET /v1/forecasts/{officeCode}` の全体レスポンス再設計
  - 多言語検索や `enName` 検索の本格対応

### 1.3 成功条件

- [x] `GET /v1/areas` で `nameMatchMode` に `exact` / `prefix` / `partial` / `suggested` を指定できる。
- [x] `GET /v1/areas` の mode 未指定時は後方互換のため `exact` と同じ挙動を維持する。
- [x] `GET /v1/areas` の `suggested` は空でない候補を relevance 順で返す。
- [x] `GET /v1/forecasts/{officeCode}/areas` が `weather` / `temperature` の統合一覧を返す。
- [x] `GET /v1/forecasts/{officeCode}/weather-areas` と `GET /v1/forecasts/{officeCode}/temperature-areas` が種別別の一覧を返す。
- [x] `GET /v1/forecasts/{officeCode}/weather-areas/{areaCode}` が `weatherArea` 専用レスポンスを返す。
- [x] `GET /v1/forecasts/{officeCode}/temperature-areas/{areaCode}` が `temperatureArea` 専用レスポンスを返す。
- [x] `GET /v1/forecasts/{officeCode}/areas:resolve` が area 名から候補 code 一覧を返す。
- [x] `GET /v1/forecasts/{officeCode}/areas/{areaCode}/latest` が対象 area の代表 1 エントリを返す。
- [x] `GET /v1/forecasts/{officeCode}/areas/{areaCode}/timeseries` が対象 area の時系列全体を返す。
- [x] 既存 `GET /v1/forecasts/{officeCode}/areas/{areaCode}` の互換方針が OpenAPI と README に明記されている。
- [x] `oapi-codegen` 再生成、`go test ./...`、`golangci-lint run` が成功する。

## 2. 前提と制約

- 期限:
  - 未指定
- 技術制約:
  - 既存構成は `openapi/openapi.yaml` を正本とする spec-first
  - 公開 OpenAPI は `openapi/openapi.yaml` と `openapi/openapi.json` の両方を同期して配信する
  - 生成コードは `oapi-codegen` により `internal/gen` に出力される
  - `tests/fixtures/area.json` には `name` / `enName` / `officeName` / `parent` / `children` はあるが、別名辞書やかな情報はない
  - `tests/fixtures/forecast-130000.json` では `weatherAreas` と `temperatureAreas` は別集合で code 体系も異なる
  - `internal/clients/jma_client.go` には full forecast 用取得経路と temperature 欠落許容の weather 用取得経路がすでにある
  - `internal/gen/spec.gen.go` は `oapi-codegen` 再生成で repository 固有補助関数が消えるため追補が必要になる
  - 第2フェーズでは `officeName` 検索は既存の exact のみ維持し、mode 拡張は `name` に限定する
- 依存チーム/外部要因:
  - JMA upstream の `area.json` と forecast JSON の構造維持
  - `oapi-codegen` 実行環境がローカルで利用可能であること

## 3. 実装方針

### 3.1 アーキテクチャ方針

- `weatherArea` と `temperatureArea` は schema 差が大きいため、専用 endpoint を主導線にする。
- `GET /v1/forecasts/{officeCode}/areas/{areaCode}` は互換 endpoint として残し、内部では専用 endpoint と同じ抽出ロジックを再利用する。
- forecast area 一覧 API は `統合一覧` と `種別別一覧` を両方提供し、探索用途と実装用途を分ける。
- `areas:resolve` は一覧 API とは責務を分け、検索・逆引き専用 endpoint とする。
- `latest` は「時刻最大」ではなく「時系列先頭要素」を返す暫定仕様とし、fixture と upstream 順序を正とする。
- `suggested` は第2フェーズでは軽量実装に限定し、別名辞書やかな変換には踏み込まない。

### 3.2 変更戦略

- OpenAPI を先に更新して生成コードを再生成する一括更新とする。
- 実装は `areas` 検索モード、forecast area 専用一覧・専用個別、resolve / latest / timeseries の順で進める。
- `suggested` の仕様は先に OpenAPI と計画書へ明記し、実装差分が出ないようにする。
- 既存 endpoint は急に削除せず、専用 endpoint を docs の主導線として案内する。

## 4. タスク分解

### [x] フェーズ1: OpenAPI 契約と互換方針の確定

- [x] タスク1.1: `areas` 検索モードを OpenAPI に追加する
  - 目的:
    - `GET /v1/areas` の検索モード契約を先に固定する。
  - 対象ファイル:
    - `openapi/openapi.yaml`
    - `openapi/openapi.json`
  - 完了条件:
    - `nameMatchMode` が enum 付きで追加されている。
    - mode 未指定時の既定値と `suggested` の説明が明記されている。

- [x] タスク1.2: forecast area 追加 endpoint 群を OpenAPI に追加する
  - 目的:
    - 一覧・専用個別・resolve・latest・timeseries の契約を固定する。
  - 対象ファイル:
    - `openapi/openapi.yaml`
    - `openapi/openapi.json`
  - 完了条件:
    - `areas` / `weather-areas` / `temperature-areas` / `areas:resolve` / `latest` / `timeseries` が定義されている。
    - 既存 `GET /v1/forecasts/{officeCode}/areas/{areaCode}` の互換方針が説明されている。

- [x] タスク1.3: 生成コードを再生成する
  - 目的:
    - OpenAPI と Go 実装側 interface / 型を整合させる。
  - 対象ファイル:
    - `internal/gen/server.gen.go`
    - `internal/gen/types.gen.go`
    - `internal/gen/spec.gen.go`
  - 完了条件:
    - `oapi-codegen` 3 コマンドが成功する。
    - 新規 query / path に対応した型と interface が生成される。
    - `internal/gen/spec.gen.go` の repository 固有補助関数が維持される。

### [x] フェーズ2: `areas` 検索モードの実装

- [x] タスク2.1: handler / usecase に検索モードを受け渡す
  - 目的:
    - `name` の mode をアプリ層へ渡せるようにする。
  - 対象ファイル:
    - `internal/handlers/areas.go`
    - `internal/usecases/list_areas.go`
    - 必要に応じて `internal/handlers/server.go`
  - 完了条件:
    - `nameMatchMode` を含む filter struct が定義されている。
    - mode 未指定時に `exact` が適用される。

- [x] タスク2.2: `AreaMapper` に exact / prefix / partial / suggested を実装する
  - 目的:
    - `name` に検索モードを適用できるようにする。
  - 対象ファイル:
    - `internal/mappers/area_mapper.go`
    - 必要に応じて新規 matcher 定義ファイル
  - 完了条件:
    - `name` に対して `exact` / `prefix` / `partial` / `suggested` が実装される。
    - `suggested` は score 降順、同 score は code 昇順で返る。
    - 既存の `parent` / `child` / `officeName` との AND 条件が維持される。

- [x] タスク2.3: `areas` 検索モードのテストを追加する
  - 目的:
    - query 契約と候補順を固定する。
  - 対象ファイル:
    - `internal/mappers/area_mapper_test.go`
    - `internal/handlers/areas_test.go`
  - 完了条件:
    - `exact` / `prefix` / `partial` / `suggested` の正常系がテスト化される。
    - mode 不正の 400 と 0 件時 `200 + items: []` が確認できる。

### [x] フェーズ3: forecast area 一覧・専用個別 API の実装

- [x] タスク3.1: forecast area 一覧レスポンスの mapper を追加する
  - 目的:
    - 統合一覧と種別別一覧を共通ロジックで生成する。
  - 対象ファイル:
    - `internal/mappers/forecast_mapper.go`
    - 必要に応じて新規型定義ファイル
  - 完了条件:
    - `kind` / `code` / `name` の統合一覧が生成できる。
    - weather / temperature の種別別一覧が生成できる。

- [x] タスク3.2: 一覧・専用個別 endpoint の usecase / handler を追加する
  - 目的:
    - weather / temperature の専用導線を実装する。
  - 対象ファイル:
    - `internal/usecases/get_forecast.go`
    - `internal/handlers/forecasts.go`
    - `internal/clients/jma_client.go`
  - 完了条件:
    - `areas` / `weather-areas` / `temperature-areas` 一覧が返せる。
    - weather / temperature の専用個別取得が返せる。
    - `OFFICE_NOT_FOUND` と area 未存在の 404 が整理される。

- [x] タスク3.3: 一覧・専用個別 API のテストを追加する
  - 目的:
    - 一覧内容と専用 schema 契約を固定する。
  - 対象ファイル:
    - `internal/mappers/forecast_mapper_test.go`
    - `internal/handlers/forecasts_test.go`
    - 必要に応じて `internal/usecases/get_forecast_test.go`
  - 完了条件:
    - 統合一覧、種別別一覧、weather 個別、temperature 個別がテスト化される。
    - 既存汎用 endpoint の互換性が確認できる。

### [x] フェーズ4: resolve / latest / timeseries API の実装

- [x] タスク4.1: resolve 用の検索ロジックを実装する
  - 目的:
    - area 名から code 候補を逆引きできるようにする。
  - 対象ファイル:
    - `internal/mappers/forecast_mapper.go`
    - `internal/usecases/get_forecast.go`
    - `internal/handlers/forecasts.go`
  - 完了条件:
    - `q` / `kind` / `matchMode` を解釈して候補一覧を返せる。
    - `suggested` の順位が fixture で固定される。

- [x] タスク4.2: latest / timeseries 抽出処理を追加する
  - 目的:
    - 軽量取得と時系列取得を専用 endpoint として提供する。
  - 対象ファイル:
    - `internal/mappers/forecast_mapper.go`
    - `internal/usecases/get_forecast.go`
    - `internal/handlers/forecasts.go`
  - 完了条件:
    - `latest` が時系列先頭要素を返す。
    - `timeseries` が対象 area の全要素を返す。
    - weather / temperature の両系統で成立する。

- [x] タスク4.3: resolve / latest / timeseries のテストを追加する
  - 目的:
    - 補助 endpoint の仕様を固定する。
  - 対象ファイル:
    - `internal/mappers/forecast_mapper_test.go`
    - `internal/handlers/forecasts_test.go`
    - 必要に応じて `internal/usecases/get_forecast_test.go`
  - 完了条件:
    - reverse-lookup、latest、timeseries の正常系 / 主要異常系がテスト化される。
    - `latest` の抽出規則が fixture ベースで固定される。

### [x] フェーズ5: 回帰確認とドキュメント整備

- [x] タスク5.1: 全体テストと生成物整合を確認する
  - 目的:
    - 仕様、生成コード、実装、テストの整合を最終確認する。
  - 対象ファイル:
    - `openapi/openapi.yaml`
    - `openapi/openapi.json`
    - `internal/gen/*`
    - `internal/**/*_test.go`
  - 完了条件:
    - `oapi-codegen` 再生成、`go test ./...`、`golangci-lint run` が成功する。
    - `/docs`、`/openapi.yaml`、`/openapi.json` に追加契約が表示される。

- [x] タスク5.2: README と関連文書を更新する
  - 目的:
    - 追加 API と検索モードの利用導線を文書化する。
  - 対象ファイル:
    - `README.md`
    - `README-JA.md`
    - 必要に応じて関連設計ドキュメント
  - 完了条件:
    - 新規 endpoint の利用例が反映される。
    - `suggested` と `latest` の定義が文書に明記される。

## 5. 影響範囲分析

### 5.1 変更対象

| 区分（直接影響/間接影響/テスト影響） | ファイル                                   | 変更種別 | 変更内容                                              | 関連タスク          |
| ------------------------------------ | ------------------------------------------ | -------- | ----------------------------------------------------- | ------------------- |
| 直接影響                             | `openapi/openapi.yaml`                     | 修正     | query / path / schema 追加                            | タスク1.1, 1.2      |
| 直接影響                             | `openapi/openapi.json`                     | 修正     | YAML と同期した公開 spec 更新                         | タスク1.1, 1.2, 1.3 |
| 直接影響                             | `internal/gen/server.gen.go`               | 生成更新 | server interface と params 更新                       | タスク1.3           |
| 直接影響                             | `internal/gen/types.gen.go`                | 生成更新 | 新規 schema 型の追加                                  | タスク1.3           |
| 直接影響                             | `internal/gen/spec.gen.go`                 | 生成更新 | embedded spec 更新と補助関数維持                      | タスク1.3           |
| 直接影響                             | `internal/handlers/areas.go`               | 修正     | 検索モード query の受け渡し                           | タスク2.1           |
| 直接影響                             | `internal/usecases/list_areas.go`          | 修正     | mode を含む filter の導入                             | タスク2.1           |
| 直接影響                             | `internal/mappers/area_mapper.go`          | 修正     | `name` 向け exact / prefix / partial / suggested 実装 | タスク2.2           |
| 直接影響                             | `internal/handlers/forecasts.go`           | 修正     | 追加 endpoint の handler 実装                         | タスク3.2, 4.1, 4.2 |
| 直接影響                             | `internal/usecases/get_forecast.go`        | 修正     | 一覧 / 専用個別 / resolve / latest / timeseries 追加  | タスク3.2, 4.1, 4.2 |
| 直接影響                             | `internal/mappers/forecast_mapper.go`      | 修正     | area 一覧・resolve・latest・timeseries の生成         | タスク3.1, 4.1, 4.2 |
| 直接影響                             | `internal/clients/jma_client.go`           | 修正     | 共通 forecast 取得経路の調整                          | タスク3.2           |
| テスト影響                           | `internal/mappers/area_mapper_test.go`     | 修正     | 検索モード別の fixture test 追加                      | タスク2.3           |
| テスト影響                           | `internal/handlers/areas_test.go`          | 修正     | query parameter の HTTP テスト追加                    | タスク2.3           |
| テスト影響                           | `internal/mappers/forecast_mapper_test.go` | 修正     | 一覧 / resolve / latest / timeseries test 追加        | タスク3.3, 4.3      |
| テスト影響                           | `internal/handlers/forecasts_test.go`      | 修正     | path/query の HTTP テスト追加                         | タスク3.3, 4.3      |
| テスト影響                           | `internal/usecases/get_forecast_test.go`   | 新規候補 | 追加 usecase の異常系 test                            | タスク3.3, 4.3      |
| 間接影響                             | `README.md`                                | 修正     | 利用例と互換方針更新                                  | タスク5.2           |
| 間接影響                             | `README-JA.md`                             | 修正     | 利用例と互換方針更新                                  | タスク5.2           |

### 5.2 依存関係

- 依存先:
  - JMA `area.json`
  - JMA forecast JSON
  - `oapi-codegen`
  - `chi` router と生成 server interface
- 影響先:
  - API 利用者の area 検索方法
  - forecast area code の探索・取得導線
  - OpenAPI UI (`/docs`) に表示される契約

### 5.3 テスト影響

- 追加:
  - `areas` 検索モード別の unit test
  - `areas` query parameter の HTTP テスト
  - forecast area 一覧・専用個別の unit / HTTP テスト
  - resolve / latest / timeseries の unit / HTTP テスト
  - 既存汎用 endpoint 互換性の回帰テスト
- 更新:
  - `forecast_mapper_test.go` の期待値拡張
  - `forecasts_test.go` の path/query ケース拡張
- 不要理由（不要な場合のみ）:
  - E2E は既存環境がないため必須外とし、`httptest` で代替する。

## 6. リスク評価

| リスク                                                                     | 影響度(1-3) | 発生確率(1-3) | スコア | 対策                                                                  | 検証                          | 対応タスク               |
| -------------------------------------------------------------------------- | ----------- | ------------- | ------ | --------------------------------------------------------------------- | ----------------------------- | ------------------------ |
| `suggested` の意味が曖昧で期待順位がぶれる                                 | 3           | 3             | 9      | score 規則を OpenAPI と計画書へ明記し、fixture で順位テストを固定する | mapper test / HTTP test       | タスク1.1, 2.2, 2.3      |
| 既存汎用 endpoint と専用 endpoint の責務が重複する                         | 3           | 2             | 6      | 汎用 endpoint は互換維持、専用 endpoint を主導線にする                | OpenAPI / README / 回帰テスト | タスク1.2, 3.2, 5.2      |
| weather / temperature の area 種別判定を誤る                               | 3           | 2             | 6      | 種別別 mapper を分け、統合 endpoint は kind を明示する                | mapper test / handler test    | タスク3.1, 3.3, 4.2, 4.3 |
| `latest` の定義が downstream 期待とずれる                                  | 2           | 2             | 4      | 第2フェーズでは「時系列先頭要素」を仕様で固定する                     | fixture test / README         | タスク1.2, 4.2, 5.2      |
| 生成コード再生成で repository 固有補助関数が消える                         | 2           | 2             | 4      | `internal/gen/spec.gen.go` の追補方針を維持する                       | build / test                  | タスク1.3, 5.1           |
| `areas:resolve` に score や matchedBy を出さないことで将来拡張が難しくなる | 2           | 2             | 4      | schema を最小構成に留め、未確定事項として明示する                     | 計画レビュー                  | タスク1.2, 4.1           |

## 7. 検証計画

- 単体テスト:
  - `go test ./...`
  - `AreaMapper` の mode 別テスト
  - `ForecastMapper` の一覧 / resolve / latest / timeseries テスト
- 結合テスト:
  - `httptest` ベースで `areas` / `forecasts` の path/query 配線を確認
  - `OFFICE_NOT_FOUND` / `FORECAST_AREA_NOT_FOUND` / bad request の異常系テスト
- E2E/手動確認:
  - `/docs` で追加 endpoint と enum が表示されることを確認
  - `/openapi.yaml` と `/openapi.json` の配信内容が一致することを確認
  - 代表ケースの `curl` 例を README で確認
- 性能/セキュリティ確認:
  - 性能: 新規 DB や外部 I/O は増やさないため専用負荷試験は必須外
  - セキュリティ: 検索 query は文字列比較のみで、機微情報出力を増やさないことを確認

## 8. リリース・運用

- リリース手順:
  - `openapi/openapi.yaml` と `openapi/openapi.json` を更新
  - `oapi-codegen` 再生成
  - 実装更新
  - `go test ./...`
  - `golangci-lint run`
  - README 更新
- ロールバック手順:
  - 追加 endpoint と検索モード拡張を含む変更一式を同一コミット単位で戻す
  - OpenAPI、生成コード、実装、README をまとめて巻き戻す
- 監視/アラート:
  - 新規監視は追加しない
  - `OFFICE_NOT_FOUND` / `FORECAST_AREA_NOT_FOUND` / mode 不正の 400 を既存ログで確認対象とする

## 9. 未確定事項

- `suggested` の response に `score` や `matchedBy` を含めるか
- `GET /v1/areas` に `limit` を追加するか
- `officeName` に対して将来的に `nameMatchMode` 相当を拡張するか

## 10. セルフレビュー結果

| 課題                                                                                | 重要度 | 対応方針                                                             |
| ----------------------------------------------------------------------------------- | ------ | -------------------------------------------------------------------- |
| 旧版は planner テンプレートの章立てとチェックボックス形式に沿っていなかった         | Low    | 本計画書でテンプレート準拠に修正済み                                 |
| `GET /v1/forecasts/{officeCode}/areas` の統合一覧 endpoint が成功条件から漏れていた | Medium | 成功条件とスコープへ反映済み                                         |
| `suggested` の定義が曖昧で実装とテストがずれる恐れがあった                          | Medium | `name` 向け正規化 + score 順の軽量実装に限定し、未確定事項を分離済み |
| 既存汎用 endpoint の扱いが docs 上で曖昧だった                                      | Medium | 互換 endpoint として残す方針を実装方針と成功条件へ反映済み           |
| `latest` の定義が「最新時刻」か「先頭要素」か曖昧だった                             | Medium | 第2フェーズでは先頭要素とする仮定を前提と制約へ明記済み              |

## 11. ユーザー確認事項（必要な場合のみ）

1. `suggested` の返却項目に検索補助情報を含めるか（重要度: Medium）
   - 背景: 現計画は最小構成の `code` / `name` を前提にしていますが、将来の UI 実装では `score` や `matchedBy` があると扱いやすくなります。
   - 選択肢: `A) 第2フェーズは補助情報なしで最小構成にする / B) 第2フェーズから score または matchedBy を response に含める`
   - 推奨: `A`
2. `GET /v1/areas` に件数制御を入れるか（重要度: Medium）
   - 背景: `suggested` や `partial` は候補が増えやすく、将来的に `limit` が必要になる可能性があります。
   - 選択肢: `A) 第2フェーズでは limit なし / B) 第2フェーズから limit を追加する`
   - 推奨: `A`
3. `officeName` の検索モード拡張を第2フェーズに含めるか（重要度: Medium）
   - 背景: 本計画は未回答時の仮定として `officeName` は exact のみ維持に寄せています。ここを同時対応すると実装とテストが増えます。
   - 選択肢: `A) 第2フェーズは name のみ mode 拡張し、officeName は exact のみ維持する / B) 第2フェーズから officeName にも同じ mode を追加する`
   - 推奨: `A`

## 12. 進捗管理

### 12.1 状態ボード

| タスク                                                                    | 状態 | 担当                | 最終更新             | 次アクション | ブロッカー |
| ------------------------------------------------------------------------- | ---- | ------------------- | -------------------- | ------------ | ---------- |
| タスク1.1 `areas` 検索モードを OpenAPI に追加する                         | DONE | operational-command | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク1.2 forecast area 追加 endpoint 群を OpenAPI に追加する             | DONE | operational-command | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク1.3 生成コードを再生成する                                          | DONE | operational-command | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク2.1 handler / usecase に検索モードを受け渡す                        | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク2.2 `AreaMapper` に exact / prefix / partial / suggested を実装する | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク2.3 `areas` 検索モードのテストを追加する                            | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク3.1 forecast area 一覧レスポンスの mapper を追加する                | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク3.2 一覧・専用個別 endpoint の usecase / handler を追加する         | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク3.3 一覧・専用個別 API のテストを追加する                           | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク4.1 resolve 用の検索ロジックを実装する                              | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク4.2 latest / timeseries 抽出処理を追加する                          | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク4.3 resolve / latest / timeseries のテストを追加する                | DONE | worker              | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク5.1 全体テストと生成物整合を確認する                                | DONE | reviewer            | 2026-03-12 16:40 JST | 完了         | なし       |
| タスク5.2 README と関連文書を更新する                                     | DONE | operational-command | 2026-03-12 16:40 JST | 完了         | なし       |

### 12.2 判断結果

- 2026-03-12 16:05 JST:
  - 実行基準は `docs/plans/api-extension-phase2-plan.md` とした。
  - サブエージェントが使えないため、`operational-command` が進捗管理を代行する。
  - 未回答時の仮定として、第2フェーズでは `officeName` の検索モード拡張は見送り、`nameMatchMode` を優先する。
- 2026-03-12 16:40 JST:
  - `GET /v1/areas` に `nameMatchMode=exact|prefix|partial|suggested` を追加し、mode 不正時の 400 を実装した。
  - forecast area の統合一覧 / 種別別一覧 / 専用個別 / resolve / latest / timeseries endpoint を追加した。
  - 既存 `GET /v1/forecasts/{officeCode}/areas/{areaCode}` は互換 endpoint として維持した。
  - `go test ./...` と `golangci-lint run` の成功を確認した。
