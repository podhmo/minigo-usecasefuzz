**今回の小さな再現例は、まず minigo 本体の決定的な回帰テストに置き、その種を difffuzz で広げるのがよいと思います。usecasefuzz には YAML などの実用シナリオを載せる。** この分担は、読んだ [Issue #162](https://github.com/podhmo/minigo/issues/162) の提案とも合っています。

特に、「difffuzz かユニットテストか」は完全な二択ではありません。現在も `testdata/difffuzz/<case>/main.go` と `want.stdout` を置くと、通常の `go test` から [`TestDiffRegressions`]($HOME/ghq/github.com/podhmo/minigo/difffuzz_test.go:21) が実行します。**差分テストで見つけた例を、普段の回帰テストとして固定する経路が既にあります。**

今回の例なら、私はこう分けます。

| 置き場所 | 主な役割 | 今回の例 |
|---|---|---|
| minigo 本体のスクリプト回帰テスト | Go の意味論を固定する | 型の同一性、addressability、コピーと共有、保持した Field の書き戻し、panic/recover |
| Go 側のユニット・統合テスト | 内部の不変条件や Engine API を確認する | `runtime.Copy`、フィールド探索、`NewSession` の設定継承 |
| minigo-usecasefuzz | ライブラリを組み合わせた実用動作を保証する | YAML のネスト・タグ・独自 marshal/unmarshal、TOML、template |
| difffuzz の生成器 | 固定例の周辺にある未知の組み合わせを探す | 型・値・操作順序を変えた reflect の比較 |

`NewSession` の設定継承は minigo 固有の API なので、Go 言語処理系との比較より、Engine を直接使うテストが自然です。一方、reflect の意味論を内部の `TypeDef` や mock Hooks だけで確認すると、コンパイラ・型解決・VM・ホスト変換を通したときの問題を取り逃がします。ここは小さな Go プログラムを実行する回帰テストを中心にしたいです。

difffuzz で対応できる範囲はかなりあります。ただ、**既存の runner を使うことと、生成器がその問題を発見できることは別**です。

現在の `corpus` は、今回の単一ファイルの再現例をそのまま Go と比較できます。例えば次の形です。

```sh
GOCACHE=/private/tmp/minigo-review-cache \
  go -C ./tools/difffuzz run ./ corpus \
  /private/tmp/minigo-reflect-mres/struct-replacement/main.go
```

一方、現在の生成ドメインは `text` と `num` です。reflect を探索するには生成対象を追加する必要があります。[生成器の構成]($HOME/ghq/github.com/podhmo/minigo/tools/difffuzz/gen.go:149) を見ると、段階的に進められます。

- **最初は固定テンプレート＋値の生成で十分**です。`IsZero`、整数幅、型比較、`CanSet` などは、既存の IIFE テンプレート方式で広げられます。
- **保持した handle と操作順序の問題には、操作列のモデルが必要**です。「Field を取得 → 親を置換 → Field に書く → 元の変数を観測」の順序を生成し、失敗に必要な順序を残して縮小します。
- 可変な初期値は probe 内で新しく作るべきです。probe 間で共有すると、実行順序や縮小によって結果が変わってしまいます。

ここで重要なのは、reflect の戻り値だけを見るのでなく、**元の変数と別名参照も観測すること**です。`SetInt` の結果だけ正しくても、別の格納先に書いている可能性があります。今回の Field 参照、Append、Bytes の問題はその典型です。

もう一つ、生成器より先に整えたいのが合否判定です。

現在の difffuzz は loud な `TRAP` を実装バックログとして扱い、`corpus` は不一致があってもレポートを出して正常終了します。`PENDING` の回帰例も、どんな失敗でも基本的には skip します。したがって、それだけでは「対応済み reflect が壊れたら CI を落とす」という保証になりません。[現在の判定コード]($HOME/ghq/github.com/podhmo/minigo/tools/difffuzz/main.go:424)

#162 の方針どおり、ケースごとに次を区別したいです。

- **対応済み**：Go と一致することを要求する。予期しない TRAP も失敗。
- **明示的に未対応**：指定した unsupported trap を要求する。
- **探索中**：差分を記録し、対応状況を管理する。

usecasefuzz 側も、期待する判定、異常時の終了コード、両側の timeout、対象 revision の確実なビルドを整えてから CI の保証に使うのがよさそうです。この問題も #162 に書かれています。

私なら、最初のマイルストーンをこうします。

1. 今回の42件・64例を不変条件ごとに整理し、まず10〜20個程度の代表的な意味論 seed を本体の回帰スイートへ入れる。残りも対応表で追跡する。
2. 修正したケースを厳密な必須テストに昇格させる。`PENDING` を消した後の `TestDiffRegressions` は既に実行エラーも不一致も失敗にできます。
3. YAML などの実用シナリオを少数追加し、usecasefuzz の CI 判定を固める。
4. seed を基に reflect の生成ドメインを作り、値の生成から操作列の生成・縮小へ広げる。

**最初から大きな reflect ファザーを完成させる必要はありません。今回の最小例を普段のテストで失わないようにするところから始めると、修正にも生成器の設計にも使えます。**
