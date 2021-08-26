package perms_manager

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

var sliceArgs = genBenchArg(64, 5, []float64{-2, 2, -1.5, 1.5}, 4, 15)

func removeNodesAppend(stack []string, needle []string) []string {
	check := func(v string) bool {
		for _, r := range needle {
			if v == r {
				return false
			}
		}
		return true
	}
	var ret []string
	for _, v := range stack {
		if check(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func BenchmarkExternal_removeNodesAppend(b *testing.B) {
	data := genBenchData(sliceArgs)
	b.ResetTimer()
	for _, c := range data {
		b.Run(c.args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				removeNodesAppend(c.haystack, c.needles)
			}
		})
	}
}

func removeNodesMakeCap(stack []string, needle []string) []string {
	check := func(v string) bool {
		for _, r := range needle {
			if v == r {
				return false
			}
		}
		return true
	}
	ret := make([]string, 0, len(stack))
	for _, v := range stack {
		if check(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func BenchmarkExternal_removeNodesMakeCap(b *testing.B) {
	data := genBenchData(sliceArgs)
	b.ResetTimer()
	for _, c := range data {
		b.Run(c.args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				removeNodesMakeCap(c.haystack, c.needles)
			}
		})
	}
}

func removeNodesDirtySwap(stack []string, needle []string) []string {
	validate := func(v string) bool {
		for _, r := range needle {
			if v == r {
				return false
			}
		}
		return true
	}
	f := 1
	for i := 0; i < len(stack)-f; {
		if validate(stack[i]) {
			i++
			continue
		}
		stack[len(stack)-f], stack[i] = stack[i], stack[len(stack)-f]
		f++
	}
	out := stack
	if f > 1 {
		out = stack[:len(stack)-f]
	}
	return out
}

func BenchmarkExternal_removeNodesDirtySwap(b *testing.B) {
	data := genBenchData(sliceArgs)
	b.ResetTimer()
	for _, c := range data {
		b.Run(c.args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				removeNodesDirtySwap(c.haystack, c.needles)
			}
		})
	}
}

func removeNodesDirtyReplace(stack []string, needle []string) []string {
	check := func(v string) bool {
		for _, r := range needle {
			if v == r {
				return false
			}
		}
		return true
	}
	i := 0
	for _, x := range stack {
		if check(x) {
			stack[i] = x
			i++
		}
	}
	stack = stack[:i]
	return stack
}

func BenchmarkExternal_removeNodesDirtyReplace(b *testing.B) {
	data := genBenchData(sliceArgs)
	b.ResetTimer()
	for _, c := range data {
		b.Run(c.args.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				removeNodesDirtyReplace(c.haystack, c.needles)
			}
		})
	}
}

type removeNodes func(stack []string, needle []string) []string

func getRemoveNodes() []removeNodes {
	return []removeNodes{removeNodesAppend, removeNodesMakeCap, removeNodesDirtySwap, removeNodesDirtyReplace}
}

func TestExternal_removeNodes(t *testing.T) {
	for _, f := range getRemoveNodes() {

		fn := strings.Split(getFunctionName(f), ".")[2]
		for _, tc := range getSliceTest() {
			t.Run(fmt.Sprintf("%s()%s", fn, tc.name), func(t *testing.T) {
				a := assert.New(t)
				st := make([]string, len(tc.in))
				copy(st, tc.in)
				r := f(st, tc.arg)

				a.ElementsMatch(tc.want, r)
				a.ElementsMatch(tc.in, st)
			})
		}
	}
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type sliceTest struct {
	name string
	in   []string
	arg  []string
	want []string
}

func getSliceTest() []sliceTest {
	return []sliceTest{
		{
			name: "simple",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"2", "3"},
			want: []string{"1", "4"},
		}, {
			name: "repeated",
			in:   []string{"1", "2", "2", "2", "4"},
			arg:  []string{"2"},
			want: []string{"1", "4"},
		}, {
			name: "multiple repeated",
			in:   []string{"1", "2", "3", "4", "4"},
			arg:  []string{"2", "3", "3", "4"},
			want: []string{"1"},
		}, {
			name: "extract center",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"2", "3"},
			want: []string{"1", "4"},
		}, {
			name: "remove all",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"2", "3", "4", "1"},
			want: []string(nil),
		}, {
			name: "nil arg",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string(nil),
			want: []string{"1", "2", "3", "4"},
		}, {
			name: "none arg",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{},
			want: []string{"1", "2", "3", "4"},
		}, {
			name: "impossible",
			in:   []string{"1", "2", "3", "4"},
			arg:  []string{"5", "6", "7", "8"},
			want: []string{"1", "2", "3", "4"},
		}, {
			name: "empty input",
			in:   []string{},
			arg:  []string{"1", "2"},
			want: []string(nil),
		}, {
			name: "nil input",
			in:   []string(nil),
			arg:  []string{"1", "2"},
			want: []string(nil),
		},
	}
}

func genBenchData(args []benchSliceArg) []benchSliceData {
	o := make([]benchSliceData, 0, len(args))
	for _, arg := range args {
		s := randSlice(arg.haystack, arg.txtLen)
		var n []string
		if arg.needles >= 0 {
			n = randNeedles(s, arg.needles)
		} else {
			n = randSlice(int(math.Abs(float64(arg.needles))), arg.txtLen+1)
		}

		o = append(o, benchSliceData{
			haystack: s,
			needles:  n,
			args:     arg,
		})
	}
	return o
}

type benchSliceData struct {
	haystack []string
	needles  []string
	args     benchSliceArg
}

func genBenchArg(startFac int, runTime int, needleFac []float64, growFac float64, txtLen int) []benchSliceArg {
	var data []benchSliceArg
	f := startFac
	for i := 0; i < runTime; i++ {
		for _, v := range needleFac {
			var fn int
			if v > 0 {
				fn = int(float64(f) / v)
			} else if v == 0 {
				fn = 0
			} else {
				fn = -int(float64(f) / math.Abs(v))
			}

			data = append(data, benchSliceArg{
				name:     fmt.Sprintf("(%d/%d)", f, fn),
				haystack: f,
				txtLen:   txtLen,
				needles:  fn,
			})
		}

		f = int(float64(f) * growFac)
	}
	return data
}

type benchSliceArg struct {
	name     string
	haystack int
	txtLen   int
	needles  int
}

func randSlice(n, len int) []string {
	r := make([]string, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, randStringBytes(len))
	}
	return r
}

func randNeedles(h []string, n int) []string {
	r := make([]string, 0, n)
	for i := 0; i < n; i++ {
		r = append(r, h[rand.Intn(len(h))])
	}
	return r
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
