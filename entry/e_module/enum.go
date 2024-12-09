package e_module

type Module int

const (
	None Module = iota
	user
	workspace
)

func (m Module) String() string {
	return [...]string{"", "user", "workspace"}[m]
}
