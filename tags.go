package metrics

import (
	"sort"
	"strings"
)

// Tag represents a single tag (or dimension) of a measurement-series. It is recommended that you
// use the abbreviated struct literal form to construct tag entries, such as:
//     Tag{"stack","prod"}
// but you can always specify the field names for clarity:
//     Tag{Key:"stack", Value:"prod"}
type Tag struct {
	Key, Value string
}

// T is a helper function to create a Tag
func T(key, value string) Tag {
	return Tag{
		Key:   key,
		Value: value,
	}
}

type TagSet []Tag

func (t *TagSet) Len() int {
	return len(*t)
}

func (t *TagSet) Less(i, j int) bool {
	return (*t)[i].Key < (*t)[j].Key
}

func (t *TagSet) Swap(i, j int) {
	(*t)[i], (*t)[j] = (*t)[j], (*t)[i]
}

// JoinTags joins individual tags into a TagSet
func JoinTags(tags ...Tag) TagSet {
	var result TagSet
	result = append(result, tags...)
	return result
}

func (t *TagSet) With(tag Tag) TagSet {
	*t = append(*t, tag)
	return *t
}

// Merge builds a new TagSet populated with the receiver's tags and then the given tags.
// The receiver can be nil.
func (t TagSet) Merge(more TagSet) TagSet {
	merged := t[:]
	return append(merged, more...)
}

func (t TagSet) HashKey() string {
	if len(t) > 0 {
		sort.Sort(&t)

		n := 0
		for _, entry := range t {
			n += len(entry.Key)
			n += len(entry.Value)
			n += 1 // for ='s
		}
		n += len(t) - 1 // for ,'s

		var builder strings.Builder
		builder.Grow(n)

		for i, entry := range t {
			if i > 0 {
				builder.WriteString(",")
			}
			builder.WriteString(entry.Key)
			builder.WriteString("=")
			builder.WriteString(entry.Value)
		}

		return builder.String()
	} else {
		return ""
	}
}
