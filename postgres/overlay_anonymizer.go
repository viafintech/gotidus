package postgres

import (
	"fmt"
	"strings"

	"github.com/Barzahlen/gotidus"
)

// OverlayAnonymizer is a gotidus.Anonymizer implementation
// which allows to overlay parts of a text field with a string based on the
// given overlayBase.
// count defines the number of repetitions of the overlayBase to be overlayed.
// start defines the first character to be overlayed.
type OverlayAnonymizer struct {
	overlayBase string
	start       int
	count       int
}

// NewOverlayAnonymizer intiializes a new OverlayAnonymizer object
func NewOverlayAnonymizer(overlayBase string, start int, count int) *OverlayAnonymizer {
	return &OverlayAnonymizer{
		overlayBase: overlayBase,
		start:       start,
		count:       count,
	}
}

// Build returns the partiql query to overlay the column content with a string.
func (a *OverlayAnonymizer) Build(tableName, columnName string) string {
	overlay := strings.Repeat(a.overlayBase, a.count)

	return fmt.Sprintf(
		`"overlay"((%s)::text, '%s'::text, %d)`,
		gotidus.FullColumnName(tableName, columnName),
		overlay,
		a.start,
	)
}
