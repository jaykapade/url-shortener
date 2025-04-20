-- +goose Up
CREATE TABLE links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_code VARCHAR(255) UNIQUE,
    full_url TEXT NOT NULL,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    click_count INT DEFAULT 0
);

-- +goose Down
DROP TABLE links;
