package service

import (
	"context"

	mclient "github.com/micro/go-micro/v2/client"
	olog "github.com/owncloud/ocis-pkg/v2/log"
	settings "github.com/owncloud/ocis-settings/pkg/proto/v0"
	ssvc "github.com/owncloud/ocis-settings/pkg/service/v0"
)

const (
	// AccountManagementPermissionID is the hardcoded setting UUID for the account management permission
	AccountManagementPermissionID   string = "8e587774-d929-4215-910b-a317b1e80f73"
	// AccountManagementPermissionName is the hardcoded setting name for the account management permission
	AccountManagementPermissionName string = "account-management"
	// GroupManagementPermissionID is the hardcoded setting UUID for the group management permission
	GroupManagementPermissionID     string = "522adfbe-5908-45b4-b135-41979de73245"
	// GroupManagementPermissionName is the hardcoded setting name for the group management permission
	GroupManagementPermissionName   string = "group-management"
)

// RegisterPermissions registers permissions for account management and group management with the settings service.
func RegisterPermissions(l *olog.Logger) {
	// TODO this won't work with a registry other than mdns. Look into Micro's client initialization.
	// https://github.com/owncloud/ocis-proxy/issues/38
	service := settings.NewBundleService("com.owncloud.api.settings", mclient.DefaultClient)

	permissionRequests := generateAccountManagementPermissionsRequests()
	for i := range permissionRequests {
		res, err := service.AddSettingToBundle(context.Background(), &permissionRequests[i])
		bundleID := permissionRequests[i].BundleId
		if err != nil {
			l.Err(err).Str("bundle", bundleID).Str("setting", permissionRequests[i].Setting.Id).Msg("error adding setting to bundle")
		} else {
			l.Info().Str("bundle", bundleID).Str("setting", res.Setting.Id).Msg("successfully added setting to bundle")
		}
	}
}

func generateAccountManagementPermissionsRequests() []settings.AddSettingToBundleRequest {
	return []settings.AddSettingToBundleRequest{
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          AccountManagementPermissionID,
				Name:        AccountManagementPermissionName,
				DisplayName: "Account Management",
				Description: "This permission gives full access to everything that is related to account management.",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_USER,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READWRITE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
		{
			BundleId: ssvc.BundleUUIDRoleAdmin,
			Setting: &settings.Setting{
				Id:          GroupManagementPermissionID,
				Name:        GroupManagementPermissionName,
				DisplayName: "Group Management",
				Description: "This permission gives full access to everything that is related to group management.",
				Resource: &settings.Resource{
					Type: settings.Resource_TYPE_GROUP,
					Id:   "all",
				},
				Value: &settings.Setting_PermissionValue{
					PermissionValue: &settings.Permission{
						Operation:  settings.Permission_OPERATION_READWRITE,
						Constraint: settings.Permission_CONSTRAINT_ALL,
					},
				},
			},
		},
	}
}
