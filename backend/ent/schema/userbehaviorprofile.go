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

// UserBehaviorProfile holds the schema definition for the UserBehaviorProfile entity.
// This tracks long-term user preferences and behavior patterns for personalized recommendations.
type UserBehaviorProfile struct {
	ent.Schema
}

// Fields of the UserBehaviorProfile.
func (UserBehaviorProfile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Unique(), // One profile per user

		// Category preferences with weights (e.g., {"electronics": 0.8, "clothing": 0.6})
		field.JSON("category_preferences", map[string]float64{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default(map[string]float64{}),

		// Price ranges by category (e.g., {"electronics": {"min": 100, "max": 1000}})
		field.JSON("price_ranges", map[string]interface{}{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default(map[string]interface{}{}),

		// Preferred brands with frequency counts
		field.JSON("brand_preferences", map[string]int{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default(map[string]int{}),

		// Communication style preferences
		field.String("communication_style").
			Default("balanced"), // "brief", "balanced", "detailed"

		// Preferred search type ("product", "specification", "comparison")
		field.String("preferred_search_type").
			Optional(),

		// Average session length in minutes
		field.Float("avg_session_duration").
			Default(0),

		// Average messages per session
		field.Float("avg_messages_per_session").
			Default(0),

		// Success rate (sessions that led to product interaction)
		field.Float("success_rate").
			Default(0),

		// Common keywords extracted from user messages
		field.JSON("common_keywords", []string{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default([]string{}),

		// Time-based patterns (preferred shopping hours)
		field.JSON("time_patterns", map[string]interface{}{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Default(map[string]interface{}{}),

		// Total number of sessions analyzed
		field.Int("total_sessions").
			Default(0),

		// Total products viewed
		field.Int("total_products_viewed").
			Default(0),

		// Total products clicked
		field.Int("total_products_clicked").
			Default(0),

		// Last learning update timestamp
		field.Time("last_learning_update").
			Default(time.Now).
			UpdateDefault(time.Now),

		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the UserBehaviorProfile.
func (UserBehaviorProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("behavior_profile").
			Field("user_id").
			Required().
			Unique(),
	}
}

// Indexes of the UserBehaviorProfile.
func (UserBehaviorProfile) Indexes() []ent.Index {
	return []ent.Index{
		// Index for quick user profile lookup
		index.Fields("user_id").
			Unique(),
		// Index for finding profiles that need learning updates
		index.Fields("last_learning_update"),
	}
}
