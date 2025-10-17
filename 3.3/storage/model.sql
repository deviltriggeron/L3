CREATE EXTENSION IF NOT EXISTS ltree;

CREATE TABLE IF NOT EXISTS comments (
    comment_id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    parent_id UUID REFERENCES comments (comment_id) ON DELETE CASCADE,
    user_name TEXT,
    comment TEXT,
    path LTREE,
    date TIMESTAMP DEFAULT NOW()
);