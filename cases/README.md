# cases/

`func main()` standalone modules — each dir runs under both `go run`
(the oracle) and `minigo`, with stdout+stderr diffed by `../run.sh`.

## Verdicts

| verdict | meaning |
|---|---|
| PASS | outputs identical |
| DIFF | both ran, outputs differ |
| TRAP | minigo panicked/trapped where Go ran |
| REJECT | minigo failed to run (resolve/compile) where Go compiled |
| ACCEPT | minigo ran what `go run` rejects — worst case |
| PASS-REJECT | both refuse |

`lim-*` cases are deliberate limitation probes: they use packages that are
not bound into the interpreter and document how the failure presents.

Latest run (2026-10-12): **44 PASS / 1 DIFF / 1 ACCEPT / 0 REJECT / 0 TRAP** —
every `lim-*` probe that used to trap now runs to completion. `lim-xml`
passes since `reflect.TypeAssert[T]` (explicit type args on a bound
builtin) works, and `lim-http` passes now that `golang.org/x/net/...`
imports resolve inside the module graph. `lim-yaml` and `lim-toml` keep
passing — `lim-toml` uses `github.com/BurntSushi/toml` since pelletier's
`internal/danger` reinterprets `reflect.SliceHeader` via
`unsafe.Pointer` (policy: prefer deps that avoid unsafe/cgo; minigo has
`unsafe.Pointer` now but not `uintptr`-minted reinterpretation).
`lim-cgo` graduated from REJECT to DIFF: the `import "C"` files are
skipped so the program runs, but the `sqlite3` driver never registers
and `sql.Open` reports "unknown driver" — a graceful runtime failure
instead of a load-time refusal.
The one ACCEPT is `inspectuse`: intentional, since
`minigo.dev/inspect` exists only inside the interpreter and `go run`
cannot compile it.

## Cases and what they exercise

| case | the use case | surface used |
|---|---|---|
| wordfreq | word frequency count over text | strings.Fields/ToLower/TrimFunc, unicode.IsLetter/IsDigit, map, sort.Slice |
| logparse | access-log analysis (status rollup, top paths, slow requests) | regexp submatch/index, strings, strconv.ParseFloat, sort |
| markdown | mini Markdown→HTML transform (headings, lists, bold/code) | strings.Builder, regexp.ReplaceAllString, html.EscapeString |
| slugify | URL slug generation incl. accents | strings.NewReplacer/ToLower/Trim, regexp.ReplaceAllString, unicode |
| texttable | aligned report table | fmt width verbs (`%-Ns`, `%N.Nf`), strings.Repeat, local `type` in composite literal |
| mustache | hand-rolled `{{var}}` templating | regexp.ReplaceAllStringFunc/FindStringSubmatch, strings.Join |
| csvsum | TSV parse → group/sum → report | strings.Split/Cut, strconv.Atoi, sort |
| jsonconfig | load JSON config into typed struct, edit, MarshalIndent | encoding/json into struct & map[string]any, json tags |
| jsonlines | JSONL transform pipeline | json.Unmarshal per line into map, json.Marshal, json.Valid |
| jsonptr | safe nested map access (JSON-pointer style) | map[string]any walks, type assertions, strconv |
| reportpipe | JSON events → aggregate → formatted report | json.Unmarshal into []struct, strings.Builder + fmt.Fprintf |
| iniconf | INI parser → sorted flat keys | strings.Cut/TrimSpace, map, sort |
| urlparse | URL building/escaping | net/url QueryEscape/PathEscape/JoinPath/Unescape, strings.Cut |
| envconfig | env-var config with typed fallbacks | os.Getenv/Setenv/Environ, strconv |
| filewalk | directory inventory | filepath.WalkDir + os.DirEntry callback, os.ReadDir, filepath.Rel |
| filemerge | read files → merge → write → verify | os.ReadFile/WriteFile/MkdirTemp/Stat, encoding/hex |
| timecalc | date/duration arithmetic | time.Parse, Duration methods (Hours/String/Sub), ParseDuration |
| errwrap | error wrapping & inspection | fmt.Errorf %w, errors.Is/As/Join/Unwrap, custom error type |
| sortrec | record sorting several ways | sort.SliceStable/Search, slices.SortFunc/IsSorted/BinarySearch |
| generic | script-defined generics (Map/Filter/Sum/keys) | generics, ~int constraints, slices.IndexFunc |
| workerpool | goroutine fan-out/fan-in | chan, go, for-range ch, sync.WaitGroup, sort |
| strbuilder | efficient string assembly | strings.Builder, fmt.Fprintf, bytes.Buffer |
| sprintf | fmt verb coverage | %v/%+v/%#v/%T/%q/%x/%b/%e/width verbs on script values |
| regexprepl | find/replace/split by pattern | *regexp.Regexp methods (FindAll, ReplaceAll*, Split, FindIndex) |
| deferrecover | panic/recover + deferred cleanup | defer, recover, errors.New, fmt |
| fnvsum | FNV-1a hash in pure script | uint64 literals/arithmetic, encoding/hex |
| inspectuse | introspect own package (minigo-only, no oracle) | minigo.dev/inspect: Decls/Kind/Doc/Signature |
| reflectset | string-map → typed struct override merge | reflect Elem/Field/CanSet/SetString/SetInt/SetBool, conf: tags |
| reflectmap | map[string]any → struct mapper (mini codec) | reflect Field/Tag/Set/AssignableTo, `db:"-"` tags, type identity |
| reflectvalid | struct-tag request validation | reflect NumField/Tag/IsZero |
| reflectdeep | semantic equality + channel/slice driving | reflect.DeepEqual, ValueOf(chan) Send/Recv/Close, Copy |
| reflectdump | recursive struct → text dump renderer | reflect Field/Kind/Elem/Interface walks over nested structs/maps/ptrs |
| reflectcall | name → handler dispatch registry (router/job style) | reflect.Value.Call with built args, error-result handling |
| reflectclone | deep-copy a config graph via reflect | reflect New/Set/SetMapIndex/MakeSlice/MakeMap, nil handling |
| lim-flag | CLI flag parsing (limitation probe) | flag |
| lim-template | text/template rendering (probe) | text/template |
| lim-http | HTTP request building (probe) | net/http |
| lim-yaml | external module dependency (probe) | gopkg.in/yaml.v3 via go.mod |
| lim-xml | stdlib XML decode (probe) | encoding/xml — reflect.TypeAssert decode |
| lim-toml | external TOML decode (probe) | github.com/BurntSushi/toml |
| lim-bufio | line scanning (probe) | bufio.Scanner |
| lim-sha | SHA-256 checksum (probe) | crypto/sha256 |
| lim-csv | CSV reading (probe) | encoding/csv |
| lim-io | io.ReadAll (probe) | io |
| lim-cgo | cgo-dependent driver (probe) | github.com/mattn/go-sqlite3 — DIFF: runs, but the cgo driver never registers (`sql: unknown driver`) |

## Usage

```sh
./run.sh            # all cases (clones podhmo/minigo next to this repo on first run)
./run.sh wordfreq   # one case
MINIGO_DIR=/path/to/minigo ./run.sh   # run against a specific checkout/branch
diff out/wordfreq.want out/wordfreq.got
```
