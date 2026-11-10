package memory

import "testing"

func TestCorte2RequiresSeedPassword(t *testing.T) {
	if _, err := Corte2(""); err == nil {
		t.Fatal("empty IDENTITY_SEED_PASSWORD must fail")
	}
	got, err := Corte2("test-seed-only")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Fatal("seed store")
	}
}
