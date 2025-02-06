package nav

type NavItem struct {
	name string
	href string
}

func getNavItems() []NavItem {
	return []NavItem{
		{
			name: "Movies",
			href: "/movies",
		},
	}
}
