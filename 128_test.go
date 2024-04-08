package int128

import (
	"encoding/json"
	"fmt"
	"strconv"
	"testing"
)

func TestParse(t *testing.T) {
	type TestData struct {
		S   string
		OkU bool
		OkI bool
	}

	tests := []TestData{
		{"0", true, true},
		{"1", true, true},
		{"-1", false, true},

		{"9223372036854775807", true, true},
		{"9223372036854775808", true, true},
		{"9223372036854775809", true, true},

		{"-9223372036854775807", false, true},
		{"-9223372036854775808", false, true},
		{"-9223372036854775809", false, true},

		{"170141183460469231731687303715884105727", true, true},
		{"170141183460469231731687303715884105728", true, false},
		{"-170141183460469231731687303715884105728", false, true},
		{"-170141183460469231731687303715884105729", false, false},

		{"340282366920938463463374607431768211455", true, false},
		{"340282366920938463463374607431768211456", false, false},
	}

	for _, test := range tests {
		// String -> Uint128
		{
			v, ok := NewUint128FromString(test.S, 10)
			if ok != test.OkU {
				t.Errorf("%s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
			}
			if test.OkU {
				vs := v.String()
				if test.S != vs {
					t.Errorf("%s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
				}
			}

			t.Logf("SuccessU: %s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
		}

		// String -> Int128
		{
			v, ok := NewInt128FromString(test.S, 10)
			if ok != test.OkI {
				t.Errorf("%s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
			}
			if test.OkI {
				vs := v.String()
				if test.S != vs {
					t.Errorf("%s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
				}
				vsNeg := (test.S[0] == '-')
				if vsNeg != v.IsNeg() {
					t.Fail()
				}
			}

			t.Logf("SuccessI: %s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
		}

		// JSON -> Uint128
		{
			quotedStr := strconv.Quote(string(test.S))
			var err error

			var v Uint128
			err = json.Unmarshal([]byte(test.S), &v)
			if (err == nil) != test.OkU {
				t.Errorf("%s -> %v: %s (%d : %d)", test.S, err, v.String(), v.lo, v.hi)
			}

			v = ZeroUint128
			err = json.Unmarshal([]byte(quotedStr), &v)
			if (err == nil) != test.OkU {
				t.Errorf("%s -> %v: %s (%d : %d)", test.S, err, v.String(), v.lo, v.hi)
			}

			if test.OkU {
				vs, err := json.Marshal(v)
				if err != nil {
					t.Errorf("%s -> %v", test.S, err)
				}
				if quotedStr != string(vs) {
					t.Errorf("%s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
				}
			}

			t.Logf("SuccessU: %s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
		}

		// JSON -> Int128
		{
			quotedStr := strconv.Quote(string(test.S))
			var err error

			var v Int128
			err = json.Unmarshal([]byte(test.S), &v)
			if (err == nil) != test.OkI {
				t.Errorf("%s -> %v: %s (%d : %d)", test.S, err, v.String(), v.lo, v.hi)
			}

			v = ZeroInt128
			err = json.Unmarshal([]byte(quotedStr), &v)
			if (err == nil) != test.OkI {
				t.Errorf("%s -> %v: %s (%d : %d)", test.S, err, v.String(), v.lo, v.hi)
			}

			if test.OkI {
				vs, err := json.Marshal(v)
				if err != nil {
					t.Errorf("%s -> %v", test.S, err)
				}
				if quotedStr != string(vs) {
					t.Errorf("%s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
				}
			}

			t.Logf("SuccessU: %s -> %s (%d : %d)", test.S, v.String(), v.lo, v.hi)
		}
	}
}

func TestU128Arith(t *testing.T) {
	type TestData struct {
		A string
		B string
		S string
		c uint64
	}

	tests := []TestData{
		{"0", "0", "0", 0},
		{"1", "340282366920938463463374607431768211455", "0", 1},
		{"170141183460469231731687303715884105727", "170141183460469231731687303715884105727", "340282366920938463463374607431768211454", 0},
		{"170141183460469231731687303715884105727", "170141183460469231731687303715884105728", "340282366920938463463374607431768211455", 0},
		{"170141183460469231731687303715884105727", "170141183460469231731687303715884105729", "0", 1},
	}

	for _, test := range tests {
		a := MustUint128FromString(test.A)
		b := MustUint128FromString(test.B)
		s := MustUint128FromString(test.S)

		sum, carry := AddCU128(a, b, 0)
		if s != sum || test.c != carry {
			t.Errorf("%v + %v != %v (%v)", a, b, sum, carry)
		}

		diff1, borrow1 := SubCU128(s, a, 0)
		if b != diff1 || test.c != borrow1 {
			t.Errorf("%v + %v != %v (%v)", a, b, sum, carry)
		}

		diff2, borrow2 := SubCU128(s, b, 0)
		if a != diff2 || test.c != borrow2 {
			t.Errorf("%v + %v != %v (%v)", a, b, sum, carry)
		}
	}
}

func TestI128Arith(t *testing.T) {
	type TestData struct {
		A string
		B string
		S string
		c uint64
	}

	tests := []TestData{
		{"0", "0", "0", 0},
		{"1", "170141183460469231731687303715884105727", "-170141183460469231731687303715884105728", 0},
		{"-1", "-1", "-2", 1},
		{"0", "-1", "-1", 0},
	}

	for _, test := range tests {
		a := MustInt128FromString(test.A)
		b := MustInt128FromString(test.B)
		s := MustInt128FromString(test.S)

		sum, carry := AddCI128(a, b, 0)
		if s != sum || test.c != carry {
			t.Errorf("%v + %v != %v (%v != %v)", a, b, sum, test.c, carry)
		}

		diff1, borrow1 := SubCI128(s, a, 0)
		if b != diff1 || test.c != borrow1 {
			t.Errorf("%v + %v != %v (%v)", a, b, sum, carry)
		}

		diff2, borrow2 := SubCI128(s, b, 0)
		if a != diff2 || test.c != borrow2 {
			t.Errorf("%v + %v != %v (%v)", a, b, sum, carry)
		}
	}
}

func TestCompareUint128(t *testing.T) {
	testCases := []struct {
		a        string
		b        string
		expected int
	}{
		{"0", "0", 0},
		{"1", "9223372036854775808", -1},
		{"170141183460469231731687303715884105727", "9223372036854775808", +1},
		{"340282366920938463463374607431768211455", "9223372036854775808", +1},
		{"340282366920938463463374607431768211455", "0", +1},
		{"0", "340282366920938463463374607431768211455", -1},
		{"1", "340282366920938463463374607431768211455", -1},
		{"340282366920938463463374607431768211455", "340282366920938463463374607431768211455", 0},
	}

	for _, tc := range testCases {
		fmt.Printf("\nRunning test: %s vs %s\n", tc.a, tc.b)
		var cmp int
		var less bool
		ia := MustUint128FromString(tc.a)
		ib := MustUint128FromString(tc.b)

		cmp = CmpUint128(ia, ib)
		if cmp != tc.expected {
			t.Errorf("Cmp(%s, %s) = %d != %d", tc.a, tc.b, cmp, tc.expected)
		}
		less = LessUint128(ia, ib)
		if less != (tc.expected == -1) {
			t.Errorf("Less(%s, %s) = %v != %d", tc.a, tc.b, less, tc.expected)
		}

		cmp = CmpUint128(ib, ia)
		if cmp != -tc.expected {
			t.Errorf("Cmp(%s, %s) = %d != %d", tc.b, tc.a, cmp, -tc.expected)
		}
		less = LessUint128(ib, ia)
		if less != (tc.expected == 1) {
			t.Errorf("Less(%s, %s) = %v != %d", tc.b, tc.a, less, -tc.expected)
		}
	}
}

func TestCompareInt128(t *testing.T) {
	testCases := []struct {
		a        string
		b        string
		expected int
	}{
		{"0", "0", 0},
		{"1", "-1", +1},
		{"-1", "0", -1},
		{"0", "-170141183460469231731687303715884105728", +1},
		{"-1", "-170141183460469231731687303715884105728", +1},
		{"-2", "-170141183460469231731687303715884105728", +1},
		{"170141183460469231731687303715884105727", "-170141183460469231731687303715884105728", +1},
		{"1", "-170141183460469231731687303715884105728", +1},
		{"170141183460469231731687303715884105727", "170141183460469231731687303715884105727", 0},
		{"-170141183460469231731687303715884105728", "-170141183460469231731687303715884105728", 0},
	}

	for _, tc := range testCases {
		fmt.Printf("\nRunning test: %s vs %s\n", tc.a, tc.b)
		var cmp int
		var less bool
		ia := MustInt128FromString(tc.a)
		ib := MustInt128FromString(tc.b)

		cmp = CmpInt128(ia, ib)
		if cmp != tc.expected {
			t.Errorf("Cmp(%s, %s) = %d != %d", tc.a, tc.b, cmp, tc.expected)
		}
		less = LessInt128(ia, ib)
		if less != (tc.expected == -1) {
			t.Errorf("Less(%s, %s) = %v != %d", tc.a, tc.b, less, tc.expected)
		}

		cmp = CmpInt128(ib, ia)
		if cmp != -tc.expected {
			t.Errorf("Cmp(%s, %s) = %d != %d", tc.b, tc.a, cmp, -tc.expected)
		}
		less = LessInt128(ib, ia)
		if less != (tc.expected == 1) {
			t.Errorf("Less(%s, %s) = %v != %d", tc.b, tc.a, less, -tc.expected)
		}
	}
}
