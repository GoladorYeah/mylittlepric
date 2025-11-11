package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// SearchHistory holds the schema definition for the SearchHistory entity.
type SearchHistory struct {
	ent.Schema
}

// Fields of the SearchHistory.
func (SearchHistory) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Optional(),
		field.String("session_id").
			Optional(),
		field.String("search_query").
			NotEmpty(),
		field.String("optimized_query").
			Optional(),
		field.String("search_type").
			NotEmpty(),
		field.String("category").
			Optional(),
		field.String("country_code").
			Default("US"),
		field.String("language_code").
			Default("en"),
		field.String("currency").
			Default("USD"),
		field.Int("result_count").
			Default(0),
		field.JSON("products_found", []map[string]interface{}{}).
			Optional().
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}),
		field.String("clicked_product_id").
			Optional(),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("expires_at").
			Optional(),
	}
}

// Edges of the SearchHistory.
func (SearchHistory) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("search_history").
			Field("user_id").
			Unique(),
	}
}

// Indexes of the SearchHistory.
func (SearchHistory) Indexes() []ent.Index {
	return []ent.Index{
		// Index for GetUserSearchHistory - filtering by user and ordering by date
		index.Fields("user_id", "created_at"),
		// Index for anonymous users - filtering by session and ordering by date
		index.Fields("session_id", "created_at"),
		// Partial index for cleanup job - only for anonymous users (user_id IS NULL)
		// This optimizes the CleanupExpiredAnonymousHistory query
		index.Fields("expires_at").
			Annotations(entsql.IndexWhere("user_id IS NULL")),
	}
}
