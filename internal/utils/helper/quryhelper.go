package helper

import (
	"fmt"
	"strings"

	"github.com/TheAmirhosssein/cool-password-manage/internal/types"
)

func MakeSearchQuery(q types.NullString, fields []string) string {
	if !q.Valid || q.String == "" || len(fields) == 0 {
		return ""
	}

	value := strings.ReplaceAll(q.String, "'", "''")

	conditions := make([]string, 0, len(fields))
	for _, field := range fields {
		conditions = append(
			conditions,
			fmt.Sprintf("%s ILIKE '%%%s%%'", field, value),
		)
	}

	return "AND (" + strings.Join(conditions, " OR ") + ")"
}
