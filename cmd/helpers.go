package cmd

import "k8s.io/client-go/util/cert"

// Returns name that will be used as certificate name.
// returns common name is organization is not provided.
// return <common-name>@<organization-name> if organization is provided
func filename(cfg cert.Config) string {
	if len(cfg.Organization) == 0 {
		return cfg.CommonName
	}
	return cfg.CommonName + "@" + cfg.Organization[0]
}
