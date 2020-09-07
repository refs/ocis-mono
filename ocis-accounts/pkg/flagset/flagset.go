package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/refs/ocis-mono/ocis-accounts/pkg/config"
	accounts "github.com/refs/ocis-mono/ocis-accounts/pkg/proto/v0"
)

// RootWithConfig applies cfg to the root flagset
func RootWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"ACCOUNTS_LOG_LEVEL"},
			Destination: &cfg.Log.Level,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-pretty",
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"ACCOUNTS_LOG_PRETTY"},
			Destination: &cfg.Log.Pretty,
		},
		&cli.BoolFlag{
			Value:       true,
			Name:        "log-color",
			Usage:       "Enable colored logging",
			EnvVars:     []string{"ACCOUNTS_LOG_COLOR"},
			Destination: &cfg.Log.Color,
		},
	}
}

// ServerWithConfig applies cfg to the root flagset
func ServerWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "http-namespace",
			Value:       "com.owncloud.web",
			Usage:       "Set the base namespace for the http namespace",
			EnvVars:     []string{"ACCOUNTS_HTTP_NAMESPACE"},
			Destination: &cfg.HTTP.Namespace,
		},
		&cli.StringFlag{
			Name:        "http-addr",
			Value:       "0.0.0.0:9181",
			Usage:       "Address to bind http server",
			EnvVars:     []string{"ACCOUNTS_HTTP_ADDR"},
			Destination: &cfg.HTTP.Addr,
		},
		&cli.StringFlag{
			Name:        "http-root",
			Value:       "/",
			Usage:       "Root path of http server",
			EnvVars:     []string{"ACCOUNTS_HTTP_ROOT"},
			Destination: &cfg.HTTP.Root,
		},
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "grpc-addr",
			Value:       "0.0.0.0:9180",
			Usage:       "Address to bind grpc server",
			EnvVars:     []string{"ACCOUNTS_GRPC_ADDR"},
			Destination: &cfg.GRPC.Addr,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.StringFlag{
			Name:        "accounts-data-path",
			Value:       "/var/tmp/ocis-accounts",
			Usage:       "accounts folder",
			EnvVars:     []string{"ACCOUNTS_DATA_PATH"},
			Destination: &cfg.Server.AccountsDataPath,
		},
		&cli.StringFlag{
			Name:        "asset-path",
			Value:       "",
			Usage:       "Path to custom assets",
			EnvVars:     []string{"ACCOUNTS_ASSET_PATH"},
			Destination: &cfg.Asset.Path,
		},
		&cli.StringFlag{
			Name:        "jwt-secret",
			Value:       "Pive-Fumkiu4",
			Usage:       "Used to create JWT to talk to reva, should equal reva's jwt-secret",
			EnvVars:     []string{"ACCOUNTS_JWT_SECRET"},
			Destination: &cfg.TokenManager.JWTSecret,
		},
	}
}

// UpdateAccountWithConfig applies update command flags to cfg
func UpdateAccountWithConfig(cfg *config.Config, a *accounts.Account) []cli.Flag {
	if a.PasswordProfile == nil {
		a.PasswordProfile = &accounts.PasswordProfile{}
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.BoolFlag{
			Name:        "enabled",
			Usage:       "Enable the account",
			Destination: &a.AccountEnabled,
		},
		&cli.StringFlag{
			Name:        "displayname",
			Usage:       "Set the displayname for the account",
			Destination: &a.DisplayName,
		},
		&cli.StringFlag{
			Name:        "preferred-name",
			Usage:       "Set the preferred-name for the account",
			Destination: &a.PreferredName,
		},
		&cli.StringFlag{
			Name:        "on-premises-sam-account-name",
			Usage:       "Set the on-premises-sam-account-name",
			Destination: &a.OnPremisesSamAccountName,
		},
		&cli.Int64Flag{
			Name:        "uidnumber",
			Usage:       "Set the uidnumber for the account",
			Destination: &a.UidNumber,
		},
		&cli.Int64Flag{
			Name:        "gidnumber",
			Usage:       "Set the gidnumber for the account",
			Destination: &a.GidNumber,
		},
		&cli.StringFlag{
			Name:        "mail",
			Usage:       "Set the mail for the account",
			Destination: &a.Mail,
		},
		&cli.StringFlag{
			Name:        "description",
			Usage:       "Set the description for the account",
			Destination: &a.Description,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "Set the password for the account",
			Destination: &a.PasswordProfile.Password,
			// TODO read password from ENV?
		},
		&cli.StringSliceFlag{
			Name:  "password-policies",
			Usage: "Possible policies: DisableStrongPassword, DisablePasswordExpiration",
		},
		&cli.BoolFlag{
			Name:        "force-password-change",
			Usage:       "Force password change on next sign-in",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignIn,
		},
		&cli.BoolFlag{
			Name:        "force-password-change-mfa",
			Usage:       "Force password change on next sign-in with mfa",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignInWithMfa,
		},
	}
}

