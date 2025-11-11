package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// ChatSession holds the schema definition for the ChatSession entity.
type ChatSession struct {
	ent.Schema
}

// Fields of the ChatSession.
func (ChatSession) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.String("session_id").
			Unique().
			NotEmpty(),
		field.UUID("user_id", uuid.UUID{}).
			Optional(),
		field.String("country_code").
			Default("US"),
		field.String("language_code").
			Default("en"),
		field.String("currency").
			Default("USD"),
		field.Int("message_count").
			Default(0),
		// JSONB fields
		field.JSON("search_state", map[string]interface{}{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default(map[string]interface{}{}),
		field.JSON("cycle_state", map[string]interface{}{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default(map[string]interface{}{}),
		field.JSON("conversation_context", map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("expires_at").
			Default(func() time.Time {
				return time.Now().Add(24 * time.Hour) // Default 24h expiry
			}),
	}
}

// Edges of the ChatSession.
func (ChatSession) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sessions").
			Field("user_id").
			Unique(),
		edge.To("messages", Message.Type),
	}
}

// Indexes of the ChatSession.
func (ChatSession) Indexes() []ent.Index {
	return []ent.Index{
		// Index for GetActiveSessionForUser - filtering by user and checking expiry
		index.Fields("user_id", "expires_at"),
		// Index for cleanup operations - finding expired sessions
		index.Fields("expires_at"),
		// Index for session_id lookups (frequently used in ProcessChat)
		index.Fields("session_id"),
	}
}
