package resources

import "shared"

type CollectionData struct {
	Property shared.Snowflake `bson:"property" json:"property" validate:"required"`
	Schema map[string]any `bson:"schema" json:"schema" validate:"required"`
}

type Collection shared.Resource[CollectionData]