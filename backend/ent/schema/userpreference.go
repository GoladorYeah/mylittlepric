package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// UserPreference holds the schema definition for the UserPreference entity.
type UserPreference struct {
	ent.Schema
}

// Fields of the UserPreference.
func (UserPreference) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Unique(), // One-to-one with User
		field.String("country").
			Optional().
			Nillable(),
		field.String("currency").
			Optional().
			Nillable(),
		field.String("language").
			Optional().
			Nillable(),
		field.String("theme").
			Optional().
			Nillable(),
		field.Bool("sidebar_open").
			Optional().
			Nillable(),
		field.String("last_active_session_id").
			Optional().
			Nillable(),
		field.JSON("saved_search", map[string]interface{}{}).
			Optional(),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the UserPreference.
func (UserPreference) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("preferences").
			Field("user_id").
			Required().
			Unique(),
	}
}
