// Code generated from sort.go using genzfunc.go; DO NOT EDIT.

// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices

func (cmp compare[E]) isSorted(list []E) bool {
	for i := len(list) - 1; i > 0; i-- {
		if cmp.Less(list[i], list[i-1]) {
			return false
		}
	}
	return true
}

func (cmp compare[E]) sortFast(list []E) {
	size := len(list)
	chance := log2Ceil(uint(size)) * 3 / 2
	if size > 50 {
		a, b, c := size/4, size/2, size*3/4
		a, ha := cmp.median(list, a-1, a, a+1)
		b, hb := cmp.median(list, b-1, b, b+1)
		c, hc := cmp.median(list, c-1, c, c+1)
		m, hint := cmp.median(list, a, b, c)
		hint &= ha & hb & hc

		pivot := list[m]
		if hint == hintRevered {
			Reverse(list)
			hint = hintSorted
		}
		if hint == hintSorted && cmp.isSorted(list) {
			return
		}

		l, r := 0, size-1
		for {
			for cmp.Less(list[l], pivot) {
				l++
			}
			for cmp.Less(pivot, list[r]) {
				r--
			}
			if l >= r {
				break
			}
			list[l], list[r] = list[r], list[l]
			l++
			r--
		}

		if l > size/2 {
			cmp.introSort(list[l:], chance)
			list = list[:l]
		} else {
			cmp.introSort(list[:l], chance)
			list = list[l:]
		}
	}
	cmp.introSort(list, chance)
}

func (cmp compare[E]) median(list []E, a, b, c int) (int, uint8) {

	if cmp.Less(list[b], list[a]) {
		if cmp.Less(list[c], list[b]) {
			return b, hintRevered
		} else if cmp.Less(list[c], list[a]) {
			return c, 0
		} else {
			return a, 0
		}
	} else {
		if cmp.Less(list[c], list[a]) {
			return a, 0
		} else if cmp.Less(list[c], list[b]) {
			return c, 0
		} else {
			return b, hintSorted
		}
	}
}

func (cmp compare[E]) sortStable(list []E) {
	if size := len(list); size < 16 {
		cmp.simpleSort(list)
	} else {
		step := 8
		a, b := 0, step
		for b <= size {
			cmp.simpleSort(list[a:b])
			a = b
			b += step
		}
		cmp.simpleSort(list[a:])

		for step < size {
			a, b = 0, step*2
			for b <= size {
				cmp.symmerge(list[a:b], step)
				a = b
				b += step * 2
			}
			if a+step < size {
				cmp.symmerge(list[a:], step)
			}
			step *= 2
		}
	}
}

func (cmp compare[E]) simpleSort(list []E) {
	if len(list) < 2 {
		return
	}
	for i := 1; i < len(list); i++ {
		curr := list[i]
		if cmp.Less(curr, list[0]) {
			for j := i; j > 0; j-- {
				list[j] = list[j-1]
			}
			list[0] = curr
		} else {
			pos := i
			for ; cmp.Less(curr, list[pos-1]); pos-- {
				list[pos] = list[pos-1]
			}
			list[pos] = curr
		}
	}
}

func (cmp compare[E]) heapSort(list []E) {
	for idx := len(list)/2 - 1; idx >= 0; idx-- {
		cmp.heapDown(list, idx)
	}
	for end := len(list) - 1; end > 0; end-- {
		list[0], list[end] = list[end], list[0]
		cmp.heapDown(list[:end], 0)
	}
}

func (cmp compare[E]) heapDown(list []E, pos int) {
	curr := list[pos]
	kid, last := pos*2+1, len(list)-1
	for kid < last {
		if cmp.Less(list[kid], list[kid+1]) {
			kid++
		}
		if !cmp.Less(curr, list[kid]) {
			break
		}
		list[pos] = list[kid]
		pos, kid = kid, kid*2+1
	}
	if kid == last && cmp.Less(curr, list[kid]) {
		list[pos], pos = list[kid], kid
	}
	list[pos] = curr
}

