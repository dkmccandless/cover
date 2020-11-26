// Package cover implements an algorithm to solve the minimum set cover problem.
package cover

import (
	"sort"

	"github.com/dkmccandless/bipartite"
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
	// in stores all added Subsets and Elements.
	// Minimize copies their contents into m to modify.
	in *bipartite.Graph

	// m holds all Subsets not yet determined to be essential or dominated,
	// and all Elements not yet determined to be covered.
	// Minimize copies the contents of m from in and modifies them during simplification.
	m *bipartite.Graph

	// essential contains the Subsets determined by Minimize to be necessary members of the covering set.
	essential sset
}

// New returns an empty Cover.
func New() *Cover {
	return &Cover{
		in: bipartite.New(),
		m:  bipartite.New(),

		essential: make(sset),
	}
}

// Add records that s contains es.
// If es is empty, Add is a no-op.
func (c *Cover) Add(s Subset, es ...Element) {
	for _, e := range es {
		c.in.Add(s, e)
	}
}

// Minimize returns all minimum-length combinations of Subsets that cover every Element.
// In general, its complexity increases exponentially with the number of Elements.
func (c *Cover) Minimize() [][]Subset {
	c.m = bipartite.Copy(c.in)
	c.essential = make(sset, c.m.NA())

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
	ss := c.m.As()
	// Sort the Subsets to search in order of coverage, starting with the largest.
	sort.Slice(ss, func(i, j int) bool { return c.m.DegA(ss[i]) > c.m.DegA(ss[j]) })

	for w := 1; w <= len(ss); w++ {
		b := make([]bool, len(ss))
		for i := 0; i < w; i++ {
			b[i] = true
		}
		for {
			var ok bool
			for _, e := range c.m.Bs() {
				// Check whether any Subsets in ss cover e.
				// b[i] indicates whether to consider ss[i].
				ok = false
				for i, s := range ss {
					if !b[i] {
						continue
					}
					if ok = c.m.Adjacent(s, e); ok {
						break
					}
				}
				if !ok {
					break
				}
			}

			if ok {
				// b encodes a valid covering set: all Elements are covered by at least one of the considered Subsets.
				cs := append(make([]Subset, 0, len(ess)+w), ess...)
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
	if len(b) < 2 {
		return false
	}
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
	// provided that each has been called at least once.
	c.reduceS()
	for c.reduceE() && c.reduceS() {
	}
	return c.m.NB() == 0
}

// reduceS reduces c by removing dominated Subsets and reports whether any Subsets were removed.
// When reduceS returns, c contains no dominated Subsets.
// The removal of a dominated Subset may reveal another Subset as essential.
func (c *Cover) reduceS() bool {
	var ok bool
	for _, d := range c.m.As() {
		for _, s := range c.m.As() {
			if d == s || !c.dominates(d, s) {
				continue
			}
			// s will not appear in any minimal covering solution because d's coverage is a proper superset.
			c.m.RemoveA(s)
			ok = true
		}
	}
	return ok
}

// dominates reports whether d dominates s; that is, whether d's Elements are a proper superset of s's.
func (c *Cover) dominates(d, s Subset) bool {
	for _, e := range c.m.AdjToA(s) {
		if !c.m.Adjacent(d, e) {
			return false
		}
	}
	return c.m.DegA(d) > c.m.DegA(s)
}

// reduceE reduces c by identifying essential Subsets, moving them from c.m to c.essential,
// and removing their Elements from c.m, and reports whether any Elements were removed.
// When reduceE returns, all Elements in c are contained by at least two Subsets.
// The removal of an Element may cause a Subset to become dominated.
func (c *Cover) reduceE() bool {
	var ok bool
	for _, e := range c.m.Bs() {
		if c.m.DegB(e) != 1 {
			continue
		}
		ok = true

		// e is contained by exactly one Subset, which is therefore essential.
		// Move it to c.essential and remove it and all Elements it covers.
		s := c.m.AdjToB(e)[0]
		for _, ee := range c.m.AdjToA(s) {
			c.m.RemoveB(ee)
		}
		c.essential[s] = struct{}{}
		c.m.RemoveA(s)
	}
	return ok
}
