package db

//ZbotDatabase DB interface for Zbot
type ZbotDatabase interface {
	GetConnectionInfo() string
	Init() error
	Close()
	Statistics(string) (string, error)
	Append(Definition) error
	Top() ([]Definition, error)
	Rand() (Definition, error)
	Last() (Definition, error)
	Get(string, string) (Definition, error)
	Set(Definition) (string, error)
	_set(string, Definition) error
	Find(string, string) ([]Definition, error)
	Search(string, string) ([]Definition, error)
	Forget(Definition) error
	UserLevel(string) (string, error)
	UserIgnoreInsert(string) error
	//UserCheckIgnore return true if the user is on the ignore_list, false if it isn´t
	UserCheckIgnore(string) bool
	UserCleanupIgnorelist() error
	UserIgnoreList() ([]UserIgnore, error)

	Lock(Definition) error
}

//Definition .
type Definition struct {
	Term    string
	Meaning string
	Author  string
	Date    string
	Chat    string
	ID      int
}

//UserIgnore .
type UserIgnore struct {
	Username string
	Since    string
	Until    string
}
