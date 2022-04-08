package resources

import "shared"

type ItemData struct {
	Collection shared.Snowflake `json:"collection" bson:"collection"`
	Data map[string]any
}

type Item = shared.Resource[ItemData]