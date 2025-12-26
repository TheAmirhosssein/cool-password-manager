package types

func NewNullString(str string) NullString {
	return NullString{String: str, Valid: str != ""}
}