// AddAccountWithConfig applies create command flags to cfg
func AddAccountWithConfig(cfg *config.Config, a *accounts.Account) []cli.Flag {
	if a.PasswordProfile == nil {
		a.PasswordProfile = &accounts.PasswordProfile{}
	}

	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
		&cli.BoolFlag{
			Name:        "enabled",
			Usage:       "Enable the account",
			Destination: &a.AccountEnabled,
		},
		&cli.StringFlag{
			Name:        "displayname",
			Usage:       "Set the displayname for the account",
			Destination: &a.DisplayName,
		},
		&cli.StringFlag{
			Name:  "username",
			Usage: "Username will be written to preferred-name and on_premises_sam_account_name",
		},
		&cli.StringFlag{
			Name:        "preferred-name",
			Usage:       "Set the preferred-name for the account",
			Destination: &a.PreferredName,
		},
		&cli.StringFlag{
			Name:        "on-premises-sam-account-name",
			Usage:       "Set the on-premises-sam-account-name",
			Destination: &a.OnPremisesSamAccountName,
		},
		&cli.Int64Flag{
			Name:        "uidnumber",
			Usage:       "Set the uidnumber for the account",
			Destination: &a.UidNumber,
		},
		&cli.Int64Flag{
			Name:        "gidnumber",
			Usage:       "Set the gidnumber for the account",
			Destination: &a.GidNumber,
		},
		&cli.StringFlag{
			Name:        "mail",
			Usage:       "Set the mail for the account",
			Destination: &a.Mail,
		},
		&cli.StringFlag{
			Name:        "description",
			Usage:       "Set the description for the account",
			Destination: &a.Description,
		},
		&cli.StringFlag{
			Name:        "password",
			Usage:       "Set the password for the account",
			Destination: &a.PasswordProfile.Password,
			// TODO read password from ENV?
		},
		&cli.StringSliceFlag{
			Name:  "password-policies",
			Usage: "Possible policies: DisableStrongPassword, DisablePasswordExpiration",
		},
		&cli.BoolFlag{
			Name:        "force-password-change",
			Usage:       "Force password change on next sign-in",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignIn,
		},
		&cli.BoolFlag{
			Name:        "force-password-change-mfa",
			Usage:       "Force password change on next sign-in with mfa",
			Destination: &a.PasswordProfile.ForceChangePasswordNextSignInWithMfa,
		},
	}
}

// ListAccountsWithConfig applies list command flags to cfg
func ListAccountsWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}

// RemoveAccountWithConfig applies remove command flags to cfg
func RemoveAccountWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}

// InspectAccountWithConfig applies inspect command flags to cfg
func InspectAccountWithConfig(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "grpc-namespace",
			Value:       "com.owncloud.api",
			Usage:       "Set the base namespace for the grpc namespace",
			EnvVars:     []string{"ACCOUNTS_GRPC_NAMESPACE"},
			Destination: &cfg.GRPC.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Value:       "accounts",
			Usage:       "service name",
			EnvVars:     []string{"ACCOUNTS_NAME"},
			Destination: &cfg.Server.Name,
		},
	}
}
