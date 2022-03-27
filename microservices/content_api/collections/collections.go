package collections

import "shared"

type CollectionData struct {
	Property shared.Snowflake
	Schema map[string]any
}

type Collection shared.Resource[CollectionData]

const COLLECTION = "collections"