# API 拡張第1フェーズ 開発計画書

## 1. 概要

### 1.1 目的

既存の `area` / `forecast` ラッパー API を、コードを知らない利用者でも探索しやすい形へ拡張する。対象は、`name` / `officeName` による code 逆引き、`child code` からの親 area 逆引き、`weatherAreas code` 単位での個別気象情報取得の 3 機能とし、`openapi/openapi.yaml` を正本とする spec-first の開発フローを維持したまま実装可能な計画へ落とし込む。

### 1.2 スコープ

- 対象:
  - `GET /v1/areas` への `name`、`officeName`、`child` query parameter 追加
  - `GET /v1/forecasts/{officeCode}/areas/{areaCode}` の追加
  - 追加 API に伴う OpenAPI、生成コード、handler / usecase / mapper / test の更新
  - README または関連設計ドキュメントの更新要否確認
- 非対象:
  - `weatherCode` など気象条件検索 API
  - JMA 上流 API の追加対応
  - DB、永続キャッシュ、認証認可、監視基盤の新設
  - 既存 `ForecastResponse` 全体の再設計

### 1.3 成功条件

- [x] `GET /v1/areas` で `parent`、`name`、`officeName`、`child` の AND 条件検索ができる。
- [x] `GET /v1/areas` の検索結果 0 件時は `200` と空の `items` を返す仕様が OpenAPI と実装で一致している。
- [x] `GET /v1/forecasts/{officeCode}/areas/{areaCode}` が `publishingOffice`、`reportDatetime`、`office`、対象 `weatherArea` を返し、temperature 系項目を含めない。
- [x] 追加 API の正常系と主要異常系が unit test と HTTP レベルテストで固定されている。
- [x] `oapi-codegen` 再生成と `go test ./...` が成功する。

## 2. 前提と制約

- 期限:
  - 未指定
- 技術制約:
  - 既存構成は `openapi/openapi.yaml` を正本とする spec-first
  - 公開 OpenAPI は `openapi/openapi.yaml` と `openapi/openapi.json` の両方を同期して配信する
  - 生成コードは `oapi-codegen` により `internal/gen` に出力される
  - `GET /v1/areas` は現状 `parent` のみを受け取り、`AreaUsecase.List` は単一引数構成である
  - `GET /v1/forecasts/{officeCode}` は JMA forecast JSON を取得後、`ForecastMapper` で `weatherAreas` と `temperatureAreas` へ正規化する
  - 現行 client は forecast JSON の `timeSeries[0..2]` を前提にしており、`timeSeries[2]` 欠落時は schema mismatch 扱いになる
  - `weatherAreas code` は office を跨いで一意と断定できないため、第1フェーズでは `officeCode` を path に含める
- 依存チーム/外部要因:
  - JMA upstream の `area.json` と forecast JSON の構造維持
  - `oapi-codegen` 実行環境がローカルで利用可能であること

## 3. 実装方針

### 3.1 アーキテクチャ方針

- `areas` の逆引きは新規エンドポイントを増やさず、既存 `GET /v1/areas` の query 拡張として扱う。
- `AreaUsecase.List` は query parameter 群を struct 化して受け取り、handler と usecase の引数増加を抑える。
- `weatherAreas` 個別取得は既存の full forecast 取得とは切り分け、weather 個別取得に必要な最小要件だけを扱う client / usecase 経路を追加する。
- `GET /v1/forecasts/{officeCode}/areas/{areaCode}` のレスポンスは専用 schema とし、少なくとも `office`、`publishingOffice`、`reportDatetime`、`weatherArea` を返す。
- 個別取得レスポンスに `temperature` 系項目は含めない。現行 mapper 上も `weatherAreas` と `temperatureAreas` は独立配列であり、1:1 対応を前提にできないためである。
- weather 個別取得では upstream の最小要件を `timeSeries[0]` 必須、`timeSeries[1]` は利用可能なら降水確率へ反映、`timeSeries[2]` は不要とする。

### 3.2 変更戦略

