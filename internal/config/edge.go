package config

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// EdgeSetting is a set of board edges,
// or directions out from the cursor.
type EdgeSetting uint8

// Edges a setting can hold, alone or combined.
const (
	EdgeNone EdgeSetting = 0

	EdgeLeft   EdgeSetting = 1 << 0
	EdgeRight  EdgeSetting = 1 << 1
	EdgeTop    EdgeSetting = 1 << 2
	EdgeBottom EdgeSetting = 1 << 3

	EdgeHorizontal = EdgeLeft | EdgeRight
	EdgeVertical   = EdgeTop | EdgeBottom
	EdgeAll        = EdgeHorizontal | EdgeVertical
)

// Has returns whether the setting covers any of the given edges.
func (edge EdgeSetting) Has(edges EdgeSetting) bool {
	return edge&edges != 0
}

// edgeNames names an edge or group of edges
// supported in the config file.
var edgeNames = map[string]EdgeSetting{
	"":           EdgeNone,
	"none":       EdgeNone,
	"left":       EdgeLeft,
	"right":      EdgeRight,
	"top":        EdgeTop,
	"bottom":     EdgeBottom,
	"horizontal": EdgeHorizontal,
	"vertical":   EdgeVertical,
	"all":        EdgeAll,
}

// UnmarshalYAML reads an edge name, or a list of them, from the config.
func (edge *EdgeSetting) UnmarshalYAML(node *yaml.Node) error {
	var names []string

	switch node.Kind {
	case yaml.ScalarNode:
		var name string
		if err := node.Decode(&name); err != nil {
			return err
		}

		names = []string{name}

	case yaml.SequenceNode:
		if err := node.Decode(&names); err != nil {
			return err
		}

	default:
		return fmt.Errorf("line %d: edges must be a name or a list of names", node.Line)
	}

	var edges EdgeSetting

	for _, name := range names {
		edge, ok := edgeNames[strings.ToLower(strings.TrimSpace(name))]
		if !ok {
			return fmt.Errorf("line %d: invalid edge %q", node.Line, name)
		}

		edges |= edge
	}

	*edge = edges

	return nil
}
