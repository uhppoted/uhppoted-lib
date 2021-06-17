package acl

import (
	"github.com/uhppoted/uhppote-core/types"
	"reflect"
	"sort"
)

func Compare(src, dst ACL) (map[uint32]Diff, error) {
	m := map[uint32]Diff{}

	for k, _ := range src {
		m[k] = Diff{}
	}

	for k, _ := range dst {
		m[k] = Diff{}
	}

	for k, _ := range m {
		p := src[k]
		q := dst[k]
		m[k] = compare(p, q)
	}

	return m, nil
}

func compare(p, q map[uint32]types.Card) Diff {
	cards := map[uint32]struct{}{}

	for k, _ := range p {
		cards[k] = struct{}{}
	}

	for k, _ := range q {
		cards[k] = struct{}{}
	}

	diff := Diff{
		Unchanged: []types.Card{},
		Added:     []types.Card{},
		Updated:   []types.Card{},
		Deleted:   []types.Card{},
	}

	for k, _ := range cards {
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
