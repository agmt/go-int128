package int128

import (
	"fmt"
	"math/big"
	"testing"
)

func BenchmarkString(b *testing.B) {
	var r Uint128
	var s string

	for _, s := range []string{"0", "170141183460469231731687303715884105727"} {
		b.Run(fmt.Sprintf("fromString(%s)", s), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				r, _ = NewUint128FromString(s, 10)
			}
		})
		b.Run(fmt.Sprintf("toString(%s)", s), func(b *testing.B) {
			for i := 0; i != b.N; i++ {
				s = r.String()
			}
		})
	}
	_ = r
	_ = s
}

func BenchmarkAdd(b *testing.B) {
	var x = 9223372036854543385
	var y = 9223372036854775574
	var z int

	b.Run("int64+int64", func(b *testing.B) {
		for i := 0; i != b.N; i++ {
			z = x + y
		}
	})
	_ = z
}

func BenchmarkBigAdd(b *testing.B) {
	var x = big.NewInt(9223372036854543385)
	var y = big.NewInt(9223372036854775574)
	var z big.Int

	b.Run("big.int+big.int", func(b *testing.B) {
		for i := 0; i != b.N; i++ {
			z.Add(x, y)
		}
	})
	_ = z
}

func BenchmarkUint128Add(b *testing.B) {
	var x = Uint128{9223372036854543385, 1}
	var y = Uint128{9223372036854775574, 2}
	var z Uint128

	b.Run("uint128+uint128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			z = AddU128(x, y)
		}
	})
	_ = z
}

func BenchmarkInt128Add(b *testing.B) {
	var x = Int128{9223372036854543385, 1}
	var y = Int128{9223372036854775574, 2}
	var z Int128

	b.Run("int128+iint128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			z = AddI128(x, y)
		}
	})
	_ = z
}