- OpenAPI を先に更新して生成コードを再生成する一括更新とする。
- 既存 `GET /v1/forecasts/{officeCode}` のレスポンス契約は変更しない。
- `GET /v1/areas` の追加 query は完全一致検索とし、部分一致や正規化比較は第1フェーズ外とする。
- 検索結果 0 件は一覧 API として自然な `200 + items: []` に統一する。
- 実装順は優先度順に `name / officeName`、`child`、`weatherAreas` とし、途中でスコープ圧縮が必要になった場合も優先順位を維持する。

## 4. タスク分解

### [x] フェーズ1: OpenAPI 契約の確定

- [x] タスク1.1: `areas` 検索条件を OpenAPI に追加する
  - 目的:
    - `GET /v1/areas` の検索仕様を先に固定する。
  - 対象ファイル:
    - `openapi/openapi.yaml`
    - `openapi/openapi.json`
  - 完了条件:
    - `name`、`officeName`、`child` query parameter が追加されている。
    - `0件時は 200 + items: []` が説明文またはレスポンス仕様で明確になっている。

- [x] タスク1.2: `weatherAreas` 個別取得 API の契約を追加する
  - 目的:
    - 実装前に新規 path と専用 schema を固定する。
  - 対象ファイル:
    - `openapi/openapi.yaml`
    - `openapi/openapi.json`
  - 完了条件:
    - `GET /v1/forecasts/{officeCode}/areas/{areaCode}` が定義されている。
    - 個別取得 schema に `office`、`publishingOffice`、`reportDatetime`、`weatherArea` が含まれている。
    - 個別取得 schema が temperature 系項目を含まないことが明記されている。
    - `404` の条件が `officeCode` 不正または `areaCode` 不一致として整理されている。

- [x] タスク1.3: 生成コードを再生成する
  - 目的:
    - OpenAPI と Go 実装側 interface を整合させる。
  - 対象ファイル:
    - `openapi/openapi.json`
    - `internal/gen/server.gen.go`
    - `internal/gen/types.gen.go`
    - `internal/gen/spec.gen.go`
  - 完了条件:
    - `oapi-codegen` 3 コマンドが成功する。
    - `openapi/openapi.yaml` と `openapi/openapi.json` が同期している。
    - 新規 query / path に対応した型と interface が生成される。

### [x] フェーズ2: `areas` 逆引き検索の実装

- [x] タスク2.1: `areas` 検索条件の受け口を handler / usecase に追加する
  - 目的:
    - 追加 query parameter をアプリ層へ渡せるようにする。
  - 対象ファイル:
    - `internal/handlers/areas.go`
    - `internal/usecases/list_areas.go`
    - `internal/handlers/server.go`
  - 完了条件:
    - `parent`、`name`、`officeName`、`child` をまとめて扱う filter struct が導入されている。
    - 既存 `GetV1Areas` の振る舞いが後方互換を保っている。

- [x] タスク2.2: `AreaMapper` に複合フィルタ処理を実装する
  - 目的:
    - `area.json` から取得した `Area` 一覧に AND 条件でフィルタを適用する。
  - 対象ファイル:
    - `internal/mappers/area_mapper.go`
    - 必要に応じて新規 filter 定義ファイル
  - 完了条件:
    - `parent`、`name`、`officeName`、`child` の各条件が完全一致で適用される。
    - 0 件時も空配列で返る。

- [x] タスク2.3: `areas` 検索のテストを追加する
  - 目的:
    - query 条件とレスポンス契約を固定する。
  - 対象ファイル:
    - `internal/mappers/area_mapper_test.go`
    - 必要に応じて `internal/handlers/areas_test.go`
  - 完了条件:
    - `name` 単独、`officeName` 単独、`child` 単独、複合条件、0 件のケースがテスト化される。
    - HTTP レベルで query parameter の受け渡しが確認できる。

### [x] フェーズ3: `weatherAreas` 個別取得の実装

