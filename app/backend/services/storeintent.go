package services

// INFO: Taskの永続化操作の意図を表す型
type StoreIntent int

const (
	StoreDeposit StoreIntent = iota
	StoreReplace
)
