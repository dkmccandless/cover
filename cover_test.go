package cover

import (
	"reflect"
	"testing"

	"github.com/dkmccandless/bipartite"
)

// smap returns an sset populated with ss.
func smap(ss ...Subset) sset {
	m := make(sset, len(ss))
	for _, s := range ss {
		m[s] = struct{}{}
	}
	return m
}

// emap returns an eset populated with es.
func emap(es ...Element) eset {
	m := make(eset, len(es))
	for _, e := range es {
		m[e] = struct{}{}
	}
	return m
}

// input holds a Subset and Elements for constructing a Cover with Add.
type input struct {
	s  Subset
	es []Element
}

// fromInputs returns a bipartite.Graph populated with inputs.
func fromInputs(inputs ...input) *bipartite.Graph {
	g := bipartite.New()
	for _, in := range inputs {
		for _, e := range in.es {
			g.Add(in.s, e)
		}
	}
	return g
}

func TestAdd(t *testing.T) {
	for _, test := range []struct {
		ins  []input
		want *Cover
	}{
		// Subset containing no Elements
		{
			[]input{
				{"empty set", []Element{}},
			},
			New(),
		},
		// Subset containing one Element
		{
			[]input{
				{true, []Element{true}},
			},
			&Cover{
				in: fromInputs(
					input{true, []Element{true}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Subset containing many elements
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{1, 2, 4, 8}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Duplicated input of Subset with no new elements
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
				{"Powers of 2", []Element{}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{1, 2, 4, 8}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Repeated input
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
				{"Powers of 2", []Element{1, 2, 4, 8}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{1, 2, 4, 8}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Subsets containing the same Element
		{
			[]input{
				{"Powers of 2", []Element{2}},
				{"Even primes", []Element{2}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{2}},
					input{"Even primes", []Element{2}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Partial overlap
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
				{"Fibonacci numbers", []Element{0, 1, 2, 3, 5, 8}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{1, 2, 4, 8}},
					input{"Fibonacci numbers", []Element{0, 1, 2, 3, 5, 8}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Add empty Subset to populated Cover
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
				{"Fibonacci numbers", []Element{0, 1, 2, 3, 5, 8}},
				{"Odd perfect numbers", []Element{}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{1, 2, 4, 8}},
					input{"Fibonacci numbers", []Element{0, 1, 2, 3, 5, 8}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
		// Add additional Elements of a Subset
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
				{"Fibonacci numbers", []Element{0, 1, 2, 3, 5, 8}},
				{"Fibonacci numbers", []Element{13}},
			},
			&Cover{
				in: fromInputs(
					input{"Powers of 2", []Element{1, 2, 4, 8}},
					input{"Fibonacci numbers", []Element{0, 1, 2, 3, 5, 8, 13}},
				),
				m: bipartite.New(),

				essential: smap(),
			},
		},
	} {
		c := New()
		for _, in := range test.ins {
			c.Add(in.s, in.es...)
		}
		if !reflect.DeepEqual(c, test.want) {
			t.Errorf("Add(%+v): got %+v, want %+v", test.ins, c, test.want)
		}
	}
}

func TestDominates(t *testing.T) {
	for _, test := range []struct {
		c   *Cover
		dom map[Subset]sset
	}{
		{
			&Cover{
				in: fromInputs(input{true, []Element{true}}),
				m:  fromInputs(input{true, []Element{true}}),

				essential: smap(),
			},
			map[Subset]sset{},
		},
		{
			&Cover{
				in: fromInputs(
					input{"A", []Element{"x"}},
					input{"B", []Element{"y"}},
				),
				m: fromInputs(
					input{"A", []Element{"x"}},
					input{"B", []Element{"y"}},
				),
				essential: smap(),
			},
			map[Subset]sset{},
		},
		{
			&Cover{
				in: fromInputs(
					input{"A", []Element{"x"}},
					input{"B", []Element{"x", "y", "z"}},
				),
				m: fromInputs(
					input{"A", []Element{"x"}},
					input{"B", []Element{"x", "y", "z"}},
				),
				essential: smap(),
			},
			map[Subset]sset{"B": smap("A")},
		},
		{
			&Cover{
				in: fromInputs(
					input{"A", []Element{2}},
					input{"B", []Element{2, 6}},
					input{"C", []Element{2, 6}},
					input{"D", []Element{1, 2, 4}},
					input{"E", []Element{3, 5, 7}},
					input{"F", []Element{0, 1, 2, 4, 7}},
				),
				m: fromInputs(
					input{"A", []Element{2}},
					input{"B", []Element{2, 6}},
					input{"C", []Element{2, 6}},
					input{"D", []Element{1, 2, 4}},
					input{"E", []Element{3, 5, 7}},
					input{"F", []Element{0, 1, 2, 4, 7}},
				),
				essential: smap(),
			},
			map[Subset]sset{
				"B": smap("A"),
				"C": smap("A"),
				"D": smap("A"),
				"F": smap("A", "D"),
			},
		},
	} {
		for a := range test.c.m.As() {
			for b := range test.c.m.As() {
				_, want := test.dom[a][b]
				if got := test.c.dominates(a, b); got != want {
					t.Errorf("dominates(%+v, %v, %v): got %v, want %v", test.c, a, b, got, want)
				}
			}
		}
	}
}

// copy copies the information in c into a new Cover and returns a pointer to it.
// The returned cover is deeply equal to c but shares no memory with it.
func (c *Cover) copy() *Cover {
	return &Cover{
		in: bipartite.Copy(c.in),
		m:  bipartite.Copy(c.m),

		essential: c.essential.copy(),
	}
}

func TestCopy(t *testing.T) {
	for name, test := range coverTests {
		for _, c := range []*Cover{test.c, test.s, test.e, test.sim} {
			if got := c.copy(); !reflect.DeepEqual(c, got) {
				t.Errorf("copy(%v, %#v): got %#v", name, c, got)
			}
		}
	}
}

func TestReduceS(t *testing.T) {
	for name, test := range coverTests {
		got := test.c.copy()
		if gotok := got.reduceS(); gotok != test.sok || !reflect.DeepEqual(got, test.s) {
			t.Errorf("reduceS(%v): got %+v, %v; want %+v, %v", name, got, gotok, test.s, test.sok)
		}
	}
}

func TestReduceE(t *testing.T) {
	for name, test := range coverTests {
		got := test.c.copy()
		if gotok := got.reduceE(); gotok != test.eok || !reflect.DeepEqual(got, test.e) {
			t.Errorf("reduceE(%v): got %+v, %v; want %+v, %v", name, got, gotok, test.e, test.eok)
		}
	}
}

func TestSimplify(t *testing.T) {
	for name, test := range coverTests {
		got := test.c.copy()
		if gotok := got.simplify(); gotok != test.simok || !reflect.DeepEqual(got, test.sim) {
			t.Errorf("simplify(%v): got %+v, %v; want %+v, %v", name, got, gotok, test.sim, test.simok)
		}
	}
}

func TestMinimize(t *testing.T) {
	for name, test := range coverTests {
		c := test.c.copy()
		// got and test.want must have identical contents, possibly in different orders.
		if got := c.Minimize(); len(got) != len(test.min) || !allMatch(got, test.min) {
			t.Errorf("Minimize(%v): got %v, want %v", name, got, test.min)
		}
	}
}

// allMatch reports whether a and b contain the same elements up to ordering.
func allMatch(a, b [][]Subset) bool {
	bms := make([]sset, len(b))
	for i, bs := range b {
		bms[i] = smap(bs...)
	}
	for _, as := range a {
		am := smap(as...)
		var ok bool
		for j, bm := range bms {
			if reflect.DeepEqual(am, bm) {
				// Match each element of b only once
				bms[j], bms = bms[len(bms)-1], bms[:len(bms)-1]
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

var coverTests = map[string]struct {
	// The input Cover. Do not mutate: use copy() and call methods on the copy.
	c *Cover

	// 	Cover after reduceS, reduceE, and simplify
	s, e, sim *Cover

	// Boolean output of reduceS, reduceE, and simplify
	sok, eok, simok bool

	// Output of Minimize
	min [][]Subset
}{
	"empty set": {
		c: New(),
		s: New(), sok: false,
		e: New(), eok: false,
		sim: New(), simok: true,
		min: [][]Subset{{}},
	},
	"tautology": {
		c: &Cover{
			in: fromInputs(input{true, []Element{true}}),
			m:  fromInputs(input{true, []Element{true}}),

			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(input{true, []Element{true}}),
			m:  fromInputs(input{true, []Element{true}}),

			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(input{true, []Element{true}}),
			m:  bipartite.New(),

			essential: smap(true),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(input{true, []Element{true}}),
			m:  bipartite.New(),

			essential: smap(true),
		},
		simok: true,
		min:   [][]Subset{{true}},
	},
	"disjoint A and B": {
		c: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			m: bipartite.New(),

			essential: smap("A", "B"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			m: bipartite.New(),

			essential: smap("A", "B"),
		},
		simok: true,
		min:   [][]Subset{{"A", "B"}},
	},
	"1 Subset contains 2 Elements": {
		c: &Cover{
			in: fromInputs(input{"A", []Element{"x", "y"}}),
			m:  fromInputs(input{"A", []Element{"x", "y"}}),

			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(input{"A", []Element{"x", "y"}}),
			m:  fromInputs(input{"A", []Element{"x", "y"}}),

			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(input{"A", []Element{"x", "y"}}),
			m:  bipartite.New(),

			essential: smap("A"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(input{"A", []Element{"x", "y"}}),
			m:  bipartite.New(),

			essential: smap("A"),
		},
		simok: true,
		min:   [][]Subset{{"A"}},
	},
	"2 Subsets contain 1 Element": {
		c: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			essential: smap(),
		},
		eok: false,
		sim: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x"}},
			),
			essential: smap(),
		},
		simok: false,
		min:   [][]Subset{{"A"}, {"B"}},
	},
	"B contains A": {
		c: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x", "y", "z"}},
			),
			m: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x", "y", "z"}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x", "y", "z"}},
			),
			m: fromInputs(
				input{"B", []Element{"x", "y", "z"}},
			),
			essential: smap(),
		},
		sok: true,
		e: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x", "y", "z"}},
			),
			m: bipartite.New(),

			essential: smap("B"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x", "y", "z"}},
			),
			m: bipartite.New(),

			essential: smap("B"),
		},
		simok: true,
		min:   [][]Subset{{"B"}},
	},
	"seven-segment A": {
		c: &Cover{
			in: fromInputs(
				input{"0-1-", []Element{2, 3, 6, 7}},
				input{"01-1", []Element{5, 7}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-11-", []Element{6, 7, 14, 15}},
				input{"100-", []Element{8, 9}},
				input{"1--0", []Element{8, 10, 12, 14}},
				input{"11-0", []Element{12, 14}},
			),
			m: fromInputs(
				input{"0-1-", []Element{2, 3, 6, 7}},
				input{"01-1", []Element{5, 7}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-11-", []Element{6, 7, 14, 15}},
				input{"100-", []Element{8, 9}},
				input{"1--0", []Element{8, 10, 12, 14}},
				input{"11-0", []Element{12, 14}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"0-1-", []Element{2, 3, 6, 7}},
				input{"01-1", []Element{5, 7}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-11-", []Element{6, 7, 14, 15}},
				input{"100-", []Element{8, 9}},
				input{"1--0", []Element{8, 10, 12, 14}},
				input{"11-0", []Element{12, 14}},
			),
			m: fromInputs(
				input{"0-1-", []Element{2, 3, 6, 7}},
				input{"01-1", []Element{5, 7}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-11-", []Element{6, 7, 14, 15}},
				input{"100-", []Element{8, 9}},
				input{"1--0", []Element{8, 10, 12, 14}},
			),
			essential: smap(),
		},
		sok: true,
		e: &Cover{
			in: fromInputs(
				input{"0-1-", []Element{2, 3, 6, 7}},
				input{"01-1", []Element{5, 7}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-11-", []Element{6, 7, 14, 15}},
				input{"100-", []Element{8, 9}},
				input{"1--0", []Element{8, 10, 12, 14}},
				input{"11-0", []Element{12, 14}},
			),
			m: fromInputs(
				input{"1--0", []Element{12}},
				input{"11-0", []Element{12}},
			),
			essential: smap("0-1-", "01-1", "-0-0", "-11-", "100-"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"0-1-", []Element{2, 3, 6, 7}},
				input{"01-1", []Element{5, 7}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-11-", []Element{6, 7, 14, 15}},
				input{"100-", []Element{8, 9}},
				input{"1--0", []Element{8, 10, 12, 14}},
				input{"11-0", []Element{12, 14}},
			),
			m: bipartite.New(),

			essential: smap("0-1-", "01-1", "-0-0", "-11-", "100-", "1--0"),
		},
		simok: true,
		min:   [][]Subset{{"0-1-", "01-1", "-0-0", "-11-", "100-", "1--0"}},
	},
	"seven-segment B": {
		c: &Cover{
			in: fromInputs(
				input{"00--", []Element{0, 1, 2, 3}},
				input{"0-00", []Element{0, 4}},
				input{"0-11", []Element{3, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"00--", []Element{0, 1, 2, 3}},
				input{"0-00", []Element{0, 4}},
				input{"0-11", []Element{3, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"1-01", []Element{9, 13}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"00--", []Element{0, 1, 2, 3}},
				input{"0-00", []Element{0, 4}},
				input{"0-11", []Element{3, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"00--", []Element{0, 1, 2, 3}},
				input{"0-00", []Element{0, 4}},
				input{"0-11", []Element{3, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"1-01", []Element{9, 13}},
			),
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(
				input{"00--", []Element{0, 1, 2, 3}},
				input{"0-00", []Element{0, 4}},
				input{"0-11", []Element{3, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"00--", []Element{1}},
				input{"-00-", []Element{1}},
			),
			essential: smap("0-00", "0-11", "-0-0", "1-01"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"00--", []Element{0, 1, 2, 3}},
				input{"0-00", []Element{0, 4}},
				input{"0-11", []Element{3, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-0", []Element{0, 2, 8, 10}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"00--", []Element{1}},
				input{"-00-", []Element{1}},
			),
			essential: smap("0-00", "0-11", "-0-0", "1-01"),
		},
		simok: false,
		min: [][]Subset{
			{"0-00", "0-11", "-0-0", "1-01", "00--"},
			{"0-00", "0-11", "-0-0", "1-01", "-00-"},
		},
	},
	"seven-segment C": {
		c: &Cover{
			in: fromInputs(
				input{"0-0-", []Element{0, 1, 4, 5}},
				input{"0--1", []Element{1, 3, 5, 7}},
				input{"01--", []Element{4, 5, 6, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-1", []Element{1, 3, 9, 11}},
				input{"--01", []Element{1, 5, 9, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
			),
			m: fromInputs(
				input{"0-0-", []Element{0, 1, 4, 5}},
				input{"0--1", []Element{1, 3, 5, 7}},
				input{"01--", []Element{4, 5, 6, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-1", []Element{1, 3, 9, 11}},
				input{"--01", []Element{1, 5, 9, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"0-0-", []Element{0, 1, 4, 5}},
				input{"0--1", []Element{1, 3, 5, 7}},
				input{"01--", []Element{4, 5, 6, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-1", []Element{1, 3, 9, 11}},
				input{"--01", []Element{1, 5, 9, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
			),
			m: fromInputs(
				input{"0-0-", []Element{0, 1, 4, 5}},
				input{"0--1", []Element{1, 3, 5, 7}},
				input{"01--", []Element{4, 5, 6, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-1", []Element{1, 3, 9, 11}},
				input{"--01", []Element{1, 5, 9, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
			),
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(
				input{"0-0-", []Element{0, 1, 4, 5}},
				input{"0--1", []Element{1, 3, 5, 7}},
				input{"01--", []Element{4, 5, 6, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-1", []Element{1, 3, 9, 11}},
				input{"--01", []Element{1, 5, 9, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
			),
			m: fromInputs(
				input{"0-0-", []Element{0}},
				input{"0--1", []Element{3}},
				input{"-00-", []Element{0}},
				input{"-0-1", []Element{3}},
			),
			essential: smap("01--", "--01", "10--"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"0-0-", []Element{0, 1, 4, 5}},
				input{"0--1", []Element{1, 3, 5, 7}},
				input{"01--", []Element{4, 5, 6, 7}},
				input{"-00-", []Element{0, 1, 8, 9}},
				input{"-0-1", []Element{1, 3, 9, 11}},
				input{"--01", []Element{1, 5, 9, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
			),
			m: fromInputs(
				input{"0-0-", []Element{0}},
				input{"0--1", []Element{3}},
				input{"-00-", []Element{0}},
				input{"-0-1", []Element{3}},
			),
			essential: smap("01--", "--01", "10--"),
		},
		simok: false,
		min: [][]Subset{
			{"01--", "--01", "10--", "0-0-", "0--1"},
			{"01--", "--01", "10--", "0-0-", "-0-1"},
			{"01--", "--01", "10--", "-00-", "0--1"},
			{"01--", "--01", "10--", "-00-", "-0-1"},
		},
	},
	"seven-segment D": {
		c: &Cover{
			in: fromInputs(
				input{"001-", []Element{2, 3}},
				input{"00-0", []Element{0, 2}},
				input{"0-10", []Element{2, 6}},
				input{"-000", []Element{0, 8}},
				input{"-011", []Element{3, 11}},
				input{"-101", []Element{5, 13}},
				input{"-110", []Element{6, 14}},
				input{"10-1", []Element{9, 11}},
				input{"1-0-", []Element{8, 9, 12, 13}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"001-", []Element{2, 3}},
				input{"00-0", []Element{0, 2}},
				input{"0-10", []Element{2, 6}},
				input{"-000", []Element{0, 8}},
				input{"-011", []Element{3, 11}},
				input{"-101", []Element{5, 13}},
				input{"-110", []Element{6, 14}},
				input{"10-1", []Element{9, 11}},
				input{"1-0-", []Element{8, 9, 12, 13}},
				input{"1-01", []Element{9, 13}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"00-0", []Element{0, 2}},
				input{"001-", []Element{2, 3}},
				input{"0-10", []Element{2, 6}},
				input{"-000", []Element{0, 8}},
				input{"-011", []Element{3, 11}},
				input{"-101", []Element{5, 13}},
				input{"-110", []Element{6, 14}},
				input{"10-1", []Element{9, 11}},
				input{"1-0-", []Element{8, 9, 12, 13}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"00-0", []Element{0, 2}},
				input{"001-", []Element{2, 3}},
				input{"0-10", []Element{2, 6}},
				input{"-000", []Element{0, 8}},
				input{"-011", []Element{3, 11}},
				input{"-101", []Element{5, 13}},
				input{"-110", []Element{6, 14}},
				input{"10-1", []Element{9, 11}},
				input{"1-0-", []Element{8, 9, 12, 13}},
			),
			essential: smap(),
		},
		sok: true,
		e: &Cover{
			in: fromInputs(
				input{"00-0", []Element{0, 2}},
				input{"001-", []Element{2, 3}},
				input{"0-10", []Element{2, 6}},
				input{"-000", []Element{0, 8}},
				input{"-011", []Element{3, 11}},
				input{"-101", []Element{5, 13}},
				input{"-110", []Element{6, 14}},
				input{"10-1", []Element{9, 11}},
				input{"1-0-", []Element{8, 9, 12, 13}},
				input{"1-01", []Element{9, 13}},
			),
			m: fromInputs(
				input{"00-0", []Element{0, 2}},
				input{"001-", []Element{2, 3}},
				input{"0-10", []Element{2}},
				input{"-000", []Element{0}},
				input{"-011", []Element{3, 11}},
				input{"10-1", []Element{11}},
				input{"1-01", []Element{}},
			),
			essential: smap("-101", "-110", "1-0-"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"00-0", []Element{0, 2}},
				input{"001-", []Element{2, 3}},
				input{"0-10", []Element{2, 6}},
				input{"-000", []Element{0, 8}},
				input{"-011", []Element{3, 11}},
				input{"-101", []Element{5, 13}},
				input{"-110", []Element{6, 14}},
				input{"10-1", []Element{9, 11}},
				input{"1-0-", []Element{8, 9, 12, 13}},
				input{"1-01", []Element{9, 13}},
			),
			m: bipartite.New(),

			essential: smap("-101", "-110", "00-0", "-011", "1-0-"),
		},
		simok: true,
		min:   [][]Subset{{"-101", "-110", "00-0", "-011", "1-0-"}},
	},
	"seven-segment G": {
		c: &Cover{
			in: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-01-", []Element{2, 3, 10, 11}},
				input{"-101", []Element{5, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
				input{"1--1", []Element{9, 11, 13, 15}},
				input{"1-1-", []Element{10, 11, 14, 15}},
			),
			m: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-01-", []Element{2, 3, 10, 11}},
				input{"-101", []Element{5, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
				input{"1--1", []Element{9, 11, 13, 15}},
				input{"1-1-", []Element{10, 11, 14, 15}},
			),
			essential: smap(),
		},
		s: &Cover{
			in: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-01-", []Element{2, 3, 10, 11}},
				input{"-101", []Element{5, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
				input{"1--1", []Element{9, 11, 13, 15}},
				input{"1-1-", []Element{10, 11, 14, 15}},
			),
			m: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-01-", []Element{2, 3, 10, 11}},
				input{"-101", []Element{5, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
				input{"1--1", []Element{9, 11, 13, 15}},
				input{"1-1-", []Element{10, 11, 14, 15}},
			),
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			in: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-01-", []Element{2, 3, 10, 11}},
				input{"-101", []Element{5, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
				input{"1--1", []Element{9, 11, 13, 15}},
				input{"1-1-", []Element{10, 11, 14, 15}},
			),
			m: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{6, 14}},
				input{"-101", []Element{5, 13}},
				input{"1--1", []Element{13, 15}},
				input{"1-1-", []Element{14, 15}},
			),
			essential: smap("-01-", "10--"),
		},
		eok: true,
		sim: &Cover{
			in: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{2, 6, 10, 14}},
				input{"-01-", []Element{2, 3, 10, 11}},
				input{"-101", []Element{5, 13}},
				input{"10--", []Element{8, 9, 10, 11}},
				input{"1--1", []Element{9, 11, 13, 15}},
				input{"1-1-", []Element{10, 11, 14, 15}},
			),
			m: fromInputs(
				input{"010-", []Element{4, 5}},
				input{"01-0", []Element{4, 6}},
				input{"--10", []Element{6, 14}},
				input{"-101", []Element{5, 13}},
				input{"1-1-", []Element{14, 15}},
				input{"1--1", []Element{13, 15}},
			),
			essential: smap("-01-", "10--"),
		},
		simok: false,
		min: [][]Subset{
			{"-01-", "10--", "010-", "--10", "1--1"},
			{"-01-", "10--", "01-0", "-101", "1-1-"},
		},
	},
}