- [x] タスク3.0: weather 個別取得用の upstream 取得方針を client に反映する
  - 目的:
    - `timeSeries[2]` 前提に引きずられず、weather 個別取得に必要な最小要件だけで upstream を扱えるようにする。
  - 対象ファイル:
    - `internal/clients/jma_client.go`
    - 必要に応じて client test
  - 完了条件:
    - weather 個別取得用の取得経路が `timeSeries[0]` 必須で成立する。
    - `timeSeries[2]` 欠落を理由に新規 endpoint が失敗しない設計方針がコード上で表現される。

- [x] タスク3.1: forecast 個別取得用の usecase / handler を追加する
  - 目的:
    - 新規 endpoint を実装する。
  - 対象ファイル:
    - `internal/clients/jma_client.go`
    - `internal/handlers/forecasts.go`
    - `internal/usecases/get_forecast.go`
    - `internal/handlers/server.go`
  - 完了条件:
    - `officeCode` と `areaCode` を受け取るメソッドが追加されている。
    - `officeCode` 不正時は既存と同様に `OFFICE_NOT_FOUND` 系の 404 を返す。
    - `areaCode` 不一致時は新規または既存エラーコードで 404 を返す。

- [x] タスク3.2: `ForecastMapper` に個別 area 抽出処理を追加する
  - 目的:
    - weather 個別取得用のデータから対象 `weatherArea` のみを返せるようにする。
  - 対象ファイル:
    - `internal/mappers/forecast_mapper.go`
    - 必要に応じて新規 mapper test
  - 完了条件:
    - 指定 `areaCode` の `weatherArea` を抽出できる。
    - 個別レスポンスに `publishingOffice`、`reportDatetime`、`office` が含まれる。
    - `temperatureAreas` を誤って混在させない。

- [x] タスク3.3: forecast 個別取得のテストを追加する
  - 目的:
    - 新規 endpoint のレスポンス契約とエラー挙動を固定する。
  - 対象ファイル:
    - 必要に応じて `internal/clients/jma_client_test.go`
    - `internal/mappers/forecast_mapper_test.go`
    - 必要に応じて `internal/handlers/forecasts_test.go`
    - 必要に応じて `internal/usecases/get_forecast_test.go`
  - 完了条件:
    - 正常系、`officeCode` 不正、`areaCode` 不一致がテスト化される。
    - `timeSeries[2]` 欠落でも新規 endpoint が成立するケース、またはその非対応方針がテストで固定される。
    - HTTP レベルで path parameter の配線が確認できる。

### [x] フェーズ4: 回帰確認とドキュメント整備

- [x] タスク4.1: 全体テストと生成物整合を確認する
  - 目的:
    - 仕様、生成コード、実装、テストの整合を最終確認する。
  - 対象ファイル:
    - `openapi/openapi.json`
    - `openapi/openapi.yaml`
    - `internal/gen/*`
    - `internal/**/*_test.go`
  - 完了条件:
    - `oapi-codegen` 再生成、`go test ./...`、必要に応じて `golangci-lint run` が成功する。
    - `openapi/openapi.yaml` と `openapi/openapi.json` の配信内容が一致している。

- [x] タスク4.2: 追加 API の利用導線を文書化する
  - 目的:
    - 新しい検索導線を README から把握できるようにする。
  - 対象ファイル:
    - `README.md`
    - `README-JA.md`
    - または関連設計ドキュメント
  - 完了条件:
    - 追加 endpoint と代表的な利用例が文書へ反映されている。

## 5. 影響範囲分析

### 5.1 変更対象

