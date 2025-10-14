CREATE TABLE IF NOT EXISTS url(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    alias TEXT NOT NULL UNIQUE,
    original_url TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);

CREATE TABLE IF NOT EXISTS analytics(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    url_id UUID NOT NULL REFERENCES url(id) ON DELETE CASCADE,
    user_Agent TEXT NOT NULL,
    ip_address TEXT,
    time_transitions TIMESTAMP DEFAULT NOW());

CREATE INDEX IF NOT EXISTS idx_analytics_url_id ON analytics(url_id);