package social

type Permission int

const (
	PermissionNone Permission = iota
	PermissionRead
	PermissionWrite
	PermissionReadWrite
)

var permissions = []string{
	"", "r", "w", "rw",
}

func (p Permission) String() string {
	return permissions[p]
}

func (p Permission) MarshalText() ([]byte, error) {
	return []byte(p.String()), nil
}

func (p *Permission) UnmarshalText(text []byte) error {
	*p = PermissionRead
	enum := string(text)
	for i, k := range permissions {
		if enum == k {
			*p = Permission(i)
			return nil
		}
	}
	return nil
}

func PermissionFromString(text string) Permission {
	for i, k := range permissions {
		if text == k {
			return Permission(i)
		}
	}
	return PermissionNone
}
