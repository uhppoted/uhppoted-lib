package acl

import (
	"reflect"
	"sort"

	"github.com/uhppoted/uhppote-core/types"
)

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
		m[k] = compare(k, p, q)
	}

	return m, nil
}

func compare(device uint32, p, q map[uint32]types.Card) Diff {
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
			if reflect.DeepEqual(u, v) {
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
