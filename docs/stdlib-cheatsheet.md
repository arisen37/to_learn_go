# Standard-library cheat sheet: JavaScript and C++ to Go

## Strings and slices

| Familiar operation | Go | Important difference |
| --- | --- | --- |
| JS `s.slice(a,b)`, C++ substr | `s[a:b]` | Byte offsets; bounds panic; no negative indexes |
| Slice by Unicode code point | `string([]rune(s)[a:b])` | Allocates; code points are not grapheme clusters |
| JS array.slice | `slices.Clone(a[lo:hi])` | `a[lo:hi]` alone shares the backing array |
| JS push | `a = append(a, value)` | Keep returned slice; may allocate |
| JS pop | Read last, clear slot, then `a = a[:len(a)-1]` | Check nonempty; clear references when retaining backing array |
| JS splice/remove | `a = slices.Delete(a, lo, hi)` | Modifies shared backing array |
| JS includes/indexOf | `strings.Contains` / `Index`; `slices.Contains` / `Index` | Linear search; -1 means absent for Index |
| JS startsWith/endsWith | `strings.HasPrefix` / `HasSuffix` | Literal prefix/suffix |
| JS trim | `strings.TrimSpace` | Unicode whitespace |
| JS split/join | `strings.Split` / `Join` | Fields is often better for whitespace tokenization |
| Split once | `left,right,ok := strings.Cut(s,"=")` | Boolean distinguishes a missing delimiter |
| JS replaceAll | `strings.ReplaceAll` | Literal match; regexp is a separate package |
| JS toLowerCase | `strings.ToLower` | Unicode lowercase; not locale-specific full case folding |
| Case-insensitive comparison | `strings.EqualFold` | Unicode simple folding |
| C++ stoi / JS Number | `strconv.Atoi` / `ParseInt(s,10,64)` | Errors must be checked; no automatic coercion |
| Number to text | `strconv.Itoa`, `FormatInt`, `FormatFloat` | Explicit base/precision |
| String concatenation loop | `strings.Builder` | WriteString then String; do not copy a nonzero Builder |
| JS sort numeric | `slices.Sort(nums)` | Numeric sorting by element type |
| Custom sort | `slices.SortFunc(values, cmp)` | Comparator sign, not boolean Less |
| JS map/filter/reduce | Explicit range loop | Error handling and allocation are visible |
| Object key existence | `value,ok := m[key]` | Missing and stored zero are distinct through ok |
| Object keys | range map, append keys, then sort | Map order unspecified |
| Object clone | `maps.Clone(m)` | Shallow, like slice clone |

## JSON and HTTP

| JavaScript | Go | Boundary |
| --- | --- | --- |
| `JSON.parse(text)` | `json.Unmarshal([]byte(text), &v)` | Usually decode to a typed struct |
| `JSON.stringify(v)` | `data,err := json.Marshal(v)` | Encoding can fail; []byte is base64 |
| `await response.json()` | Check status, defer Body.Close, bounded `json.Decoder` | No `.json()` method; lifecycle is explicit |
| `fetch(url,{signal})` | `http.NewRequestWithContext` + reusable `client.Do` | Non-2xx is not necessarily an error |
| Property optional/missing | Pointer field or custom presence wrapper | Zero, absent and null require explicit design |
| Big integer JSON | `json.Number` via Decoder.UseNumber or int64 struct field | Unstructured default float64 can lose precision |
| Deferred JSON subtree | `json.RawMessage` | Validate/decode later |

```go
type Payload struct { Name string `json:"name"` }
var value Payload
if err := json.Unmarshal(data, &value); err != nil { return err }
encoded, err := json.Marshal(value)
_ = encoded
return err
```

For a Reader, use json.NewDecoder(r); for a Writer, use json.NewEncoder(w).
Encoder.Encode appends a newline. Decoder accepts one value at a time, so require a second
Decode to return io.EOF when your API requires exactly one value. For HTTP, bound the body
before or during decoding. DisallowUnknownFields catches unexpected struct keys but does
not reject duplicate keys. Domain validation remains your code's responsibility.

## CSV, files, time and errors

| Need | Standard API | Pitfall |
| --- | --- | --- |
| CSV records | `csv.NewReader(r).Read()` / `ReadAll()` | Quoted commas/newlines defeat Split |
| CSV output | `csv.NewWriter(w)`, Write, Flush, Error | Check errors after flushing |
| File bytes | `os.ReadFile` / `WriteFile` | Whole-file memory use; WriteFile is not a durability transaction |
| Stream file | `os.Open`, defer Close, bufio/io | Check Scanner.Err and buffered writer Flush |
| Paths | `filepath.Join` | URL paths use path, not filepath |
| Copy stream | `io.Copy(dst,src)` | Cancellation depends on underlying streams |
| Exact frame read | `io.ReadFull(r,b)` | Partial I/O is normal |
| Limit reads | `io.LimitReader`, `http.MaxBytesReader` | limit+1 helps detect oversized responses |
| URL encode | `url.Values.Encode`, `url.QueryEscape` | Do not concatenate unescaped user parameters |
| Time parse/format | `time.Parse(time.RFC3339,s)`, t.Format | Uses Go reference layout, not yyyy/mm tokens |
| Error wrap/inspect | `fmt.Errorf("context: %w",err)`, errors.Is/As | Comparing strings loses identity |
| Structured logs | `log/slog` | Log metadata, avoid secrets and unbounded bodies |
| Env presence | `os.LookupEnv` | Empty and absent may differ |
| Process | `exec.CommandContext(ctx,exe,args...)` | Avoid shell interpolation |

Architecture: `source bytes -> bounded Reader -> format parser -> typed validation -> domain -> encoder -> Writer`.
Runnable examples are in examples/parsing and examples/system; graded parsing tasks are easy 02/06/07,
medium 04 and hard 04/06. Official package docs are linked in [resources](resources.md).
