package postgres

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/viafintech/gotidus"
)

// TextAnonymizer is a gotidus.Anonymizer interface implementation.
// It can be used to randomize the text of a column.
type TextAnonymizer struct {
	buildMapping func(string) string
}

// NewTextAnonymizer initializes a new TextAnonymizer object
func NewTextAnonymizer() *TextAnonymizer {
	return &TextAnonymizer{
		buildMapping: buildMapping,
	}
}

func (a *TextAnonymizer) base() string {
	return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZßüäöÜÄÖ0123456789"
}

func (a *TextAnonymizer) mapping() string {
	return a.buildMapping(a.base())
}

// Build returns the partial query to randomize the text of a column while still keeping
// the general structure.
func (a *TextAnonymizer) Build(tableName, columnName string) string {
	return fmt.Sprintf(
		`translate(%s::TEXT, '%s'::TEXT, '%s'::TEXT)`,
		gotidus.FullColumnName(tableName, columnName),
		a.base(),
		a.mapping(),
	)
}

func buildMapping(base string) string {
	list := strings.Split(base, "")
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(list), func(i, j int) { list[i], list[j] = list[j], list[i] })

	return strings.Join(list, "")
}