| 区分（直接影響/間接影響/テスト影響） | ファイル                                   | 変更種別       | 変更内容                                               | 関連タスク          |
| ------------------------------------ | ------------------------------------------ | -------------- | ------------------------------------------------------ | ------------------- |
| 直接影響                             | `openapi/openapi.json`                     | 修正           | YAML と同期した公開 spec 更新                          | タスク1.1, 1.2, 1.3 |
| 直接影響                             | `openapi/openapi.yaml`                     | 修正           | query / path / schema 追加                             | タスク1.1, 1.2      |
| 直接影響                             | `internal/clients/jma_client.go`           | 修正           | weather 個別取得向けの upstream 取得経路追加または緩和 | タスク3.0, 3.1      |
| 直接影響                             | `internal/gen/server.gen.go`               | 生成更新       | server interface と params の更新                      | タスク1.3           |
| 直接影響                             | `internal/gen/types.gen.go`                | 生成更新       | 新規 schema 型の追加                                   | タスク1.3           |
| 直接影響                             | `internal/gen/spec.gen.go`                 | 生成更新       | embedded spec 更新                                     | タスク1.3           |
| 直接影響                             | `internal/handlers/areas.go`               | 修正           | query filter の受け渡し更新                            | タスク2.1           |
| 直接影響                             | `internal/usecases/list_areas.go`          | 修正           | filter struct の導入と検索処理呼び出し                 | タスク2.1           |
| 直接影響                             | `internal/mappers/area_mapper.go`          | 修正           | `name` / `officeName` / `child` 条件のフィルタ追加     | タスク2.2           |
| 直接影響                             | `internal/handlers/forecasts.go`           | 修正           | 個別 area endpoint 追加                                | タスク3.1           |
| 直接影響                             | `internal/usecases/get_forecast.go`        | 修正           | 個別 `weatherArea` 取得メソッド追加                    | タスク3.1           |
| 直接影響                             | `internal/mappers/forecast_mapper.go`      | 修正           | 個別 area 抽出と専用レスポンス生成                     | タスク3.2           |
| 間接影響                             | `internal/handlers/server.go`              | 修正           | usecase 依存の受け口調整                               | タスク2.1, 3.1      |
| 間接影響                             | `cmd/server/main.go`                       | 修正可能性あり | 新規 interface / usecase 初期化の整合確認              | タスク3.1           |
| テスト影響                           | `internal/mappers/area_mapper_test.go`     | 修正           | 検索条件別の fixture test 追加                         | タスク2.3           |
| テスト影響                           | `internal/mappers/forecast_mapper_test.go` | 修正           | 個別 area 抽出の test 追加                             | タスク3.3           |
| テスト影響                           | `internal/clients/jma_client_test.go`      | 新規候補       | weather 個別取得向け upstream パースの test            | タスク3.3           |
| テスト影響                           | `internal/handlers/areas_test.go`          | 新規候補       | query parameter の HTTP テスト                         | タスク2.3           |
| テスト影響                           | `internal/handlers/forecasts_test.go`      | 新規候補       | path parameter の HTTP テスト                          | タスク3.3           |
| テスト影響                           | `internal/usecases/get_forecast_test.go`   | 新規候補       | `officeCode` / `areaCode` 異常系の usecase test        | タスク3.3           |
| 間接影響                             | `README.md`                                | 修正候補       | API 利用例更新                                         | タスク4.2           |
| 間接影響                             | `README-JA.md`                             | 修正候補       | API 利用例更新                                         | タスク4.2           |

### 5.2 依存関係

- 依存先:
  - JMA `area.json`
  - JMA forecast JSON
  - `oapi-codegen`
  - `chi` router と生成 server interface
- 影響先:
  - API 利用者の query / path 利用方法
  - OpenAPI UI (`/docs`) に表示される契約
  - handler と usecase の結合点

### 5.3 テスト影響

- 追加:
  - `areas` 検索条件別の unit test
  - `areas` query parameter の HTTP テスト
  - weather 個別取得向け client / usecase 異常系テスト
  - `weatherAreas` 個別取得の unit test
  - `weatherAreas` path parameter の HTTP テスト
- 更新:
  - `forecast_mapper_test.go` に専用レスポンスの期待値追加

## 6. 進捗管理

### 6.1 状態ボード

