package app

// All exercises the full builtin surface of the type mapping.
type All struct {
	B    bool          `json:"b"`
	I    int           `json:"i"`
	I8   int8          `json:"i8"`
	I16  int16         `json:"i16"`
	I32  int32         `json:"i32"`
	I64  int64         `json:"i64"`
	U    uint          `json:"u"`
	U8   uint8         `json:"u8"`
	U16  uint16        `json:"u16"`
	U32  uint32        `json:"u32"`
	U64  uint64        `json:"u64"`
	UP   uintptr       `json:"up"`
	F32  float32       `json:"f32"`
	F64  float64       `json:"f64"`
	C64  complex64     `json:"c64"`
	C128 complex128    `json:"c128"`
	S    string        `json:"s"`
	By   byte          `json:"by"`
	R    rune          `json:"r"`
	Any  any           `json:"any"`
}
