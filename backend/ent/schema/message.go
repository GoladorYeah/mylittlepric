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

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("session_id", uuid.UUID{}),
		field.String("role").
			NotEmpty(), // "user" or "assistant"
		field.Text("content").
			NotEmpty(),
		field.String("response_type").
			Optional(),
		field.Strings("quick_replies").
			Optional(),
		field.JSON("products", []map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}),
		field.JSON("search_info", map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("session", ChatSession.Type).
			Ref("messages").
			Field("session_id").
			Required().
			Unique(),
	}
}

// Indexes of the Message.
func (Message) Indexes() []ent.Index {
	return []ent.Index{
		// Index for GetMessagesFromDB - filtering by session_id and ordering by created_at
		index.Fields("session_id", "created_at"),
	}
}
