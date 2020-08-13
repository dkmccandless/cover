package cover

import (
	"reflect"
	"testing"
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

// fromInputs returns an initialized Cover populated with inputs.
func fromInputs(inputs ...input) *Cover {
	c := New()
	for _, in := range inputs {
		c.Add(in.s, in.es...)
	}
	c.initialize()
	return c
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
				inss: map[Subset]eset{true: emap(true)},
				ines: map[Element]sset{true: smap(true)},
				ss:   map[Subset]eset{},
				es:   map[Element]sset{},

				essential: smap(),
			},
		},
		// Subset containing many elements
		{
			[]input{
				{"Powers of 2", []Element{1, 2, 4, 8}},
			},
			&Cover{
				inss: map[Subset]eset{
					"Powers of 2": emap(1, 2, 4, 8),
				},
				ines: map[Element]sset{
					1: smap("Powers of 2"),
					2: smap("Powers of 2"),
					4: smap("Powers of 2"),
					8: smap("Powers of 2"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
				inss: map[Subset]eset{
					"Powers of 2": emap(1, 2, 4, 8),
				},
				ines: map[Element]sset{
					1: smap("Powers of 2"),
					2: smap("Powers of 2"),
					4: smap("Powers of 2"),
					8: smap("Powers of 2"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
				inss: map[Subset]eset{
					"Powers of 2": emap(1, 2, 4, 8),
				},
				ines: map[Element]sset{
					1: smap("Powers of 2"),
					2: smap("Powers of 2"),
					4: smap("Powers of 2"),
					8: smap("Powers of 2"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
				inss: map[Subset]eset{
					"Powers of 2": emap(2),
					"Even primes": emap(2),
				},
				ines: map[Element]sset{
					2: smap("Powers of 2", "Even primes"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
				inss: map[Subset]eset{
					"Powers of 2":       emap(1, 2, 4, 8),
					"Fibonacci numbers": emap(0, 1, 2, 3, 5, 8),
				},
				ines: map[Element]sset{
					0: smap("Fibonacci numbers"),
					1: smap("Powers of 2", "Fibonacci numbers"),
					2: smap("Powers of 2", "Fibonacci numbers"),
					3: smap("Fibonacci numbers"),
					4: smap("Powers of 2"),
					5: smap("Fibonacci numbers"),
					8: smap("Powers of 2", "Fibonacci numbers"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
				inss: map[Subset]eset{
					"Powers of 2":       emap(1, 2, 4, 8),
					"Fibonacci numbers": emap(0, 1, 2, 3, 5, 8),
				},
				ines: map[Element]sset{
					0: smap("Fibonacci numbers"),
					1: smap("Powers of 2", "Fibonacci numbers"),
					2: smap("Powers of 2", "Fibonacci numbers"),
					3: smap("Fibonacci numbers"),
					4: smap("Powers of 2"),
					5: smap("Fibonacci numbers"),
					8: smap("Powers of 2", "Fibonacci numbers"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
				inss: map[Subset]eset{
					"Powers of 2":       emap(1, 2, 4, 8),
					"Fibonacci numbers": emap(0, 1, 2, 3, 5, 8, 13),
				},
				ines: map[Element]sset{
					0:  smap("Fibonacci numbers"),
					1:  smap("Powers of 2", "Fibonacci numbers"),
					2:  smap("Powers of 2", "Fibonacci numbers"),
					3:  smap("Fibonacci numbers"),
					4:  smap("Powers of 2"),
					5:  smap("Fibonacci numbers"),
					8:  smap("Powers of 2", "Fibonacci numbers"),
					13: smap("Fibonacci numbers"),
				},
				ss: map[Subset]eset{},
				es: map[Element]sset{},

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
			fromInputs(input{true, []Element{true}}),
			map[Subset]sset{},
		},
		{
			fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"y"}},
			),
			map[Subset]sset{},
		},
		{
			fromInputs(
				input{"A", []Element{"x"}},
				input{"B", []Element{"x", "y", "z"}},
			),
			map[Subset]sset{"B": smap("A")},
		},
		{
			fromInputs(
				input{"A", []Element{2}},
				input{"B", []Element{2, 6}},
				input{"C", []Element{2, 6}},
				input{"D", []Element{1, 2, 4}},
				input{"E", []Element{3, 5, 7}},
				input{"F", []Element{0, 1, 2, 4, 7}},
			),
			map[Subset]sset{
				"B": smap("A"),
				"C": smap("A"),
				"D": smap("A"),
				"F": smap("A", "D"),
			},
		},
	} {
		for a := range test.c.ss {
			for b := range test.c.ss {
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
	cc := New()
	for s := range c.inss {
		cc.inss[s] = c.inss[s].copy()
	}
	for e := range c.ines {
		cc.ines[e] = c.ines[e].copy()
	}
	for s := range c.ss {
		cc.ss[s] = c.ss[s].copy()
	}
	for e := range c.es {
		cc.es[e] = c.es[e].copy()
	}
	cc.essential = c.essential.copy()
	return cc
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
		c: fromInputs(input{true, []Element{true}}),
		s: &Cover{
			inss: map[Subset]eset{true: emap(true)},
			ines: map[Element]sset{true: smap(true)},
			ss:   map[Subset]eset{true: emap(true)},
			es:   map[Element]sset{true: smap(true)},

			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{true: emap(true)},
			ines: map[Element]sset{true: smap(true)},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap(true),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{true: emap(true)},
			ines: map[Element]sset{true: smap(true)},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap(true),
		},
		simok: true,
		min:   [][]Subset{{true}},
	},
	"disjoint A and B": {
		c: fromInputs(
			input{"A", []Element{"x"}},
			input{"B", []Element{"y"}},
		),
		s: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("y")},
			ines: map[Element]sset{"x": smap("A"), "y": smap("B")},
			ss:   map[Subset]eset{"A": emap("x"), "B": emap("y")},
			es:   map[Element]sset{"x": smap("A"), "y": smap("B")},

			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("y")},
			ines: map[Element]sset{"x": smap("A"), "y": smap("B")},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap("A", "B"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("y")},
			ines: map[Element]sset{"x": smap("A"), "y": smap("B")},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap("A", "B"),
		},
		simok: true,
		min:   [][]Subset{{"A", "B"}},
	},
	"1 Subset contains 2 Elements": {
		c: fromInputs(input{"A", []Element{"x", "y"}}),
		s: &Cover{
			inss: map[Subset]eset{"A": emap("x", "y")},
			ines: map[Element]sset{"x": smap("A"), "y": smap("A")},
			ss:   map[Subset]eset{"A": emap("x", "y")},
			es:   map[Element]sset{"x": smap("A"), "y": smap("A")},

			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{"A": emap("x", "y")},
			ines: map[Element]sset{"x": smap("A"), "y": smap("A")},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap("A"),
		},
		eok: true,
		sim: &Cover{inss: map[Subset]eset{"A": emap("x", "y")},
			ines: map[Element]sset{"x": smap("A"), "y": smap("A")},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap("A"),
		},
		simok: true,
		min:   [][]Subset{{"A"}},
	},
	"2 Subsets contain 1 Element": {
		c: fromInputs(
			input{"A", []Element{"x"}},
			input{"B", []Element{"x"}},
		),
		s: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("x")},
			ines: map[Element]sset{"x": smap("A", "B")},
			ss:   map[Subset]eset{"A": emap("x"), "B": emap("x")},
			es:   map[Element]sset{"x": smap("A", "B")},

			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("x")},
			ines: map[Element]sset{"x": smap("A", "B")},
			ss:   map[Subset]eset{"A": emap("x"), "B": emap("x")},
			es:   map[Element]sset{"x": smap("A", "B")},

			essential: smap(),
		},
		eok: false,
		sim: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("x")},
			ines: map[Element]sset{"x": smap("A", "B")},
			ss:   map[Subset]eset{"A": emap("x"), "B": emap("x")},
			es:   map[Element]sset{"x": smap("A", "B")},

			essential: smap(),
		},
		simok: false,
		min:   [][]Subset{{"A"}, {"B"}},
	},
	"B contains A": {
		c: fromInputs(
			input{"A", []Element{"x"}},
			input{"B", []Element{"x", "y", "z"}},
		),
		s: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("x", "y", "z")},
			ines: map[Element]sset{"x": smap("A", "B"), "y": smap("B"), "z": smap("B")},
			ss:   map[Subset]eset{"B": emap("x", "y", "z")},
			es:   map[Element]sset{"x": smap("B"), "y": smap("B"), "z": smap("B")},

			essential: smap(),
		},
		sok: true,
		e: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("x", "y", "z")},
			ines: map[Element]sset{"x": smap("A", "B"), "y": smap("B"), "z": smap("B")},
			ss:   map[Subset]eset{"A": emap()},
			es:   map[Element]sset{},

			essential: smap("B"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{"A": emap("x"), "B": emap("x", "y", "z")},
			ines: map[Element]sset{"x": smap("A", "B"), "y": smap("B"), "z": smap("B")},
			ss:   map[Subset]eset{},
			es:   map[Element]sset{},

			essential: smap("B"),
		},
		simok: true,
		min:   [][]Subset{{"B"}},
	},
	"seven-segment A": {
		c: fromInputs(
			input{"0-1-", []Element{2, 3, 6, 7}},
			input{"01-1", []Element{5, 7}},
			input{"-0-0", []Element{0, 2, 8, 10}},
			input{"--10", []Element{2, 6, 10, 14}},
			input{"-11-", []Element{6, 7, 14, 15}},
			input{"100-", []Element{8, 9}},
			input{"1--0", []Element{8, 10, 12, 14}},
			input{"11-0", []Element{12, 14}},
		),
		s: &Cover{
			inss: map[Subset]eset{
				"0-1-": emap(2, 3, 6, 7),
				"01-1": emap(5, 7),
				"-0-0": emap(0, 2, 8, 10),
				"--10": emap(2, 6, 10, 14),
				"-11-": emap(6, 7, 14, 15),
				"100-": emap(8, 9),
				"1--0": emap(8, 10, 12, 14),
				"11-0": emap(12, 14),
			},
			ines: map[Element]sset{
				0:  smap("-0-0"),
				2:  smap("0-1-", "-0-0", "--10"),
				3:  smap("0-1-"),
				5:  smap("01-1"),
				6:  smap("0-1-", "--10", "-11-"),
				7:  smap("0-1-", "01-1", "-11-"),
				8:  smap("-0-0", "100-", "1--0"),
				9:  smap("100-"),
				10: smap("-0-0", "--10", "1--0"),
				12: smap("1--0", "11-0"),
				14: smap("--10", "-11-", "1--0", "11-0"),
				15: smap("-11-"),
			},
			ss: map[Subset]eset{
				"0-1-": emap(2, 3, 6, 7),
				"01-1": emap(5, 7),
				"-0-0": emap(0, 2, 8, 10),
				"--10": emap(2, 6, 10, 14),
				"-11-": emap(6, 7, 14, 15),
				"100-": emap(8, 9),
				"1--0": emap(8, 10, 12, 14),
			},
			es: map[Element]sset{
				0:  smap("-0-0"),
				2:  smap("0-1-", "-0-0", "--10"),
				3:  smap("0-1-"),
				5:  smap("01-1"),
				6:  smap("0-1-", "--10", "-11-"),
				7:  smap("0-1-", "01-1", "-11-"),
				8:  smap("-0-0", "100-", "1--0"),
				9:  smap("100-"),
				10: smap("-0-0", "--10", "1--0"),
				12: smap("1--0"),
				14: smap("--10", "-11-", "1--0"),
				15: smap("-11-"),
			},
			essential: smap(),
		},
		sok: true,
		e: &Cover{
			inss: map[Subset]eset{
				"0-1-": emap(2, 3, 6, 7),
				"01-1": emap(5, 7),
				"-0-0": emap(0, 2, 8, 10),
				"--10": emap(2, 6, 10, 14),
				"-11-": emap(6, 7, 14, 15),
				"100-": emap(8, 9),
				"1--0": emap(8, 10, 12, 14),
				"11-0": emap(12, 14),
			},
			ines: map[Element]sset{
				0:  smap("-0-0"),
				2:  smap("0-1-", "-0-0", "--10"),
				3:  smap("0-1-"),
				5:  smap("01-1"),
				6:  smap("0-1-", "--10", "-11-"),
				7:  smap("0-1-", "01-1", "-11-"),
				8:  smap("-0-0", "100-", "1--0"),
				9:  smap("100-"),
				10: smap("-0-0", "--10", "1--0"),
				12: smap("1--0", "11-0"),
				14: smap("--10", "-11-", "1--0", "11-0"),
				15: smap("-11-"),
			},
			ss: map[Subset]eset{"--10": emap(), "1--0": emap(12), "11-0": emap(12)},
			es: map[Element]sset{12: smap("1--0", "11-0")},

			essential: smap("0-1-", "01-1", "-0-0", "-11-", "100-"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{
				"0-1-": emap(2, 3, 6, 7),
				"01-1": emap(5, 7),
				"-0-0": emap(0, 2, 8, 10),
				"--10": emap(2, 6, 10, 14),
				"-11-": emap(6, 7, 14, 15),
				"100-": emap(8, 9),
				"1--0": emap(8, 10, 12, 14),
				"11-0": emap(12, 14),
			},
			ines: map[Element]sset{
				0:  smap("-0-0"),
				2:  smap("0-1-", "-0-0", "--10"),
				3:  smap("0-1-"),
				5:  smap("01-1"),
				6:  smap("0-1-", "--10", "-11-"),
				7:  smap("0-1-", "01-1", "-11-"),
				8:  smap("-0-0", "100-", "1--0"),
				9:  smap("100-"),
				10: smap("-0-0", "--10", "1--0"),
				12: smap("1--0", "11-0"),
				14: smap("--10", "-11-", "1--0", "11-0"),
				15: smap("-11-"),
			},
			ss: map[Subset]eset{"--10": emap()},
			es: map[Element]sset{},

			essential: smap("0-1-", "01-1", "-0-0", "-11-", "100-", "1--0"),
		},
		simok: true,
		min:   [][]Subset{{"0-1-", "01-1", "-0-0", "-11-", "100-", "1--0"}},
	},
	"seven-segment B": {
		c: fromInputs(
			input{"00--", []Element{0, 1, 2, 3}},
			input{"0-00", []Element{0, 4}},
			input{"0-11", []Element{3, 7}},
			input{"-00-", []Element{0, 1, 8, 9}},
			input{"-0-0", []Element{0, 2, 8, 10}},
			input{"1-01", []Element{9, 13}},
		),
		s: &Cover{
			inss: map[Subset]eset{
				"00--": emap(0, 1, 2, 3),
				"0-00": emap(0, 4),
				"0-11": emap(3, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-0": emap(0, 2, 8, 10),
				"1-01": emap(9, 13),
			},
			ines: map[Element]sset{
				0:  smap("00--", "0-00", "-00-", "-0-0"),
				1:  smap("00--", "-00-"),
				2:  smap("00--", "-0-0"),
				3:  smap("00--", "0-11"),
				4:  smap("0-00"),
				7:  smap("0-11"),
				8:  smap("-00-", "-0-0"),
				9:  smap("-00-", "1-01"),
				10: smap("-0-0"),
				13: smap("1-01"),
			},
			ss: map[Subset]eset{
				"00--": emap(0, 1, 2, 3),
				"0-00": emap(0, 4),
				"0-11": emap(3, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-0": emap(0, 2, 8, 10),
				"1-01": emap(9, 13),
			},
			es: map[Element]sset{
				0:  smap("00--", "0-00", "-00-", "-0-0"),
				1:  smap("00--", "-00-"),
				2:  smap("00--", "-0-0"),
				3:  smap("00--", "0-11"),
				4:  smap("0-00"),
				7:  smap("0-11"),
				8:  smap("-00-", "-0-0"),
				9:  smap("-00-", "1-01"),
				10: smap("-0-0"),
				13: smap("1-01"),
			},
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{
				"00--": emap(0, 1, 2, 3),
				"0-00": emap(0, 4),
				"0-11": emap(3, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-0": emap(0, 2, 8, 10),
				"1-01": emap(9, 13),
			},
			ines: map[Element]sset{
				0:  smap("00--", "0-00", "-00-", "-0-0"),
				1:  smap("00--", "-00-"),
				2:  smap("00--", "-0-0"),
				3:  smap("00--", "0-11"),
				4:  smap("0-00"),
				7:  smap("0-11"),
				8:  smap("-00-", "-0-0"),
				9:  smap("-00-", "1-01"),
				10: smap("-0-0"),
				13: smap("1-01"),
			},
			ss: map[Subset]eset{
				"00--": emap(1),
				"-00-": emap(1),
			},
			es: map[Element]sset{
				1: smap("00--", "-00-"),
			},
			essential: smap("0-00", "0-11", "-0-0", "1-01"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{
				"00--": emap(0, 1, 2, 3),
				"0-00": emap(0, 4),
				"0-11": emap(3, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-0": emap(0, 2, 8, 10),
				"1-01": emap(9, 13),
			},
			ines: map[Element]sset{
				0:  smap("00--", "0-00", "-00-", "-0-0"),
				1:  smap("00--", "-00-"),
				2:  smap("00--", "-0-0"),
				3:  smap("00--", "0-11"),
				4:  smap("0-00"),
				7:  smap("0-11"),
				8:  smap("-00-", "-0-0"),
				9:  smap("-00-", "1-01"),
				10: smap("-0-0"),
				13: smap("1-01"),
			},
			ss: map[Subset]eset{
				"00--": emap(1),
				"-00-": emap(1),
			},
			es: map[Element]sset{
				1: smap("00--", "-00-"),
			},
			essential: smap("0-00", "0-11", "-0-0", "1-01"),
		},
		simok: false,
		min: [][]Subset{
			{"0-00", "0-11", "-0-0", "1-01", "00--"},
			{"0-00", "0-11", "-0-0", "1-01", "-00-"},
		},
	},
	"seven-segment C": {
		c: fromInputs(
			input{"0-0-", []Element{0, 1, 4, 5}},
			input{"0--1", []Element{1, 3, 5, 7}},
			input{"01--", []Element{4, 5, 6, 7}},
			input{"-00-", []Element{0, 1, 8, 9}},
			input{"-0-1", []Element{1, 3, 9, 11}},
			input{"--01", []Element{1, 5, 9, 13}},
			input{"10--", []Element{8, 9, 10, 11}},
		),
		s: &Cover{
			inss: map[Subset]eset{
				"0-0-": emap(0, 1, 4, 5),
				"0--1": emap(1, 3, 5, 7),
				"01--": emap(4, 5, 6, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-1": emap(1, 3, 9, 11),
				"--01": emap(1, 5, 9, 13),
				"10--": emap(8, 9, 10, 11),
			},
			ines: map[Element]sset{
				0:  smap("0-0-", "-00-"),
				1:  smap("0-0-", "0--1", "-00-", "-0-1", "--01"),
				3:  smap("0--1", "-0-1"),
				4:  smap("0-0-", "01--"),
				5:  smap("0-0-", "0--1", "01--", "--01"),
				6:  smap("01--"),
				7:  smap("0--1", "01--"),
				8:  smap("-00-", "10--"),
				9:  smap("-00-", "-0-1", "--01", "10--"),
				10: smap("10--"),
				11: smap("-0-1", "10--"),
				13: smap("--01"),
			},
			ss: map[Subset]eset{
				"0-0-": emap(0, 1, 4, 5),
				"0--1": emap(1, 3, 5, 7),
				"01--": emap(4, 5, 6, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-1": emap(1, 3, 9, 11),
				"--01": emap(1, 5, 9, 13),
				"10--": emap(8, 9, 10, 11),
			},
			es: map[Element]sset{
				0:  smap("0-0-", "-00-"),
				1:  smap("0-0-", "0--1", "-00-", "-0-1", "--01"),
				3:  smap("0--1", "-0-1"),
				4:  smap("0-0-", "01--"),
				5:  smap("0-0-", "0--1", "01--", "--01"),
				6:  smap("01--"),
				7:  smap("0--1", "01--"),
				8:  smap("-00-", "10--"),
				9:  smap("-00-", "-0-1", "--01", "10--"),
				10: smap("10--"),
				11: smap("-0-1", "10--"),
				13: smap("--01"),
			},
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{
				"0-0-": emap(0, 1, 4, 5),
				"0--1": emap(1, 3, 5, 7),
				"01--": emap(4, 5, 6, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-1": emap(1, 3, 9, 11),
				"--01": emap(1, 5, 9, 13),
				"10--": emap(8, 9, 10, 11),
			},
			ines: map[Element]sset{
				0:  smap("0-0-", "-00-"),
				1:  smap("0-0-", "0--1", "-00-", "-0-1", "--01"),
				3:  smap("0--1", "-0-1"),
				4:  smap("0-0-", "01--"),
				5:  smap("0-0-", "0--1", "01--", "--01"),
				6:  smap("01--"),
				7:  smap("0--1", "01--"),
				8:  smap("-00-", "10--"),
				9:  smap("-00-", "-0-1", "--01", "10--"),
				10: smap("10--"),
				11: smap("-0-1", "10--"),
				13: smap("--01"),
			},
			ss: map[Subset]eset{
				"0-0-": emap(0),
				"0--1": emap(3),
				"-00-": emap(0),
				"-0-1": emap(3),
			},
			es: map[Element]sset{
				0: smap("0-0-", "-00-"),
				3: smap("0--1", "-0-1"),
			},
			essential: smap("01--", "--01", "10--"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{
				"0-0-": emap(0, 1, 4, 5),
				"0--1": emap(1, 3, 5, 7),
				"01--": emap(4, 5, 6, 7),
				"-00-": emap(0, 1, 8, 9),
				"-0-1": emap(1, 3, 9, 11),
				"--01": emap(1, 5, 9, 13),
				"10--": emap(8, 9, 10, 11),
			},
			ines: map[Element]sset{
				0:  smap("0-0-", "-00-"),
				1:  smap("0-0-", "0--1", "-00-", "-0-1", "--01"),
				3:  smap("0--1", "-0-1"),
				4:  smap("0-0-", "01--"),
				5:  smap("0-0-", "0--1", "01--", "--01"),
				6:  smap("01--"),
				7:  smap("0--1", "01--"),
				8:  smap("-00-", "10--"),
				9:  smap("-00-", "-0-1", "--01", "10--"),
				10: smap("10--"),
				11: smap("-0-1", "10--"),
				13: smap("--01"),
			},
			ss: map[Subset]eset{
				"0-0-": emap(0),
				"0--1": emap(3),
				"-00-": emap(0),
				"-0-1": emap(3),
			},
			es: map[Element]sset{
				0: smap("0-0-", "-00-"),
				3: smap("0--1", "-0-1"),
			},
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
		c: fromInputs(
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
		s: &Cover{
			inss: map[Subset]eset{
				"00-0": emap(0, 2),
				"001-": emap(2, 3),
				"0-10": emap(2, 6),
				"-000": emap(0, 8),
				"-011": emap(3, 11),
				"-101": emap(5, 13),
				"-110": emap(6, 14),
				"10-1": emap(9, 11),
				"1-0-": emap(8, 9, 12, 13),
				"1-01": emap(9, 13),
			},
			ines: map[Element]sset{
				0:  smap("00-0", "-000"),
				2:  smap("00-0", "001-", "0-10"),
				3:  smap("001-", "-011"),
				5:  smap("-101"),
				6:  smap("0-10", "-110"),
				8:  smap("-000", "1-0-"),
				9:  smap("10-1", "1-0-", "1-01"),
				11: smap("-011", "10-1"),
				12: smap("1-0-"),
				13: smap("-101", "1-0-", "1-01"),
				14: smap("-110"),
			},
			ss: map[Subset]eset{
				"00-0": emap(0, 2),
				"001-": emap(2, 3),
				"0-10": emap(2, 6),
				"-000": emap(0, 8),
				"-011": emap(3, 11),
				"-101": emap(5, 13),
				"-110": emap(6, 14),
				"10-1": emap(9, 11),
				"1-0-": emap(8, 9, 12, 13),
			},
			es: map[Element]sset{
				0:  smap("00-0", "-000"),
				2:  smap("00-0", "001-", "0-10"),
				3:  smap("001-", "-011"),
				5:  smap("-101"),
				6:  smap("0-10", "-110"),
				8:  smap("-000", "1-0-"),
				9:  smap("10-1", "1-0-"),
				11: smap("-011", "10-1"),
				12: smap("1-0-"),
				13: smap("-101", "1-0-"),
				14: smap("-110"),
			},
			essential: smap(),
		},
		sok: true,
		e: &Cover{
			inss: map[Subset]eset{
				"00-0": emap(0, 2),
				"001-": emap(2, 3),
				"0-10": emap(2, 6),
				"-000": emap(0, 8),
				"-011": emap(3, 11),
				"-101": emap(5, 13),
				"-110": emap(6, 14),
				"10-1": emap(9, 11),
				"1-0-": emap(8, 9, 12, 13),
				"1-01": emap(9, 13),
			},
			ines: map[Element]sset{
				0:  smap("00-0", "-000"),
				2:  smap("00-0", "001-", "0-10"),
				3:  smap("001-", "-011"),
				5:  smap("-101"),
				6:  smap("0-10", "-110"),
				8:  smap("-000", "1-0-"),
				9:  smap("10-1", "1-0-", "1-01"),
				11: smap("-011", "10-1"),
				12: smap("1-0-"),
				13: smap("-101", "1-0-", "1-01"),
				14: smap("-110"),
			},
			ss: map[Subset]eset{
				"00-0": emap(0, 2),
				"001-": emap(2, 3),
				"0-10": emap(2),
				"-000": emap(0),
				"-011": emap(3, 11),
				"10-1": emap(11),
				"1-01": emap(),
			},
			es: map[Element]sset{
				0:  smap("00-0", "-000"),
				2:  smap("00-0", "001-", "0-10"),
				3:  smap("001-", "-011"),
				11: smap("-011", "10-1"),
			},
			essential: smap("-101", "-110", "1-0-"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{
				"00-0": emap(0, 2),
				"001-": emap(2, 3),
				"0-10": emap(2, 6),
				"-000": emap(0, 8),
				"-011": emap(3, 11),
				"-101": emap(5, 13),
				"-110": emap(6, 14),
				"10-1": emap(9, 11),
				"1-0-": emap(8, 9, 12, 13),
				"1-01": emap(9, 13),
			},
			ines: map[Element]sset{
				0:  smap("00-0", "-000"),
				2:  smap("00-0", "001-", "0-10"),
				3:  smap("001-", "-011"),
				5:  smap("-101"),
				6:  smap("0-10", "-110"),
				8:  smap("-000", "1-0-"),
				9:  smap("10-1", "1-0-", "1-01"),
				11: smap("-011", "10-1"),
				12: smap("1-0-"),
				13: smap("-101", "1-0-", "1-01"),
				14: smap("-110"),
			},
			ss: map[Subset]eset{"001-": emap()},
			es: map[Element]sset{},

			essential: smap("-101", "-110", "00-0", "-011", "1-0-"),
		},
		simok: true,
		min:   [][]Subset{{"-101", "-110", "00-0", "-011", "1-0-"}},
	},
	"seven-segment G": {
		c: fromInputs(
			input{"010-", []Element{4, 5}},
			input{"01-0", []Element{4, 6}},
			input{"--10", []Element{2, 6, 10, 14}},
			input{"-01-", []Element{2, 3, 10, 11}},
			input{"-101", []Element{5, 13}},
			input{"10--", []Element{8, 9, 10, 11}},
			input{"1--1", []Element{9, 11, 13, 15}},
			input{"1-1-", []Element{10, 11, 14, 15}},
		),
		s: &Cover{
			inss: map[Subset]eset{
				"010-": emap(4, 5),
				"01-0": emap(4, 6),
				"--10": emap(2, 6, 10, 14),
				"-01-": emap(2, 3, 10, 11),
				"-101": emap(5, 13),
				"10--": emap(8, 9, 10, 11),
				"1--1": emap(9, 11, 13, 15),
				"1-1-": emap(10, 11, 14, 15),
			},
			ines: map[Element]sset{
				2:  smap("--10", "-01-"),
				3:  smap("-01-"),
				4:  smap("010-", "01-0"),
				5:  smap("-101", "010-"),
				6:  smap("01-0", "--10"),
				8:  smap("10--"),
				9:  smap("10--", "1--1"),
				10: smap("--10", "-01-", "10--", "1-1-"),
				11: smap("-01-", "10--", "1--1", "1-1-"),
				13: smap("-101", "1--1"),
				14: smap("--10", "1-1-"),
				15: smap("1--1", "1-1-"),
			},
			ss: map[Subset]eset{
				"010-": emap(4, 5),
				"01-0": emap(4, 6),
				"--10": emap(2, 6, 10, 14),
				"-01-": emap(2, 3, 10, 11),
				"-101": emap(5, 13),
				"10--": emap(8, 9, 10, 11),
				"1--1": emap(9, 11, 13, 15),
				"1-1-": emap(10, 11, 14, 15),
			},
			es: map[Element]sset{
				2:  smap("--10", "-01-"),
				3:  smap("-01-"),
				4:  smap("010-", "01-0"),
				5:  smap("-101", "010-"),
				6:  smap("01-0", "--10"),
				8:  smap("10--"),
				9:  smap("10--", "1--1"),
				10: smap("--10", "-01-", "10--", "1-1-"),
				11: smap("-01-", "10--", "1--1", "1-1-"),
				13: smap("-101", "1--1"),
				14: smap("--10", "1-1-"),
				15: smap("1--1", "1-1-"),
			},
			essential: smap(),
		},
		sok: false,
		e: &Cover{
			inss: map[Subset]eset{
				"010-": emap(4, 5),
				"01-0": emap(4, 6),
				"--10": emap(2, 6, 10, 14),
				"-01-": emap(2, 3, 10, 11),
				"-101": emap(5, 13),
				"10--": emap(8, 9, 10, 11),
				"1--1": emap(9, 11, 13, 15),
				"1-1-": emap(10, 11, 14, 15),
			},
			ines: map[Element]sset{
				2:  smap("--10", "-01-"),
				3:  smap("-01-"),
				4:  smap("010-", "01-0"),
				5:  smap("-101", "010-"),
				6:  smap("01-0", "--10"),
				8:  smap("10--"),
				9:  smap("10--", "1--1"),
				10: smap("--10", "-01-", "10--", "1-1-"),
				11: smap("-01-", "10--", "1--1", "1-1-"),
				13: smap("-101", "1--1"),
				14: smap("--10", "1-1-"),
				15: smap("1--1", "1-1-"),
			},
			ss: map[Subset]eset{
				"010-": emap(4, 5),
				"01-0": emap(4, 6),
				"--10": emap(6, 14),
				"-101": emap(5, 13),
				"1--1": emap(13, 15),
				"1-1-": emap(14, 15),
			},
			es: map[Element]sset{
				4:  smap("010-", "01-0"),
				5:  smap("-101", "010-"),
				6:  smap("01-0", "--10"),
				13: smap("-101", "1--1"),
				14: smap("--10", "1-1-"),
				15: smap("1--1", "1-1-"),
			},
			essential: smap("-01-", "10--"),
		},
		eok: true,
		sim: &Cover{
			inss: map[Subset]eset{
				"010-": emap(4, 5),
				"01-0": emap(4, 6),
				"--10": emap(2, 6, 10, 14),
				"-01-": emap(2, 3, 10, 11),
				"-101": emap(5, 13),
				"10--": emap(8, 9, 10, 11),
				"1--1": emap(9, 11, 13, 15),
				"1-1-": emap(10, 11, 14, 15),
			},
			ines: map[Element]sset{
				2:  smap("--10", "-01-"),
				3:  smap("-01-"),
				4:  smap("010-", "01-0"),
				5:  smap("-101", "010-"),
				6:  smap("01-0", "--10"),
				8:  smap("10--"),
				9:  smap("10--", "1--1"),
				10: smap("--10", "-01-", "10--", "1-1-"),
				11: smap("-01-", "10--", "1--1", "1-1-"),
				13: smap("-101", "1--1"),
				14: smap("--10", "1-1-"),
				15: smap("1--1", "1-1-"),
			},
			ss: map[Subset]eset{
				"010-": emap(4, 5),
				"01-0": emap(4, 6),
				"--10": emap(6, 14),
				"-101": emap(5, 13),
				"1-1-": emap(14, 15),
				"1--1": emap(13, 15),
			},
			es: map[Element]sset{
				4:  smap("010-", "01-0"),
				5:  smap("-101", "010-"),
				6:  smap("01-0", "--10"),
				13: smap("-101", "1--1"),
				14: smap("--10", "1-1-"),
				15: smap("1--1", "1-1-"),
			},

			essential: smap("-01-", "10--"),
		},
		simok: false,
		min: [][]Subset{
			{"-01-", "10--", "010-", "--10", "1--1"},
			{"-01-", "10--", "01-0", "-101", "1-1-"},
		},
	},
}
