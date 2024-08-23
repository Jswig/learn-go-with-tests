package iteration

import "fmt"
import "testing"

func TestRepeat(t *testing.T) {
	repeated := Repeat("a", 5)
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}

func BenchmarkRepeat(b *testing.B) {
	for range b.N {
		Repeat("a", 5)
	}
}

func ExampleRepeat() {
	result := Repeat("z", 3)
	fmt.Println(result)
	// Output: zzz
}
