package int128

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"
	"math/bits"
)

type Uint128 struct {
	lo uint64
	hi uint64
}

var (
	ZeroUint128 = Uint128{
		lo: 0,
		hi: 0,
	}

	bone      = big.NewInt(1)
	b2pow64   = big.NewInt(0).Lsh(bone, 64) // 2**64 - 1st number not fitting in uint64
	b2pow127  = big.NewInt(0).Lsh(bone, 127)
	b2pow127N = big.NewInt(0).Neg(b2pow127)
	b2pow128  = big.NewInt(0).Lsh(bone, 128) // 2**128 - 1st number not fitting in uint128
)

func NewUint128FromString(s string, base int) (r Uint128, ok bool) {
	var bi, rem big.Int
	_, ok = bi.SetString(s, base)
	if !ok {
		return
	}

	if bi.Sign() < 0 {
		ok = false
		return
	}
	if bi.Cmp(b2pow128) >= 0 {
		ok = false
		return
	}

	bi.QuoRem(&bi, b2pow64, &rem)

	r = Uint128{
		lo: rem.Uint64(),
		hi: bi.Uint64(),
	}

	return
}

func MustUint128FromString(s string) Uint128 {
	r, ok := NewUint128FromString(s, 10)
	if !ok {
		panic(s)
	}
	return r
}

func (lhs Uint128) BigInt() *big.Int {
	var blo, bhi big.Int
	bhi.SetUint64(lhs.hi)
	blo.SetUint64(lhs.lo)

	bhi.Mul(&bhi, b2pow64)
	bhi.Add(&bhi, &blo)

	return &bhi
}

func (lhs Uint128) String() string {
	return lhs.BigInt().String()
}

func (lhs Uint128) PutLittleEndian(buf []byte) {
	_ = buf[15]
	binary.LittleEndian.PutUint64(buf[0:8], lhs.lo)
	binary.LittleEndian.PutUint64(buf[8:16], lhs.hi)
}

func (lhs Uint128) PutBigEndian(buf []byte) {
	_ = buf[15]
	binary.BigEndian.PutUint64(buf[8:16], lhs.lo)
	binary.BigEndian.PutUint64(buf[0:8], lhs.hi)
}

func (lhs *Uint128) FromLittleEndian(buf []byte) {
	_ = buf[15]
	*lhs = Uint128{
		lo: binary.LittleEndian.Uint64(buf[0:8]),
		hi: binary.LittleEndian.Uint64(buf[8:16]),
	}
}

func (lhs *Uint128) FromBigEndian(buf []byte) {
	_ = buf[15]
	*lhs = Uint128{
		lo: binary.BigEndian.Uint64(buf[8:16]),
		hi: binary.BigEndian.Uint64(buf[0:8]),
	}
}

func (lhs *Uint128) UnmarshalJSON(p []byte) (err error) {
	var s json.Number
	err = json.Unmarshal(p, &s)
	if err != nil {
		return err
	}

	r, ok := NewUint128FromString(s.String(), 10)
	if !ok {
		return ErrUnmarshalError
	}
	*lhs = r
	return
}

func (lhs Uint128) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, lhs.String())), nil
}

func (lhs Uint128) IsNeg() bool {
	return false
}

func (lhs Uint128) IsZero() bool {
	return lhs == ZeroUint128
}

func (lhs Uint128) Neg() (neg Uint128) {
	neg, _ = SubCU128(ZeroUint128, lhs, 0)
	return
}

func (lhs Uint128) Add(rhs Uint128) Uint128 {
	return AddU128(lhs, rhs)
}

func (lhs *Uint128) AddAssign(rhs Uint128) *Uint128 {
	*lhs = AddU128(*lhs, rhs)
	return lhs
}

func (lhs Uint128) Sub(rhs Uint128) Uint128 {
	return SubU128(lhs, rhs)
}

func (lhs *Uint128) SubAssign(rhs Uint128) *Uint128 {
	*lhs = SubU128(*lhs, rhs)
	return lhs
}

func (lhs Uint128) Equal(rhs Uint128) bool {
	return lhs == rhs
}

func (lhs Uint128) Less(rhs Uint128) bool {
	/*if lhs.hi < rhs.hi {
		return true
	} else if lhs.hi > rhs.hi {
		return false
	} else {
		return lhs.lo < rhs.lo
	}*/
	_, c := SubCU128(lhs, rhs, 0)
	return (c != 0)
}

func CmpUint128(lhs, rhs Uint128) int {
	/*if lhs.hi < rhs.hi {
		return -1
	} else if lhs.hi > rhs.hi {
		return +1
	} else if lhs.lo < rhs.lo {
		return -1
	} else if lhs.lo > rhs.lo {
		return +1
	} else {
		return 0
	}*/
	if lhs.Equal(rhs) {
		return 0
	} else if lhs.Less(rhs) {
		return -1
	} else {
		return +1
	}
}

func LessUint128(lhs Uint128, rhs Uint128) bool {
	return lhs.Less(rhs)
}

func AddCU128(lhs Uint128, rhs Uint128, carry uint64) (sum Uint128, carryOut uint64) {
	lo, carry := bits.Add64(lhs.lo, rhs.lo, carry)
	hi, carry := bits.Add64(lhs.hi, rhs.hi, carry)
	return Uint128{
		lo: lo,
		hi: hi,
	}, carry
}

func AddU128(lhs Uint128, rhs Uint128) (sum Uint128) {
	sum, _ = AddCU128(lhs, rhs, 0)
	return
}

func SubCU128(lhs Uint128, rhs Uint128, borrow uint64) (diff Uint128, borrowOut uint64) {
	lo, borrow := bits.Sub64(lhs.lo, rhs.lo, borrow)
	hi, borrow := bits.Sub64(lhs.hi, rhs.hi, borrow)
	return Uint128{
		lo: lo,
		hi: hi,
	}, borrow
}

func SubU128(lhs Uint128, rhs Uint128) (diff Uint128) {
	diff, _ = SubCU128(lhs, rhs, 0)
	return
}
