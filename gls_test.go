package gls

import "testing"

func TestSetGetGo(t *testing.T) {
	Set("V", 1)
	v := Get("V").(int)
	if v != 1 {
		t.Errorf("V not equal 1")
	}

	q := make(chan bool, 0)
	go func(){
		Set("V", 1)
		v := Get("V").(int)
		if v != 1 {
			t.Errorf("V not equal 1")
        }

		Go(func() {
			v := Get("V").(int)
			if v != 1 {
				t.Errorf("V not equal 1")
			}

			Set("V", 2)
			v = Get("V").(int)
			if v != 2 {
				t.Errorf("V not equal 2")
			}

			Set("X", &v)
			pv := Get("X").(*int)
			if *pv != 2 {
				t.Errorf("X not equal 2")
			}

			q <- true
        })
		<-q

		if Get("X") != nil {
			t.Errorf("X exists in parent.")
		}
		v = Get("V").(int)
		if v != 1 {
			t.Errorf("V not equal 1")
		}

		q <- true
    }()

	<-q
}

func Benchmark_Set(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Set("V", 1)
    }
}

func Benchmark_Get(b *testing.B) {
	Set("V", 1)
	for i := 0; i < b.N; i++ {
		Get("V")
    }
}

func Benchmark_GetNil(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Get("V")
    }
}

func Benchmark_Set_4Threads(b *testing.B) {
	q := make(chan bool, 0)
	for c := 0; c < 4; c++ {
		go func() {
			for i := 0; i < b.N; i++ {
				Set("V", 1)
			}
			q <- true
		}()
    }
	for c := 0; c < 4; c++ {
		<-q
	}
}

func Benchmark_Get_4Threads(b *testing.B) {
	q := make(chan bool, 0)
	for c := 0; c < 4; c++ {
		go func() {
			Set("V", 1)
			for i := 0; i < b.N; i++ {
				Get("V")
			}
			q <- true
		}()
    }
	for c := 0; c < 4; c++ {
		<-q
	}
}

func Benchmark_GetNil_4Threads(b *testing.B) {
	q := make(chan bool, 0)
	for c := 0; c < 4; c++ {
		go func() {
			for i := 0; i < b.N; i++ {
				Get("V")
			}
			q <- true
		}()
    }
	for c := 0; c < 4; c++ {
		<-q
	}
}

