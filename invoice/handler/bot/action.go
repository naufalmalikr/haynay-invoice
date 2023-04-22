package bot

type processFunc func(input string, user *user) error
type rollbackFunc func(user *user) error
type action struct {
	NextStep string
	Process  processFunc
	Rollback rollbackFunc
}
