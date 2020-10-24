package mod

type Dialog struct {
	Show        bool
	Green       func()
	GreenLabel  string
	Orange      func()
	OrangeLabel string
	Red         func()
	RedLabel    string
	CustomField func()
	Title       string
	Text        string
}