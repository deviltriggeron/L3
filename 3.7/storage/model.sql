CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    pass TEXT NOT NULL,
    role TEXT NOT NULL CHECK (
        role IN ('admin', 'manager', 'viewer')
    )
);

INSERT INTO
    users (username, pass, role)
VALUES ('admin', 'admin', 'admin'),
    (
        'manager',
        'manager',
        'manager'
    ),
    ('viewer', 'viewer', 'viewer');

CREATE TABLE IF NOT EXISTS item (
    id UUID PRIMARY KEY,
    product TEXT NOT NULL,
    price FLOAT NOT NULL,
    description TEXT,
    count BIGINT,
    create_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS items_history (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    item_id UUID NOT NULL,
    action TEXT NOT NULL,
    old_data JSONB,
    new_data JSONB,
    changed_by TEXT NOT NULL,
    changed_at TIMESTAMP DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION log_item_changes() RETURNS trigger AS $$
DECLARE
    actor TEXT;
BEGIN
    actor := current_setting('jwt.user', true);

    IF actor IS NULL THEN
        actor := 'unknown';
    END IF;

    IF TG_OP = 'INSERT' THEN
        INSERT INTO items_history(item_id, action, new_data, changed_by)
        VALUES (NEW.id, 'INSERT', to_jsonb(NEW), actor);

        RETURN NEW;

    ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO items_history(item_id, action, old_data, new_data, changed_by)
        VALUES (NEW.id, 'UPDATE', to_jsonb(OLD), to_jsonb(NEW), actor);

        RETURN NEW;

    ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO items_history(item_id, action, old_data, changed_by)
        VALUES (OLD.id, 'DELETE', to_jsonb(OLD), actor);

        RETURN OLD;
    END IF;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER items_history_trigger
AFTER INSERT OR UPDATE OR DELETE ON item
FOR EACH ROW EXECUTE FUNCTION log_item_changes();