| タスク                                                               | 状態 | 担当                | 最終更新             | 次アクション | ブロッカー |
| -------------------------------------------------------------------- | ---- | ------------------- | -------------------- | ------------ | ---------- |
| タスク1.1 `areas` 検索条件を OpenAPI に追加する                      | DONE | operational-command | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク1.2 `weatherAreas` 個別取得 API の契約を追加する               | DONE | operational-command | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク1.3 生成コードを再生成する                                     | DONE | operational-command | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク2.1 `areas` 検索条件の受け口を handler / usecase に追加する    | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク2.2 `AreaMapper` に複合フィルタ処理を実装する                  | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク2.3 `areas` 検索のテストを追加する                             | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク3.0 weather 個別取得用の upstream 取得方針を client に反映する | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク3.1 forecast 個別取得用の usecase / handler を追加する         | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク3.2 `ForecastMapper` に個別 area 抽出処理を追加する            | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク3.3 forecast 個別取得のテストを追加する                        | DONE | worker              | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク4.1 全体テストと生成物整合を確認する                           | DONE | reviewer            | 2026-03-12 14:30 JST | 完了         | なし       |
| タスク4.2 追加 API の利用導線を文書化する                            | DONE | operational-command | 2026-03-12 14:30 JST | 完了         | なし       |

### 6.2 判断結果

- 2026-03-12 13:00 JST:
  - 実行基準は `docs/plans/api-extension-phase1-plan.md` とした。
  - 新規の進捗管理 Markdown は作成せず、本計画書内で状態を管理する。
  - 必要に応じて OpenAPI 存在確認以外の契約テスト追加
  - `openapi/openapi.json` の同期確認
- 2026-03-12 14:30 JST:
  - `GET /v1/areas` に `name` / `officeName` / `child` を追加し、完全一致の AND 条件と `200 + items: []` を OpenAPI と実装で一致させた。
  - `GET /v1/forecasts/{officeCode}/areas/{areaCode}` を追加し、`WEATHER_AREA_NOT_FOUND` を 404 として採用した。
  - weather 個別取得は `timeSeries[0]` 必須・`timeSeries[2]` 任意の専用取得経路を client に追加した。
  - `go test ./...` と `golangci-lint run` の成功を確認した。
- 2026-03-12 15:05 JST:
  - `GET /v1/forecasts/{officeCode}/areas/{areaCode}` は `weatherAreas.code` に加えて `temperatureAreas.code` でも検索可能に拡張した。
  - 個別レスポンスは `weatherArea` または `temperatureArea` の一方を返す契約へ更新した。
  - 404 エラーコードは `FORECAST_AREA_NOT_FOUND` に一般化した。
- 不要理由（不要な場合のみ）:
  - E2E は既存環境がないため必須外。ただし handler レベルの HTTP テストで代替する。

## 6. リスク評価

| リスク                                                                         | 影響度(1-3) | 発生確率(1-3) | スコア | 対策                                                                                          | 検証                                                      | 対応タスク               |
| ------------------------------------------------------------------------------ | ----------- | ------------- | ------ | --------------------------------------------------------------------------------------------- | --------------------------------------------------------- | ------------------------ |
| `GET /v1/areas` の複合検索仕様が曖昧で実装解釈が割れる                         | 3           | 2             | 6      | `200 + items: []`、AND 条件、完全一致を計画と OpenAPI に明記する                              | query 条件ごとの unit / HTTP テスト                       | タスク1.1, 2.3           |
| weather 個別取得が既存 client の `timeSeries[2]` 前提に引きずられて 502 化する | 3           | 2             | 6      | weather 個別取得専用の upstream 取得経路を用意し、最小要件を `timeSeries[0]` 必須へ切り分ける | client / usecase test で `timeSeries[2]` 欠落時の挙動確認 | タスク3.0, 3.3           |
| `weatherAreas` / `temperatureAreas` 個別取得で対象種別を誤って関連付ける       | 3           | 2             | 6      | 個別レスポンスを `weatherArea` または `temperatureArea` の片側返却に限定する                  | mapper test で返却種別を確認                              | タスク1.2, 3.2, 3.3      |
| 生成コード更新漏れで handler interface が不整合になる                          | 2           | 2             | 4      | OpenAPI 更新直後に `oapi-codegen` 再生成を固定手順にする                                      | `go test ./...` とビルド確認                              | タスク1.3, 4.1           |
| `openapi/openapi.yaml` と `openapi/openapi.json` が乖離する                    | 2           | 2             | 4      | 対象ファイルと手順に JSON 同期を含める                                                        | `/openapi.yaml` と `/openapi.json` の内容確認             | タスク1.1, 1.2, 1.3, 4.1 |
| query parameter 増加で `AreaUsecase.List` が保守しづらくなる                   | 2           | 2             | 4      | filter struct を導入して引数を整理する                                                        | usecase / handler テストで配線確認                        | タスク2.1                |
| JMA upstream のレスポンス構造変更で fixture と実装が乖離する                   | 2           | 2             | 4      | 既存 fixture を使ったテストに加え、必要時に fixture 更新を明示する                            | mapper test の失敗検知                                    | タスク2.3, 3.3           |

