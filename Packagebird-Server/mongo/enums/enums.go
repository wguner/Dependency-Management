package enums

type Collection int

const (
	Packages Collection = iota
	Projects
	Graphs
	PackagesMetadata
	Users
	Authentications
	Sources
	Scripts
)

func (c Collection) String() string {
	return []string{
		"packages",
		"projects",
		"graphs",
		"packagesMetadata",
		"users",
		"authentications",
		"sources",
		"scripts",
	}[c]
}
