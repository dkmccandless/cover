// Package cover implements an algorithm to solve the minimum set cover problem.
package cover

import (
	"sort"
)

// Element is contained by one or more Subsets.
type Element interface{}

// Subset contains one or more Elements.
type Subset interface{}

// eset is a set of Elements.
type eset map[Element]struct{}

// copy returns a copy of es that shares no memory with it.
func (es eset) copy() eset {
	m := make(eset)
	for e := range es {
		m[e] = struct{}{}
	}
	return m
}

// sset is a set of Subsets.
type sset map[Subset]struct{}

// copy returns a copy of ss that shares no memory with it.
func (ss sset) copy() sset {
	m := make(sset)
	for s := range ss {
		m[s] = struct{}{}
	}
	return m
}

// Cover records Subsets and the Elements they contain.
type Cover struct {
	// inss and ines store all added Subsets and Elements.
	// Minimize copies their contents into ss and es to modify.
	inss map[Subset]eset
	ines map[Element]sset

	// ss holds all Subsets not yet determined to be essential or dominated.
	// Minimize copies the contents of ss from inss and modifies them during simplification.
	ss map[Subset]eset

	// es holds all Elements not yet determined to be covered.
	// Minimize copies the contents of es from ines and modifies them during simplification.
	es map[Element]sset

	// essential contains the Subsets determined by Minimize to be necessary members of the covering set.
	essential sset
}

// New returns an empty Cover.
func New() *Cover {
	return &Cover{
		inss: make(map[Subset]eset),
		ines: make(map[Element]sset),
		ss:   make(map[Subset]eset),
		es:   make(map[Element]sset),

		essential: make(sset),
	}
}

// Add records that s contains es.
// If es is empty, Add is a no-op.
func (c *Cover) Add(s Subset, es ...Element) {
	if len(es) == 0 {
		return
	}
	if _, ok := c.inss[s]; !ok {
		c.inss[s] = make(eset, len(es))
	}
	for _, e := range es {
		if _, ok := c.ines[e]; !ok {
			c.ines[e] = make(sset)
		}
		c.inss[s][e] = struct{}{}
		c.ines[e][s] = struct{}{}
	}
	// Invariant: All maps are non-empty, and inss[s] contains e if and only if ines[e] contains s.
}

// initialize prepares c for minimization by copying c.inss into c.ss and c.ines into c.es and clearing c.essential.
func (c *Cover) initialize() {
	c.ss = make(map[Subset]eset, len(c.inss))
	for s := range c.inss {
		c.ss[s] = c.inss[s].copy()
	}
	c.es = make(map[Element]sset, len(c.ines))
	for e := range c.ines {
		c.es[e] = c.ines[e].copy()
	}
	c.essential = make(sset, len(c.ss))
}

// Minimize returns all minimum-length combinations of Subsets that cover every Element.
// In general, its complexity increases exponentially with the number of Elements.
func (c *Cover) Minimize() [][]Subset {
	c.initialize()

	isUnique := c.simplify()

	// ess holds the essential Subsets for returning as a slice.
	var ess []Subset
	for s := range c.essential {
		ess = append(ess, s)
	}
	if isUnique {
		// The essential Subsets constitute a unique covering set.
		return [][]Subset{ess}
	}

	// At least one non-essential Subset is required to cover at least one Element.
	// Search all Subset unions of length 1, then 2, and so on until covering sets are found.
	var covers [][]Subset
	ss := make([]Subset, 0, len(c.ss))
	for s := range c.ss {
		ss = append(ss, s)
	}
	// Sort the Subsets to search in order of coverage, starting with the largest.
	sort.Slice(ss, func(i, j int) bool { return len(c.ss[ss[i]]) > len(c.ss[ss[j]]) })

	for w := 1; w <= len(ss); w++ {
		b := make([]bool, len(ss))
		for i := 0; i < w; i++ {
			b[i] = true
		}
		for {
			var ok bool
			for e := range c.es {
				// Check whether any Subsets in ss cover e.
				// b[i] indicates whether to consider ss[i].
				ok = false
				for i := range ss {
					if !b[i] {
						continue
					}
					if _, ok = c.es[e][ss[i]]; ok {
						break
					}
				}
				if !ok {
					break
				}
			}

			if ok {
				// n encodes a valid covering set: all Elements are covered by at least one of the considered Subsets.
				cs := append([]Subset{}, ess...)
				for i := range ss {
					if !b[i] {
						continue
					}
					cs = append(cs, ss[i])
				}
				covers = append(covers, cs)
			}
			if !nextPerm(b) {
				break
			}
		}
		if len(covers) > 0 {
			break
		}
	}

	return covers
}

