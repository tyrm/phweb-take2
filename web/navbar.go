package web

import "regexp"

func makeNavbar(path string) (navbar *[]templateNavbarNode) {
	newNavbar := []templateNavbarNode{
			{
				Text:     "Home",
				MatchStr: "^/home$",
				FAIcon:   "home",
				URL:      "/home",
			},
			{
				Text:     "Files",
				MatchStr: "^/web/files/.*$",
				FAIcon:   "file",
				URL:      "/web/files/",
			},
			{
				Text:   "Admin",
				FAIcon: "hammer",
				URL:    "#",
				Children: []*templateNavbarNode{
					{
						Text:     "Job Runner",
						MatchStr: "^/web/admin/jobrunner/.*$",
						FAIcon:   "clock",
						URL:      "/web/admin/jobrunner/",
					},
					{
						Text:     "Oauth Clients",
						MatchStr: "^/web/admin/oauth-clients/.*$",
						FAIcon:   "desktop",
						URL:      "/web/admin/oauth-clients/",
						Disabled: true,
					},
					{
						Text:     "Registry",
						MatchStr: "^/web/admin/registry/.*$",
						FAIcon:   "book",
						URL:      "/web/admin/registry/",
					},
					{
						Text:     "Users",
						MatchStr: "^/web/admin/users/.*$",
						FAIcon:   "user",
						URL:      "/web/admin/users/",
					},
					{
						Text:     "Something else here",
						FAIcon:   "paw",
						URL:      "#",
						Disabled: true,
					},
				},
			},
		}

	for i := 0; i < len(newNavbar); i++ {
		if newNavbar[i].MatchStr != "" {
			match, err := regexp.MatchString(newNavbar[i].MatchStr, path)
			if err != nil {
				logger.Errorf("makeNavbar:Error matching regex: %v", err)
			}
			if match {
				newNavbar[i].Active = true
			}

		}

		if newNavbar[i].Children != nil {
			for j := 0; j < len(newNavbar[i].Children); j++ {

				if newNavbar[i].Children[j].MatchStr != "" {
					subMatch, err := regexp.MatchString(newNavbar[i].Children[j].MatchStr, path)
					if err != nil {
						logger.Errorf("makeNavbar:Error matching regex: %v", err)
					}

					if subMatch {
						newNavbar[i].Active = true
						newNavbar[i].Children[j].Active = true
					}

				}

			}
		}
	}

	return &newNavbar
}