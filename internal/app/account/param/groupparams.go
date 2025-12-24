package param

import "github.com/TheAmirhosssein/cool-password-manage/internal/types"

type ReadGroupParams struct {
	MemberID types.ID
	Limit    int
	Offset   int
}
