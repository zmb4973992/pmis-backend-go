package util

import (
	"pmis-backend-go/global"
	"pmis-backend-go/model"
)

func GetOrganizationIDsForDataRange(userID int) (organizationIDsForDataRange []int) {
	organizationIDs := GetOrganizationIDs(userID)
	global.DB.Model(&model.OrganizationAndDataRange{}).Where("organization_id in ?", organizationIDs).
		Select("data_range_id").Find(&organizationIDsForDataRange)
	organizationIDsForDataRange = RemoveDuplication(organizationIDsForDataRange)
	return
}

// GetOrganizationIDs 获得所有的组织id(含子组织)
func GetOrganizationIDs(userID int) (organizationIDs []int) {
	global.DB.Model(&model.OrganizationAndUser{}).Where("user_id = ?", userID).
		Select("organization_id").Find(&organizationIDs)
	for i := range organizationIDs {
		res := getInferiorOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	organizationIDs = RemoveDuplication(organizationIDs)
	return
}

func getInferiorOrganizationIDs(organizationID int) (organizationIDs []int) {
	global.DB.Model(&model.Organization{}).Where("superior_id = ?", organizationID).
		Select("id").Find(&organizationIDs)
	for i := range organizationIDs {
		res := getInferiorOrganizationIDs(organizationIDs[i])
		organizationIDs = append(organizationIDs, res...)
	}
	return
}
