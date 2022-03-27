package resources

import "shared"

type CollectionData struct {
	Property shared.Snowflake `bson:"property" json:"property"`
	Schema []map[string]any `bson:"schema" json:"schema"`
}

type Collection shared.Resource[CollectionData]