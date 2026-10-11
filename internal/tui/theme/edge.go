package theme

import "github.com/vinegm/termsweeper/internal/config"

// Glyphs marking a window edge the board continues past.
const (
	ArrowUp    = "▲"
	ArrowDown  = "▼"
	ArrowLeft  = "◀"
	ArrowRight = "▶"
)

type Edge = config.EdgeSetting

const (
	EdgeLeft   = config.EdgeLeft
	EdgeRight  = config.EdgeRight
	EdgeTop    = config.EdgeTop
	EdgeBottom = config.EdgeBottom
)
