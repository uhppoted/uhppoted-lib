package acl

import (
	"github.com/uhppoted/uhppote-core/types"
)

type SystemDiff map[uint32]Diff

type Diff struct {
	Unchanged []types.Card
	Updated   []types.Card
	Added     []types.Card
	Deleted   []types.Card
}

type ConsolidatedDiff struct {
	Unchanged []uint32 `json:"unchanged"`
	Updated   []uint32 `json:"updated"`
	Added     []uint32 `json:"added"`
	Deleted   []uint32 `json:"deleted"`
}

func (diff *SystemDiff) Consolidate() *ConsolidatedDiff {
	consolidated := map[uint32]*struct {
		updated bool
		added   bool
		deleted bool
	}{}

	for _, d := range *diff {
		lists := [][]types.Card{d.Unchanged, d.Updated, d.Added, d.Deleted}
		for _, l := range lists {
			for _, card := range l {
				consolidated[card.CardNumber] = &struct {
					updated bool
					added   bool
					deleted bool
				}{}
			}
		}
	}

	for _, d := range *diff {
		for _, card := range d.Updated {
			consolidated[card.CardNumber].updated = true
		}
	}

	for _, d := range *diff {
		for _, card := range d.Added {
			// A card that has been updated on one controller and added on another is regarded as 'updated on the system'
			if !consolidated[card.CardNumber].updated {
				consolidated[card.CardNumber].added = true
			}
		}
	}

	for _, d := range *diff {
		for _, card := range d.Deleted {
			consolidated[card.CardNumber].deleted = true
		}
	}

	unchanged := []uint32{}
	updated := []uint32{}
	added := []uint32{}
	deleted := []uint32{}

	for card, s := range consolidated {
		if !s.updated && !s.added && !s.deleted {
			unchanged = append(unchanged, card)
		}

		if s.updated {
			updated = append(updated, card)
		}

		if s.added {
			added = append(added, card)
		}

		if s.deleted {
			deleted = append(deleted, card)
		}
	}

	usort(unchanged)
	usort(updated)
	usort(added)
	usort(deleted)

	return &ConsolidatedDiff{
		Unchanged: unchanged,
		Updated:   updated,
		Added:     added,
		Deleted:   deleted,
	}
}

func (diff *SystemDiff) HasChanges() bool {
	for _, d := range *diff {
		if d.HasChanges() {
			return true
		}
	}

	return false
}

func (d *Diff) HasChanges() bool {
	if len(d.Updated) > 0 {
		return true
	}

	if len(d.Added) > 0 {
		return true
	}

	if len(d.Deleted) > 0 {
		return true
	}

	return false
}
