package acl

import (
	"sort"

	"github.com/uhppoted/uhppote-core/types"
)

type equivalent = func(types.Card, types.Card) bool

func Compare(src, dst ACL) (map[uint32]Diff, error) {
	m := map[uint32]Diff{}

	for k := range src {
		m[k] = Diff{}
	}

	for k := range dst {
		m[k] = Diff{}
	}

	for k := range m {
		p := src[k]
		q := dst[k]
		m[k] = compare(k, p, q, equals)
	}

	return m, nil
}

func CompareWithPIN(src, dst ACL) (map[uint32]Diff, error) {
	m := map[uint32]Diff{}

	for k := range src {
		m[k] = Diff{}
	}

	for k := range dst {
		m[k] = Diff{}
	}

	for k := range m {
		p := src[k]
		q := dst[k]
		m[k] = compare(k, p, q, equalsWithPIN)
	}

	return m, nil
}

func compare(device uint32, p, q map[uint32]types.Card, eq equivalent) Diff {
	cards := map[uint32]struct{}{}

	for k := range p {
		cards[k] = struct{}{}
	}

	for k := range q {
		cards[k] = struct{}{}
	}

	diff := Diff{
		Unchanged: []types.Card{},
		Added:     []types.Card{},
		Updated:   []types.Card{},
		Deleted:   []types.Card{},
	}

	for k := range cards {
		u, hasu := p[k]
		v, hasv := q[k]

		if hasu && hasv {
			if eq(u, v) {
				diff.Unchanged = append(diff.Unchanged, u)
			} else {
				diff.Updated = append(diff.Updated, v)
			}
		} else if !hasu && hasv {
			diff.Added = append(diff.Added, v)
		} else if hasu && !hasv {
			diff.Deleted = append(diff.Deleted, u)
		}
	}

	for _, list := range [][]types.Card{
		diff.Unchanged,
		diff.Added,
		diff.Updated,
		diff.Deleted} {
		sort.Slice(list, func(i, j int) bool { return list[i].CardNumber < list[j].CardNumber })
	}

	return diff
}

/*
 * Compares two cards, ignoring PIN
 */
func equals(p, q types.Card) bool {
	if p.CardNumber != q.CardNumber {
		return false
	}

	if p.From != nil && q.From != nil {
		if !p.From.Equals(*q.From) {
			return false
		}
	} else if p.From != nil || q.From != nil {
		return false
	}

	if p.To != nil && q.To != nil {
		if !p.To.Equals(*q.To) {
			return false
		}
	} else if p.To != nil || q.To != nil {
		return false
	}

	for _, i := range []uint8{1, 2, 3, 4} {
		if p.Doors[i] != q.Doors[i] {
			return false
		}
	}

	return true
}

/*
 * Compares two cards, including PIN
 */
func equalsWithPIN(p, q types.Card) bool {
	if !equals(p, q) {
		return false
	} else if p.PIN != q.PIN {
		return false
	}

	return true
}
