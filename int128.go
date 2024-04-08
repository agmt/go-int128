package int128

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/bits"
)

var (
	ErrUnmarshalError = errors.New("wrong i128 value")
)

type Int128 struct {
	lo uint64
	hi uint64
}

var (
	ZeroInt128 = Int128{
		lo: 0,
		hi: 0,
	}
)

func NewInt128FromString(s string, base int) (r Int128, ok bool) {
	var bi, rem big.Int

	_, ok = bi.SetString(s, base)
	if !ok {
		return
	}

	if bi.Cmp(b2pow127N) < 0 {
		ok = false
		return
	}
	if bi.Cmp(b2pow127) >= 0 {
		ok = false
		return
	}
	sign := bi.Sign()

	if sign < 0 {
		bi.Neg(&bi)
	}

	bi.QuoRem(&bi, b2pow64, &rem)

	r = Int128{
		lo: rem.Uint64(),
		hi: bi.Uint64(),
	}

	if sign < 0 {
		r = r.Neg()
	}

	return
}

func MustInt128FromString(s string) Int128 {
	r, ok := NewInt128FromString(s, 10)
	if !ok {
		panic(s)
	}
	return r
}

func (lhs Int128) BigInt() *big.Int {
	var blo, bhi big.Int
	bhi.SetUint64(lhs.hi)
	blo.SetUint64(lhs.lo)

	bhi.Mul(&bhi, b2pow64)
	bhi.Add(&bhi, &blo)

	if lhs.IsNeg() {
		bhi.Sub(b2pow128, &bhi)
		bhi.Neg(&bhi)
	}

	return &bhi
}

func (lhs Int128) String() string {
	return lhs.BigInt().String()
}

func (lhs Int128) PutLittleEndian(buf []byte) {
	_ = buf[15]
	binary.LittleEndian.PutUint64(buf[0:8], lhs.lo)
	binary.LittleEndian.PutUint64(buf[8:16], lhs.hi)
}

func (lhs Int128) PutBigEndian(buf []byte) {
	_ = buf[15]
	binary.BigEndian.PutUint64(buf[8:16], lhs.lo)
	binary.BigEndian.PutUint64(buf[0:8], lhs.hi)
}

func (lhs *Int128) FromLittleEndian(buf []byte) {
	_ = buf[15]
	*lhs = Int128{
		lo: binary.LittleEndian.Uint64(buf[0:8]),
		hi: binary.LittleEndian.Uint64(buf[8:16]),
	}
}

func (lhs *Int128) FromBigEndian(buf []byte) {
	_ = buf[15]
	*lhs = Int128{
		lo: binary.BigEndian.Uint64(buf[8:16]),
		hi: binary.BigEndian.Uint64(buf[0:8]),
	}
}

func (lhs *Int128) UnmarshalJSON(p []byte) (err error) {
	var s json.Number
	err = json.Unmarshal(p, &s)
	if err != nil {
		return err
	}

	r, ok := NewInt128FromString(s.String(), 10)
	if !ok {
		return ErrUnmarshalError
	}
	*lhs = r
	return
}

func (lhs Int128) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, lhs.String())), nil
}

func (lhs Int128) IsNeg() bool {
	return lhs.hi > math.MaxInt64
}

func (lhs Int128) IsZero() bool {
	return lhs == ZeroInt128
}

func (lhs Int128) Neg() (neg Int128) {
	neg, _ = SubCI128(ZeroInt128, lhs, 0)
	return
}

func (lhs Int128) Add(rhs Int128) Int128 {
	return AddI128(lhs, rhs)
}

func (lhs *Int128) AddAssign(rhs Int128) *Int128 {
	*lhs = AddI128(*lhs, rhs)
	return lhs
}

func (lhs Int128) Sub(rhs Int128) Int128 {
	return SubI128(lhs, rhs)
}

func (lhs *Int128) SubAssign(rhs Int128) *Int128 {
	*lhs = SubI128(*lhs, rhs)
	return lhs
}

func (lhs Int128) Equal(rhs Int128) bool {
	return lhs == rhs
}

func (lhs Int128) Less(rhs Int128) bool {
	r, _ := SubCI128(lhs, rhs, 0)
	SF := r.IsNeg()
	OF := (lhs.IsNeg() != rhs.IsNeg()) && (r.IsNeg() == rhs.IsNeg())
	return SF != OF
}

func CmpInt128(lhs, rhs Int128) int {
	if lhs.Equal(rhs) {
		return 0
	} else if lhs.Less(rhs) {
		return -1
	} else {
		return +1
	}
}

func LessInt128(lhs Int128, rhs Int128) bool {
	return lhs.Less(rhs)
}

func AddCI128(lhs Int128, rhs Int128, carry uint64) (sum Int128, carryOut uint64) {
	lo, carry := bits.Add64(lhs.lo, rhs.lo, carry)
	hi, carry := bits.Add64(lhs.hi, rhs.hi, carry)
	return Int128{
		lo: lo,
		hi: hi,
	}, carry
}

func AddI128(lhs Int128, rhs Int128) (sum Int128) {
	sum, _ = AddCI128(lhs, rhs, 0)
	return
}

func SubCI128(lhs Int128, rhs Int128, borrow uint64) (diff Int128, borrowOut uint64) {
	lo, borrow := bits.Sub64(lhs.lo, rhs.lo, borrow)
	hi, borrow := bits.Sub64(lhs.hi, rhs.hi, borrow)
	return Int128{
		lo: lo,
		hi: hi,
	}, borrow
}

func SubI128(lhs Int128, rhs Int128) (diff Int128) {
	diff, _ = SubCI128(lhs, rhs, 0)
	return
}
