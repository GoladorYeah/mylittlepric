package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("email").
			Unique().
			NotEmpty(),
		field.String("password_hash").
			Optional().
			Sensitive(),
		field.String("google_id").
			Optional().
			Unique(),
		field.String("name").
			Optional(),
		field.String("avatar_url").
			Optional(),
		field.String("provider").
			Default("email"), // "email" or "google"
		field.Time("created_at").
			Immutable().
			Default(func() time.Time { return time.Now() }),
		field.Time("updated_at").
			Default(func() time.Time { return time.Now() }).
			UpdateDefault(func() time.Time { return time.Now() }),
		field.Time("last_login").
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sessions", ChatSession.Type),
		edge.To("search_history", SearchHistory.Type),
		edge.To("preferences", UserPreference.Type).
			Unique(), // One-to-one relationship
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return []ent.Index{
		// Index for email lookups during login
		// Email is already unique, but this improves lookup performance
		index.Fields("email"),
		// Index for OAuth provider lookups (Google login)
		index.Fields("provider", "google_id"),
	}
}
