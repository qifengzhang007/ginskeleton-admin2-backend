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
		if value.HasSubNode > 0 {
			value.Children = []auth.AuthOrganizationPostTree{}
			data[key] = value
		}
	}
	return
}
