**前の5件を含め、修正対象を42件確認しました。** 各問題の独立した最小再現コード、Go の期待結果、minigo の実際の結果を保存しています。

- [再現コード一式・64例（ZIP）](/private/tmp/minigo-reflect-review-mres.zip)
- [全42件の一覧・コード・実行結果](/private/tmp/minigo-reflect-mres/review.txt)

例えば、既存動作の退行はこれだけで再現できます。

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	var a any = time.Unix(0, 0).UTC().Weekday()
	var b any = int(4)
	fmt.Println(a == b)
}
```

Go と比較元の minigo は `false`、今回のブランチは `true` です。

各再現例は次の形式で比較できます。コード内の `recover` は、panic の有無を比較するために入れています。

```sh
GOCACHE=/private/tmp/minigo-review-cache \
  go -C /private/tmp/minigo-reflect-mres/host-numeric-regression run ./

/private/tmp/minigo-yaml-review run \
  /private/tmp/minigo-reflect-mres/host-numeric-regression
```

前の5件にも独立した再現コードを用意しました。

| 指摘 | 最小再現コード | Go → 今回のブランチ |
|---|---|---|
| 1. Field の addressability | [コード](/private/tmp/minigo-reflect-mres/field-addressability/main.go) | `false false` → `true true` |
| 2. スライス型の同一性 | [コード](/private/tmp/minigo-reflect-mres/slice-type-identity/main.go) | `true` → `false` |
| 3. MakeSlice の容量 | [コード](/private/tmp/minigo-reflect-mres/make-slice-capacity/main.go) | `1 4` → `1 1` |
| 4. IsZero の判定 | [コード](/private/tmp/minigo-reflect-mres/is-zero/main.go) | `true false false` → `false true true` |
| 5. SetInt の整数幅 | [コード](/private/tmp/minigo-reflect-mres/set-int-width/main.go) | `1` → `257` |

追加の指摘は以下です。同じ原因で起きる派生ケースはまとめています。

::code-comment{title="[P2] 6. 型エイリアスを参照先と同じ型として扱う" body="keyOf がエイリアス自身の Name を使うため、TypeFor[A]() と TypeFor[int]() が type A = int でも別の型になります。byte と uint8 も同様です。interning 前にエイリアスを参照先へ解決してください。[MRE](/private/tmp/minigo-reflect-mres/alias-identity/main.go)：Go は true true、今回は false false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=88 end=92 priority=2}

::code-comment{title="[P2] 7. ローカル型とジェネリック型の識別情報を保持する" body="keyOf の package path と Name だけのキーでは、別々の関数内で宣言された同名型や S[int] と S[string] が衝突します。同じ RType が返り、後から調べた型の情報まで先に登録した型になります。宣言の同一性と型引数をキーに含めてください。[MRE](/private/tmp/minigo-reflect-mres/generic-identity/main.go)：Go は false、今回は true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=88 end=92 priority=2}

::code-comment{title="[P2] 8. ホスト整数のアンラップで型の同一性を失わない" body="binaryOp が equality より前にホストの名前付き整数を int64 に変えるため、any(time.Weekday(4)) == any(int(4)) が false から true に退行します。また reflect.Int+1 の型が reflect.Kind から int に変わります。interface equality では動的型を保持し、算術結果にも元の名前付き型を復元してください。[MRE](/private/tmp/minigo-reflect-mres/host-numeric-regression/main.go)。" file="$HOME/ghq/github.com/podhmo/minigo/vm/vm.go" start=5075 end=5079 priority=2}

::code-comment{title="[P2] 9. NewSession にパッケージモードを引き継ぐ" body="WithPackageModes の設定が NewSession の Engine 初期化にコピーされません。strings を ModeDeny にした親では import が拒否されますが、子セッションでは同じスクリプトが成功します。NewSession に pkgModes を引き継いでください。[ホスト側MRE](/private/tmp/minigo-reflect-mres/session-policy/main.go) は parent denied: true、session denied: false を出力します。" file="$HOME/ghq/github.com/podhmo/minigo/minigo.go" start=107 end=108 priority=2}

::code-comment{title="[P2] 10. Indirect はポインタを一段だけ解除する" body="Indirect がループで全ポインタを解除します。**int に適用すると Go は *int を返しますが、今回は int まで進みます。ポインタの場合に一度だけ Elem を呼んでください。[MRE](/private/tmp/minigo-reflect-mres/indirect-depth/main.go)：ptr → int。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=279 end=285 priority=2}

::code-comment{title="[P2] 11. ValueOf に渡した reflect.Value 自体を反射する" body="valueOfValue が入力の RValue をそのまま返すため、reflect.ValueOf(reflect.ValueOf(1)) が reflect.Value の構造体を表さず、中身の int を表します。TypeOf も同じ経路なので影響します。公開 API への入力と、内部の facade アンラップを区別してください。[MRE](/private/tmp/minigo-reflect-mres/value-of-value/main.go)：struct → int。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=246 end=249 priority=2}

::code-comment{title="[P2] 12. ValueOf の構造体引数をコピーする" body="ValueOf が構造体オブジェクトをそのまま保持するため、v := reflect.ValueOf(s) の後に s.X を変更すると v の値も変化します。Go では interface 引数への格納時に構造体がコピーされます。runtime.Copy または vc.Copy で引数をスナップショットしてください。[MRE](/private/tmp/minigo-reflect-mres/valueof-struct-copy/main.go)：1 → 9。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=255 end=255 priority=2}

::code-comment{title="[P2] 13. Elem・Addr・Slice に非公開フィールド由来の制約を引き継ぐ" body="Elem などの wrap 呼び出しで ro が失われます。非公開の *int フィールドを Elem で辿ると CanSet と CanInterface が true になり、外から書き換えられます。非公開スライスの Slice でも同様です。派生 Value に read-only provenance を保持してください。[MRE](/private/tmp/minigo-reflect-mres/private-pointer/main.go)：Go は false false と panic、今回は true true と 9。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=336 end=337 priority=2}

::code-comment{title="[P2] 14. Field の参照を構造体の格納先に結び付ける" body="FieldRef の Base が現在の Struct オブジェクトに固定されます。フィールド Value を取得した後、親 Value.Set で構造体全体を置換すると、フィールドへの書き込みが古いオブジェクトへ向かいます。親の addressable storage を保持する参照にしてください。[MRE](/private/tmp/minigo-reflect-mres/struct-replacement/main.go)：Go は 3、今回は 2。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=490 end=491 priority=2}

::code-comment{title="[P2] 15. Set と SetMapIndex で代入可能性を検証する" body="Set が入力と格納先の型を検証せず、int の変数に string を格納できます。SetMapIndex も map[string]int に string を格納して正常終了します。書き込む前に型の代入可能性と入力の read-only 制約を検証してください。[MRE](/private/tmp/minigo-reflect-mres/set-type-check/main.go)：Go は panic、今回は bad。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1030 end=1034 priority=2}

::code-comment{title="[P2] 16. 各 Value API が受け付ける Kind を検証する" body="script 側の SetString は int に文字列を書き込み、Uint は int を受け入れ、IsNil は int に false を返します。また Field はポインタを暗黙に解除します。Go ではいずれも panic する操作です。内部の値表現だけで判断せず、各 API の許可された Kind を検証してください。[MRE](/private/tmp/minigo-reflect-mres/set-kind-check/main.go)：Go は panic、今回は bad。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1091 end=1098 priority=2}

::code-comment{title="[P2] 17. reflect.Copy で構造体要素をコピーする" body="copy_ は runtime.Value の参照をコピーするため、[]Struct のコピー後も要素が共有されます。コピー元のフィールドを変更するとコピー先まで変わります。通常の copy builtin と同様に、重なりを考慮したスナップショットと要素ごとの値コピーを行ってください。[MRE](/private/tmp/minigo-reflect-mres/copy-struct/main.go)：1 → 9。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=540 end=544 priority=2}

::code-comment{title="[P2] 18. Append・AppendSlice の追加要素を値コピーする" body="append_ と appendSlice は構造体要素の参照をそのまま追加します。addressable な構造体 Value を追加してから元の構造体を変更すると、追加済みの要素も変わります。追加する構造体・配列要素に runtime.Copy を適用してください。[MRE](/private/tmp/minigo-reflect-mres/append-struct/main.go)：1 → 9。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=492 end=493 priority=2}

::code-comment{title="[P2] 19. Append で既存の空き容量を利用する" body="append_ と appendSlice は常に既存要素を新しい backing array にコピーします。容量に余裕がある場合も元のスライスとの共有が切れ、追加後の既存要素への変更が元に反映されません。通常の append と同じ容量・共有規則を使ってください。[MRE](/private/tmp/minigo-reflect-mres/append-backing/main.go)：Go は 9、今回は 0。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=492 end=492 priority=2}

::code-comment{title="[P2] 20. SetMapIndex で格納する構造体をコピーする" body="SetMapIndex が addressable な構造体 Value を参照のまま格納します。呼び出し後に元の構造体を変更すると map の格納済みの値まで変わります。型検証後にキーと値へ Go の値コピーを適用してください。[MRE](/private/tmp/minigo-reflect-mres/map-write-copy/main.go)：1 → 9。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=742 end=742 priority=2}

::code-comment{title="[P2] 21. Send・TrySend で送信する構造体をコピーする" body="Send が構造体の runtime.Value をそのまま channel に入れます。バッファ付き channel への送信後に元の構造体を変更すると、受信した値まで変わります。Send と TrySend の送信境界で値コピーしてください。[MRE](/private/tmp/minigo-reflect-mres/reflect-send-copy/main.go)：1 → 9。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1389 end=1390 priority=2}

::code-comment{title="[P2] 22. 戻り値のない Call は空の結果を返す" body="Call は Tuple 以外の結果を常に要素一つのスライスに包むため、func(){} の呼び出しでも len(result) が 1 になります。関数シグネチャの戻り値数を使い、戻り値なしと nil を一つ返す関数を区別してください。[MRE](/private/tmp/minigo-reflect-mres/call-void/main.go)：0 → 1。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1167 end=1167 priority=2}

::code-comment{title="[P2] 23. CallSlice で末尾のスライスを展開する" body="CallSlice がそのまま Call を呼ぶため、可変長引数の末尾スライスを一つの要素として渡します。func(...int) に []int{1,2} を渡す有効な呼び出しが panic します。CallSlice の可変長引数経路を実装してください。[MRE](/private/tmp/minigo-reflect-mres/call-slice/main.go)：Go は 3、今回は panic。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1171 end=1172 priority=2}

::code-comment{title="[P2] 24. Comparable は構造体・配列の内部も判定する" body="Comparable は最上位の Kind だけを見ているため、スライスを含む構造体や [1][]int に true を返します。これを信じて比較や map キー作成へ進むコードが誤動作します。構造体フィールドと配列要素の comparability を再帰判定してください。[MRE](/private/tmp/minigo-reflect-mres/comparable-struct/main.go)：false false → true true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=705 end=709 priority=2}

::code-comment{title="[P2] 25. AssignableTo に underlying type の代入規則を反映する" body="AssignableTo は型の完全一致と interface への代入しか認めません。type S []int の S から []int への代入は Go では可能ですが、false になります。片方が名前付き型でない場合の underlying type 一致など、Go の代入規則を判定してください。[MRE](/private/tmp/minigo-reflect-mres/assignable-underlying/main.go)：true → false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=670 end=676 priority=2}

::code-comment{title="[P2] 26. ConvertibleTo が不可能な変換を許可しないようにする" body="ConvertibleTo は整数から complex への変換や string から []int への変換にも true を返します。Kind の大きな範囲だけで判断せず、数値カテゴリとスライス要素型を確認してください。[MRE](/private/tmp/minigo-reflect-mres/convertible-invalid/main.go)：false false → true true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=688 end=695 priority=2}

::code-comment{title="[P2] 27. Equal に Go の型・comparability・値比較規則を適用する" body="Equal が構造比較へ直接進むため、異なる名前付き構造体の同じフィールド値を true、同じ名前付き整数同士を false と判定します。またスライス同士を比較しても panic しません。型一致と comparability を確認してから、両側の名前付き値を対称に処理してください。[MRE](/private/tmp/minigo-reflect-mres/equal-types/main.go)：false → true。派生例も一覧に収録しています。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1274 end=1277 priority=2}

::code-comment{title="[P2] 28. Bytes が返すスライスを元の値と共有させる" body="Bytes が常に新しい []byte を作るため、返されたスライスの変更が元の []byte に反映されません。Go の Bytes は backing array を共有する API です。script 値への共有 view を返す経路を用意してください。[MRE](/private/tmp/minigo-reflect-mres/bytes-alias/main.go)：Go は 9、今回は 1。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=926 end=931 priority=2}

::code-comment{title="[P2] 29. 合成した配列・map の型情報を値の生成後も保持する" body="ArrayOf の配列長と MapOf のキー型が RType にだけ保存され、Zero・New・MakeMap が使う TypeDef にありません。ArrayOf から Zero を作ると Kind が slice になり、MapOf から MakeMap を作ると Type().Key().Kind() が invalid になります。合成型の形状を runtime の型記述にも保持してください。[MRE](/private/tmp/minigo-reflect-mres/array-constructor/main.go)：array array true → array slice false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=404 end=407 priority=2}

::code-comment{title="[P2] 30. channel の方向を型情報に保持する" body="script 側の ChanOf は方向引数を無視します。また AST 由来の channel 型のキーも方向を区別しません。そのため受信専用と送信専用、双方向と受信専用の型が同一になります。TypeDef と canonical key に channel direction を保持してください。[MRE](/private/tmp/minigo-reflect-mres/chan-direction/main.go)：false → true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=389 end=389 priority=2}

::code-comment{title="[P2] 31. 関数の動的型にシグネチャを保持する" body="typeOfValue は関数を調べるたびに、シグネチャのない新しい TypeDef を作ります。そのため同じ関数を二度 TypeOf しても別の型になります。Function・Closure・BoundMethod・BuiltinFunc から実際のシグネチャを取得し、同じ関数型を intern してください。[MRE](/private/tmp/minigo-reflect-mres/function-identity/main.go)：true → false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=150 end=151 priority=2}

::code-comment{title="[P2] 32. Implements でメソッドのシグネチャを比較する" body="Implements はメソッド名だけを比較するため、F(int) を持つ型が F(string) を要求する interface を実装すると判定します。さらに interface でないターゲットにも true を返します。ターゲットの Kind を検証し、パラメータ・戻り値・非公開メソッドの所属も比較してください。[MRE](/private/tmp/minigo-reflect-mres/interface-signature/main.go)：false → true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=615 end=623 priority=2}

::code-comment{title="[P2] 33. reflect 用のメソッド集合に receiver と公開範囲を反映する" body="NumMethod などが engine の名前集合を直接使うため、値型 S に *S のポインタ receiver メソッドを含め、非公開メソッドも数えます。Implements の結果にも影響します。reflect の規則に沿って receiver と公開範囲を選別してください。[MRE](/private/tmp/minigo-reflect-mres/method-pointer-set/main.go)：0 false → 1 true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=561 end=565 priority=2}

::code-comment{title="[P2] 34. Method に正しい Type と PkgPath を設定する" body="script 型の Method と MethodByName は Type を設定せず、公開メソッドにも宣言パッケージの PkgPath を付けます。シグネチャ検査と公開メソッド判定が誤ります。Type を埋め、公開メソッドの PkgPath は空にしてください。[MRE](/private/tmp/minigo-reflect-mres/method-metadata/main.go)：false true → true false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=579 end=579 priority=2}

::code-comment{title="[P2] 35. Convert の string 変換に表示用 String を使わない" body="Convert が string の変換先で Value.String を使うため、[]byte{65,66} の変換結果が AB ではなく <[]uint8 Value> になります。整数から string への変換も同じ問題です。元の Kind に応じた Go の変換処理を実装してください。[MRE](/private/tmp/minigo-reflect-mres/convert-string/main.go)。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=1241 end=1242 priority=2}

::code-comment{title="[P2] 36. FieldByName で埋め込みとポインタ埋め込みを探索する" body="RValue.FieldByName は直接フィールドだけを探し、RType.FieldByName は埋め込み *Struct を探索しません。Go で見つかる昇格フィールドが IsValid=false または found=false になります。型側で深さ・曖昧性を考慮したフィールドパスを解決し、値側でもそのパスを辿ってください。[MRE](/private/tmp/minigo-reflect-mres/field-type-embedded-pointer/main.go)：Go は true と [0 0]、今回は false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=546 end=550 priority=2}

::code-comment{title="[P2] 37. 非公開の匿名フィールドにも PkgPath を設定する" body="Field が embedded の場合を除外しているため、非公開型 s を埋め込んだ匿名フィールドの PkgPath が空になります。空の PkgPath を公開フィールドの判定に使うコードが誤ります。匿名かどうかにかかわらず、非公開フィールドには宣言パッケージを設定してください。[MRE](/private/tmp/minigo-reflect-mres/private-struct-name/main.go)：false → true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rtype.go" start=503 end=505 priority=2}

::code-comment{title="[P2] 38. MapIter.Next で削除済みのキーを飛ばす" body="MapRange が開始時のキーを保存し、Next が現在の map を確認せず進むため、未訪問のキーを削除しても列挙されます。Go の map iteration では訪問前に削除されたキーは生成されません。Next でキーの存在を確認して飛ばしてください。[MRE](/private/tmp/minigo-reflect-mres/map-iter-delete/main.go)：false → true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=780 end=781 priority=2}

::code-comment{title="[P2] 39. 配列を Slice した結果にスライス型を設定する" body="Slice が元の Typ と td をそのまま引き継ぐため、配列を切り出した結果も array 型として報告されます。結果の要素型を保持した []T の型記述へ変更してください。[MRE](/private/tmp/minigo-reflect-mres/array-slice-type/main.go)：slice true → array false。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=656 end=657 priority=2}

::code-comment{title="[P2] 40. FieldRef の昇格探索で最も浅いフィールドを選ぶ" body="追加された Get の深さ優先探索は、先に現れた深い埋め込みのフィールドを、後に現れる浅いフィールドより優先します。同じ深さの曖昧性ではなく、有効な Go プログラムでも &s.X の読み取りが別フィールドになります。探索深度を比較して最も浅い一致を選んでください。[MRE](/private/tmp/minigo-reflect-mres/promoted-depth/main.go)：Go は 2 2、今回は 1 2。" file="$HOME/ghq/github.com/podhmo/minigo/runtime/value.go" start=227 end=229 priority=2}

::code-comment{title="[P3] 41. ArrayOf の負の長さを拒否する" body="script 側の ArrayOf は負の長さを検証せず、ArrayOf(-1,t).Len() が -1 を返します。Go では生成時に panic します。型を登録する前に長さを検証してください。[MRE](/private/tmp/minigo-reflect-mres/array-of-invalid/main.go)：panic → -1。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/minireflect.go" start=403 end=407 priority=3}

::code-comment{title="[P2] 42. 配列の Index・Slice に親の addressability を反映する" body="Index は配列でも常に IndexRef を作るため、ValueOf([1]int{1}).Index(0) が変更可能になります。Slice も addressable でない配列を受け入れます。スライス要素と配列要素を区別し、配列では親の addressability を引き継いでください。[MRE](/private/tmp/minigo-reflect-mres/array-index-address/main.go)：false false → true true。" file="$HOME/ghq/github.com/podhmo/minigo/minireflect/rvalue.go" start=634 end=638 priority=2}

`make test`・`make lint` は通過しています。usecasefuzz の37ケースに新しい退行はなく、`lim-yaml` は TRAP から PASS に改善しています。上記は追加の比較用コードで確認した問題です。リポジトリは変更していません。
