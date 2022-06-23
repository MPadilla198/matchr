package matchr

import "unicode/utf8"

// Levenshtein computes the Levenshtein distance between two
// strings. The returned value - distance - is the number of insertions,
// deletions, and substitutions it takes to transform one
// string (s1) into another (s2). Each step in the transformation "costs"
// one distance point.
func Levenshtein(s1 string, s2 string) (distance int) {
	// index by code point, not byte
	r1 := []rune(s1)
	r2 := []rune(s2)

	rows := len(r1) + 1
	cols := len(r2) + 1

	var d1 int
	var d2 int
	var d3 int
	var i int
	var j int
	dist := make([]int, rows*cols)

	for i = 0; i < rows; i++ {
		dist[i*cols] = i
	}

	for j = 0; j < cols; j++ {
		dist[j] = j
	}

	for j = 1; j < cols; j++ {
		for i = 1; i < rows; i++ {
			if r1[i-1] == r2[j-1] {
				dist[(i*cols)+j] = dist[((i-1)*cols)+(j-1)]
			} else {
				d1 = dist[((i-1)*cols)+j] + 1
				d2 = dist[(i*cols)+(j-1)] + 1
				d3 = dist[((i-1)*cols)+(j-1)] + 1

				dist[(i*cols)+j] = min(d1, min(d2, d3))
			}
		}
	}

	distance = dist[(cols*rows)-1]

	return
}

// Found at: https://en.wikibooks.org/wiki/Algorithm_Implementation/Strings/Levenshtein_distance#Go
// Also see: https://xlinux.nist.gov/dads/HTML/Levenshtein.html
// This version uses dynamic programming with time complexity of O(mn) where m and n are lengths of a and b,
// and the space complexity is n + 1 of integers plus some constant  space(i.e. O(n)).
func Levenshtein_DP(a,b string) int {
	f := make([]int, utf8.RuneCountInString(b)+1)

	for j := range f {
		f[j] = j
	}

	for _, ca := range a {
		j := 1
		fj1 := f[0] // fj1 is the value of f[j - 1] in last iteration
		f[0]++
		for _, cb := range b {
			mn := min(f[j]+1, f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, f[j] = f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	return f[len(f)-1]
}
