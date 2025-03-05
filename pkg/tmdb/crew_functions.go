package tmdb

import "sort"

type BilledCrew struct {
	Name  string
	Roles []string
}

func (credits *Credits) BilledCrew() []BilledCrew {
	crewmap := make(map[string][]string)
	billedcrew := []BilledCrew{}
	for _, crew := range credits.Crew {
		if crew.Job == "Director" ||
			crew.Job == "Screenplay" ||
			crew.Job == "Writer" ||
			crew.Job == "Novel" ||
			crew.Job == "Story" {
			crewmap[crew.Name] = append(crewmap[crew.Name], crew.Job)
		}
	}

	for name, jobs := range crewmap {
		billedcrew = append(billedcrew, BilledCrew{Name: name, Roles: jobs})
	}
	for i := range billedcrew {
		sort.Strings(billedcrew[i].Roles)
	}
	sort.Slice(billedcrew, func(i, j int) bool {
		return billedcrew[i].Roles[0] < billedcrew[j].Roles[0]
	})
	return billedcrew
}

func (billedcrew *BilledCrew) FRoles() string {
	jobs := ""
	for _, job := range billedcrew.Roles {
		jobs += job + ", "
	}
	return jobs[:len(jobs)-2]
}
