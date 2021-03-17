package auth_organization_post

import (
	"goskeleton/app/model/auth"
)

type AuthOrganizationPostService struct {
}

func (a *AuthOrganizationPostService) GetOrgByFid(fid int) (err error, data []auth.AuthOrganizationPostTree) {
	models := auth.CreateAuthOrganizationFactory("")
	err = models.GetByFid(fid, &data)
	for key, value := range data {
		has := []auth.AuthOrganizationPostTree{}
		id := value.Id
		models.GetByFid(id, &has)
		if len(has) != 0 {
			value.Children = []auth.AuthOrganizationPostTree{}
			data[key] = value
		}
	}
	return
}
