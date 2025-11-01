CREATE TABLE IF NOT EXISTS item (
    id UUID PRIMARY KEY,
    type TEXT NOT NULL,
    amount DOUBLE PRECISION NOT NULL,
    category TEXT NOT NULL,
    description TEXT NOT NULL,
    create_date TIMESTAMP
)