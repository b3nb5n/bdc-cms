package resources

import "shared"

type PropertyData struct {
	Hosts []string `bson:"hosts" json:"hosts" validate:"required,dive,url"`
}

type Property shared.Resource[PropertyData]