package processing

import "testing"

func BenchmarkRun2(b *testing.B) {
	var n int = 10000
	for i := 0; i < b.N; i++ {
		states := Run(1234, int64(n), 1, 1)
		if len(states[0]) != n {
			println("Hello")
			b.Fail()
		}
	}
}
func BenchmarkRun(b *testing.B) {
	var n int = 10000
	for i := 0; i < b.N; i++ {
		states := RunNoCache(1234, int64(n), 1, 1)
		if len(states[0]) != n {
			println("Hello")
			b.Fail()
		}
	}
}