// nextPerm implements Knuth's Algorithm L to generate the next lexicographic permutation of b.
// It reports whether there are more permutations remaining.
func nextPerm(b []bool) bool {
	j := len(b) - 2
	for ; !b[j] || b[j+1]; j-- {
		if j == 0 {
			return false
		}
	}
	l := len(b) - 1
	for b[l] {
		l--
	}
	b[j], b[l] = b[l], b[j]
	for k, l := j+1, len(b)-1; k < l; k, l = k+1, l-1 {
		b[k], b[l] = b[l], b[k]
	}
	return true
}

// simplify simplifies c by identifying all essential Subsets.
// It reports whether the essential Subsets are sufficient to cover all Elements by themselves
// (and the covering set is therefore unique).
func (c *Cover) simplify() bool {
	// reduceS removes all dominated Subsets but may reveal another Subset as essential;
	// reduceE removes all essential Subsets and the Elements they contain, but may cause another Subset to become dominated.
	// Call them in alternation: c is fully simplified when either does not apply any reductions,
	// provided that they have both been called at least once.
	c.reduceS()
	for c.reduceE() && c.reduceS() {
	}
	return len(c.es) == 0
}

// reduceS reduces c by removing dominated Subsets and reports whether any Subsets were removed.
// When reduceP returns, c contains no dominated Subsets.
// The removal of a dominated Subset may expose another Subset as essential.
func (c *Cover) reduceS() bool {
	var ok bool
	for d := range c.ss {
		for s := range c.ss {
			if d == s || !c.dominates(d, s) {
				continue
			}
			// s will not appear in any minimal covering solution because d's coverage is a proper superset.
			c.removeS(s)
			ok = true
		}
	}
	return ok
}

// removeS removes s from c.ss and c.es.
func (c *Cover) removeS(s Subset) {
	for e := range c.ss[s] {
		delete(c.es[e], s)
	}
	delete(c.ss, s)
}

// dominates reports whether Subset a dominates Subset b; that is, whether a's Elements are a proper superset of b's.
func (c *Cover) dominates(a, b Subset) bool {
	for e := range c.ss[b] {
		if _, ok := c.ss[a][e]; !ok {
			return false
		}
	}
	return len(c.ss[a]) > len(c.ss[b])
}

// reduceE reduces c by identifying essential Subsets, moving them from c.ss to c.essential,
// and removing their Elements from c.es, and reports whether any Elements were removed.
// When reduceM returns, all Elements in c are contained by at least two Subsets.
// The removal of an Element may cause a Subset to become dominated.
func (c *Cover) reduceE() bool {
	var ok bool
	for e := range c.es {
		if len(c.es[e]) != 1 {
			continue
		}
		ok = true

		// e is contained by exactly one Subset, which is therefore essential.
		// Move it to c.essential and remove it and all Elements it covers.
		var s Subset
		for s = range c.es[e] {
		}
		for ee := range c.ss[s] {
			c.removeE(ee)
		}
		c.essential[s] = struct{}{}
		c.removeS(s)
	}
	return ok
}

// removeE removes e from c.
func (c *Cover) removeE(e Element) {
	for s := range c.es[e] {
		delete(c.ss[s], e)
	}
	delete(c.es, e)
}