## 7. 検証計画

- 単体テスト:
  - `go test ./...`
  - `AreaMapper` の検索条件別テスト
  - `JMAClient` の weather 個別取得向けパーステスト
  - `ForecastMapper` の個別 area 抽出テスト
- 結合テスト:
  - `httptest` ベースで handler の query / path parameter 配線を確認
  - usecase の `officeCode` / `areaCode` 異常系テスト
- E2E/手動確認:
  - `/docs` または OpenAPI を参照し、追加 endpoint が表示されることを確認
  - `/openapi.yaml` と `/openapi.json` の両方で追加契約が配信されることを確認
  - ローカル起動時に代表ケースを手動確認
- 性能/セキュリティ確認:
  - 性能: 新規 DB や外部 I/O は増やさないため専用負荷試験は必須外
  - セキュリティ: 検索条件は文字列比較のみで、機微情報出力を増やさないことを確認

## 8. リリース・運用

- リリース手順:
  - `openapi/openapi.yaml` と `openapi/openapi.json` を更新
  - `oapi-codegen` 再生成
  - 実装更新
  - `go test ./...`
  - 必要に応じて `golangci-lint run`
  - README 更新
- ロールバック手順:
  - 追加 endpoint と query 拡張を含む変更一式を git で巻き戻す
  - OpenAPI と生成コードを同じコミット単位で戻す
- 監視/アラート:
  - 新規監視は追加しない
  - 既存ログで `AREA_NOT_FOUND` / `OFFICE_NOT_FOUND` / 新規 `FORECAST_AREA_NOT_FOUND` 系の頻度を確認対象とする

## 9. 未確定事項

- 現時点の未確定事項なし

## 10. セルフレビュー結果

| 課題                                                                      | 重要度 | 対応方針                                                                           |
| ------------------------------------------------------------------------- | ------ | ---------------------------------------------------------------------------------- |
| 旧計画書が `docs/designs` 配下で Planner テンプレートに準拠していなかった | Low    | `docs/plans` 配下へ Planner 準拠版を新規作成して反映済み                           |
| `weatherAreas` 個別取得レスポンスの必須項目が曖昧だった                   | Medium | `office`、`publishingOffice`、`reportDatetime`、`weatherArea` を返す方針で修正済み |
| `GET /v1/areas` の 0 件時挙動が曖昧だった                                 | Medium | `200 + items: []` に統一する前提へ修正済み                                         |
| テスト計画が unit test 中心で HTTP 契約確認が弱かった                     | Medium | handler レベルの HTTP テスト追加を計画へ反映済み                                   |
| weather 個別取得が既存 client の `timeSeries[2]` 前提に依存していた       | Medium | client 層の変更を計画に追加し、新規 endpoint の最小要件を明記済み                  |
| `openapi/openapi.json` の更新手順が漏れていた                             | Medium | 対象ファイル、検証、リリース手順へ反映済み                                         |

## 11. ユーザー確認事項（必要な場合のみ）

- 現時点で確認待ち事項なし