func (cmp compare[E]) sortIndex5(list []E,
	a, b, c, d, e int) (int, int, int, int, int) {
	if cmp.Less(list[b], list[a]) {
		a, b = b, a
	}
	if cmp.Less(list[d], list[c]) {
		c, d = d, c
	}
	if cmp.Less(list[c], list[a]) {
		a, c = c, a
		b, d = d, b
	}
	if cmp.Less(list[c], list[e]) {
		if cmp.Less(list[d], list[e]) {
			if cmp.Less(list[b], list[d]) {
				if cmp.Less(list[c], list[b]) {
					return a, c, b, d, e
				} else {
					return a, b, c, d, e
				}
			} else if cmp.Less(list[b], list[e]) {
				return a, c, d, b, e
			} else {
				return a, c, d, e, b
			}
		} else {
			if cmp.Less(list[b], list[e]) {
				if cmp.Less(list[c], list[b]) {
					return a, c, b, e, d
				} else {
					return a, b, c, e, d
				}
			} else if cmp.Less(list[b], list[d]) {
				return a, c, e, b, d
			} else {
				return a, c, e, d, b
			}
		}
	} else {
		if cmp.Less(list[b], list[c]) {
			if cmp.Less(list[e], list[a]) {
				return e, a, b, c, d
			} else if cmp.Less(list[e], list[b]) {
				return a, e, b, c, d
			} else {
				return a, b, e, c, d
			}
		} else {
			if cmp.Less(list[a], list[e]) {
				a, e = e, a
			}
			if cmp.Less(list[d], list[b]) {
				b, d = d, b
			}
			return e, a, c, b, d
		}
	}
}

func (cmp compare[E]) triPartition(list []E) (l, r int) {
	size := len(list)
	m, s := size/2, size/4

	x, l, _, r, y := cmp.sortIndex5(list, m-s, m-1, m, m+1, m+s)

	s = size - 1
	pivotL, pivotR := list[l], list[r]
	list[l], list[r] = list[0], list[s]
	list[1], list[x] = list[x], list[1]
	list[s-1], list[y] = list[y], list[s-1]

	l, r = 2, s-2
	for {
		for cmp.Less(list[l], pivotL) {
			l++
		}
		for cmp.Less(pivotR, list[r]) {
			r--
		}
		if cmp.Less(pivotR, list[l]) {
			list[l], list[r] = list[r], list[l]
			r--
			if cmp.Less(list[l], pivotL) {
				l++
				continue
			}
		}
		break
	}

	for k := l + 1; k <= r; k++ {
		if cmp.Less(pivotR, list[k]) {
			for cmp.Less(pivotR, list[r]) {
				r--
			}
			if k >= r {
				break
			}
			if cmp.Less(list[r], pivotL) {
				list[l], list[k], list[r] = list[r], list[l], list[k]
				l++
			} else {
				list[k], list[r] = list[r], list[k]
			}
			r--
		} else if cmp.Less(list[k], pivotL) {
			list[k], list[l] = list[l], list[k]
			l++
		}
	}

	l--
	r++
	list[0], list[l] = list[l], pivotL
	list[s], list[r] = list[r], pivotR
	return l, r
}

func (cmp compare[E]) introSort(list []E, chance int) {
	for len(list) > 14 {
		if chance--; chance < 0 {
			cmp.heapSort(list)
			return
		}

		l, r := cmp.triPartition(list)
		cmp.introSort(list[:l], chance)
		cmp.introSort(list[r+1:], chance)
		if !cmp.Less(list[l], list[r]) {
			return
		}
		list = list[l+1 : r]
	}
	cmp.simpleSort(list)
}

func (cmp compare[E]) symmerge(list []E, border int) {
	size := len(list)

	if border == 1 {
		curr := list[0]
		a, b := 1, size
		for a < b {
			m := int(uint(a+b) / 2)
			if cmp.Less(list[m], curr) {
				a = m + 1
			} else {
				b = m
			}
		}
		for i := 1; i < a; i++ {
			list[i-1] = list[i]
		}
		list[a-1] = curr
		return
	}

	if border == size-1 {
		curr := list[border]
		a, b := 0, border
		for a < b {
			m := int(uint(a+b) / 2)
			if cmp.Less(curr, list[m]) {
				b = m
			} else {
				a = m + 1
			}
		}
		for i := border; i > a; i-- {
			list[i] = list[i-1]
		}
		list[a] = curr
		return
	}

	half := size / 2
	n := border + half
	a, b := 0, border
	if border > half {
		a, b = n-size, half
	}

	p := n - 1
	for a < b {
		m := int(uint(a+b) / 2)
		if cmp.Less(list[p-m], list[m]) {
			b = m
		} else {
			a = m + 1
		}
	}
	b = n - a

	if a < border && border < b {
		rotateLeft(list[a:b], border-a)
	}
	if 0 < a && a < half {
		cmp.symmerge(list[:half], a)
	}
	if half < b && b < size {
		cmp.symmerge(list[half:], b-half)
	}
}
