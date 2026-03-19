package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type OAuthClient struct {
	ent.Schema
}

func (OAuthClient) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().Immutable(), // public client_id
		field.String("name").NotEmpty(),
		field.String("secret").Sensitive(), // hashed client_secret
		field.String("redirect_uri").NotEmpty(),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}