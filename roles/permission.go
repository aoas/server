package roles

type PermissionType int

const (
	All PermissionType = iota + 1
	Create
	Read
	Update
	Delete
)

func (p *PermissionType) String() string {
	switch p {
	case Read:
		return "Read"
	case Create:
		return "Create"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case All:
		return "All"
	default:
		"Unknown"
	}
}
