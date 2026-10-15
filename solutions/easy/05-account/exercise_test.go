package lab

import (
	"math"
	"testing"
)

func TestAccount(t *testing.T) {
	if _, err := NewAccount(-1); err == nil {
		t.Fatal("negative starting balance")
	}
	a, err := NewAccount(100)
	if err != nil {
		t.Fatal(err)
	}
	if err = a.Deposit(50); err != nil {
		t.Fatal(err)
	}
	if err = a.Withdraw(25); err != nil || a.Balance() != 125 {
		t.Fatal(a.Balance(), err)
	}
	for _, n := range []int64{0, -1, 126} {
		before := a.Balance()
		if err = a.Withdraw(n); err == nil || a.Balance() != before {
			t.Fatal("invalid withdrawal changed balance")
		}
	}
	for _, n := range []int64{0, -1, math.MaxInt64} {
		before := a.Balance()
		if err = a.Deposit(n); err == nil || a.Balance() != before {
			t.Fatal("invalid deposit changed balance")
		}
	}
	if err = a.Withdraw(125); err != nil || a.Balance() != 0 {
		t.Fatal("withdraw entire balance")
	}
}
