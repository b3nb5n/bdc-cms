package resources

import "shared"

type PropertyData struct {
	Hosts []string `bson:"hosts" json:"hosts"`
}

type Property shared.Resource[PropertyData]