package postgres

import (
	"fmt"
	"strings"

	"github.com/viafintech/gotidus"
)

// RemoveJSONKeysAnonymizer is a gotidus.Anonymizer interface implementation which allows
// removing specific keys from a JSON object.
type RemoveJSONKeysAnonymizer struct {
	keys []string
}

// NewRemoveJSONKeysAnonymizer initializes a new RemoveJSONKeysAnonymizer object.
// The slice of strings given is the keys that are removed from the JSON
// content if they exist.
func NewRemoveJSONKeysAnonymizer(keys []string) *RemoveJSONKeysAnonymizer {
	return &RemoveJSONKeysAnonymizer{
		keys: keys,
	}
}

// Build returns the partial query to remove specific keys from the given column.
func (a *RemoveJSONKeysAnonymizer) Build(tableName, columnName string) string {
	removedKeys := make([]string, len(a.keys))
	for i, key := range a.keys {
		removedKeys[i] = fmt.Sprintf("key <> '%s'", key)
	}

	removedKeysString := strings.Join(removedKeys, " AND ")

	return fmt.Sprintf(
		`(
      SELECT
        concat('{', string_agg(to_json("key") || ':' || "value", ','), '}')::JSON
      FROM json_each(%[1]s::JSON) WHERE %[2]s
    )`,
		gotidus.FullColumnName(tableName, columnName),
		removedKeysString,
	)
}
