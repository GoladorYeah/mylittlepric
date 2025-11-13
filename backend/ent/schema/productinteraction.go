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

// ProductInteraction holds the schema definition for the ProductInteraction entity.
// This tracks user interactions with products for better recommendations.
type ProductInteraction struct {
	ent.Schema
}

// Fields of the ProductInteraction.
func (ProductInteraction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Immutable(),
		field.UUID("user_id", uuid.UUID{}).
			Optional(), // Optional for anonymous users
		field.String("session_id").
			NotEmpty(),

		// Product information
		field.String("product_id").
			NotEmpty(), // External product ID from SERP API
		field.String("product_name").
			NotEmpty(),
		field.Float("product_price").
			Optional(),
		field.String("product_currency").
			Optional(),
		field.String("product_category").
			Optional(),
		field.String("product_brand").
			Optional(),
		field.String("product_url").
			Optional(),

		// Product data snapshot (for later analysis)
		field.JSON("product_data", map[string]interface{}{}).
			SchemaType(map[string]string{
				dialect.Postgres: "jsonb",
			}).
			Optional(),

		// Interaction type: "viewed", "clicked", "compared", "dismissed"
		field.String("interaction_type").
			NotEmpty(),

		// Context: where in the conversation this happened
		field.Int("message_position").
			Default(0), // Which message number in the session

		// Engagement metrics
		field.Int("view_duration_seconds").
			Default(0), // How long the product was visible
		field.Int("click_count").
			Default(0), // Number of clicks on this product
		field.Bool("opened_details").
			Default(false),
		field.Bool("added_to_comparison").
			Default(false),

		// User feedback
		field.String("feedback").
			Optional(), // "too_expensive", "not_relevant", "perfect", etc.
		field.Float("implicit_score").
			Default(0), // 0-1 score based on interaction patterns

		// Search context that led to this product
		field.String("search_query").
			Optional(),
		field.String("search_type").
			Optional(),

		// Position in results
		field.Int("position_in_results").
			Default(0),

		// Interaction sequence (to understand browsing patterns)
		field.Int("interaction_sequence").
			Default(0), // Order of this interaction in the session

		// Timestamp of interaction
		field.Time("interacted_at").
			Default(time.Now),

		field.Time("created_at").
			Immutable().
			Default(time.Now),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
	}
}

// Edges of the ProductInteraction.
func (ProductInteraction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("product_interactions").
			Field("user_id").
			Unique(),
	}
}

// Indexes of the ProductInteraction.
func (ProductInteraction) Indexes() []ent.Index {
	return []ent.Index{
		// Index for user interaction history
		index.Fields("user_id", "interacted_at"),
		// Index for session-based queries
		index.Fields("session_id", "interaction_sequence"),
		// Index for product popularity analysis
		index.Fields("product_id", "interaction_type", "interacted_at"),
		// Index for category-based recommendations
		index.Fields("user_id", "product_category", "implicit_score"),
		// Composite index for user preferences by brand
		index.Fields("user_id", "product_brand"),
	}
}
