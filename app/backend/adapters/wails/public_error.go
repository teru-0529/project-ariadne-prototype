package wails

// INFO: Frontendに公開するエラー
type PublicError struct {
	Processed bool
	Message   string
}

// FUNCTION: PublicErrorのエラーメッセージを返す関数
func (e *PublicError) Error() string {
	return e.Message
}
