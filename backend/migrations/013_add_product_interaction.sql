-- migrations/013_add_product_interaction.sql
-- Track user interactions with products for better recommendations

-- Product interaction table
CREATE TABLE IF NOT EXISTS product_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255) NOT NULL,

    -- Product information
    product_id VARCHAR(255) NOT NULL,
    product_name VARCHAR(500) NOT NULL,
    product_price DOUBLE PRECISION,
    product_currency VARCHAR(3),
    product_category VARCHAR(100),
    product_brand VARCHAR(255),
    product_url TEXT,

    -- Product data snapshot (JSONB for later analysis)
    product_data JSONB,

    -- Interaction type: "viewed", "clicked", "compared", "dismissed"
    interaction_type VARCHAR(50) NOT NULL,

    -- Context: where in the conversation this happened
    message_position INTEGER DEFAULT 0,

    -- Engagement metrics
    view_duration_seconds INTEGER DEFAULT 0,
    click_count INTEGER DEFAULT 0,
    opened_details BOOLEAN DEFAULT false,
    added_to_comparison BOOLEAN DEFAULT false,

    -- User feedback: "too_expensive", "not_relevant", "perfect", etc.
    feedback VARCHAR(100),

    -- Implicit score based on interaction patterns (0-1)
    implicit_score DOUBLE PRECISION DEFAULT 0,

    -- Search context that led to this product
    search_query TEXT,
    search_type VARCHAR(50),

    -- Position in results
    position_in_results INTEGER DEFAULT 0,

    -- Interaction sequence (order of this interaction in the session)
    interaction_sequence INTEGER DEFAULT 0,

    -- Timestamps
    interacted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_product_interactions_user_time
    ON product_interactions(user_id, interacted_at);

CREATE INDEX IF NOT EXISTS idx_product_interactions_session
    ON product_interactions(session_id, interaction_sequence);

CREATE INDEX IF NOT EXISTS idx_product_interactions_product
    ON product_interactions(product_id, interaction_type, interacted_at);

CREATE INDEX IF NOT EXISTS idx_product_interactions_category
    ON product_interactions(user_id, product_category, implicit_score);

CREATE INDEX IF NOT EXISTS idx_product_interactions_brand
    ON product_interactions(user_id, product_brand);

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_product_interactions_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to auto-update updated_at
CREATE TRIGGER trigger_product_interactions_updated_at
    BEFORE UPDATE ON product_interactions
    FOR EACH ROW
    EXECUTE FUNCTION update_product_interactions_updated_at();
