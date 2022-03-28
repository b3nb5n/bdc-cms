package shared

import (
	"encoding/json"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Visibility int

const (
	Live Visibility = iota
	Draft
	Archived
	Deleted
)

type ResourceMeta struct {
	Created time.Time `bson:"created" json:"created"`
	Edited time.Time `bson:"edited" json:"edited"`
	Visibility Visibility `bson:"visibility" json:"visibility"`
}

func NewResourceMeta() *ResourceMeta {
	now := time.Now()
	return &ResourceMeta{
		Created: now,
		Edited: now,
	}
}

func (meta *ResourceMeta) MarshalJSON() ([]byte, error) {
	type Alias ResourceMeta
	return json.Marshal(&struct {
		*Alias
		Created int64 `json:"created"`
		Edited int64 `json:"edited"`
	} {
		Alias: (*Alias)(meta),
		Created: meta.Created.Unix(),
		Edited: meta.Edited.Unix(),
	})
}

func (meta *ResourceMeta) UnmarshalJSON(data []byte) error {
	type Alias ResourceMeta
	aux := &struct {
		*Alias
		Created int64 `json:"created"`
		Edited int64 `json:"edited"`
	} {
		Alias: (*Alias)(meta),
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	meta.Created = time.Unix(aux.Created, 0)
	meta.Edited = time.Unix(aux.Edited, 0)
	return nil
}

func (meta *ResourceMeta) MarshalBSON() ([]byte, error) {
	type Alias ResourceMeta
	return bson.Marshal(&struct {
		*Alias
		Created primitive.DateTime `bson:"created"`
		Edited primitive.DateTime `bson:"edited"`
	} {
		Alias: (*Alias)(meta),
		Created: primitive.NewDateTimeFromTime(meta.Created),
		Edited: primitive.NewDateTimeFromTime(meta.Edited),
	})
}

func (meta *ResourceMeta) UnmarshalBSON(data []byte) error {
	type Alias ResourceMeta
	aux := &struct {
		*Alias
		Created primitive.DateTime `bson:"created"`
		Edited primitive.DateTime `bson:"edited"`
	} {
		Alias: (*Alias)(meta),
	}

	err := bson.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	meta.Created = aux.Created.Time()
	meta.Edited = aux.Edited.Time()
	return nil
}

type Resource[T any] struct {
	ID Snowflake `bson:"_id" json:"id"`
	Meta ResourceMeta `bson:"meta" json:"meta"`
	Data T `bson:"data" json:"data"`
}

func NewResource[T any](data T) (*Resource[T], error) {
	snowflake, err := NewSnowflake()
	if err != nil {
		return nil, fmt.Errorf("Error generating snowflake: %v", err)
	}

	return &Resource[T]{
		ID: snowflake,
		Meta: *NewResourceMeta(),
		Data: data,
	}, err
}