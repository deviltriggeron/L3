CREATE TABLE IF NOT EXISTS users (
    user_id UUID PRIMARY KEY,
    user_name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    event_id UUID PRIMARY KEY,
    date TIMESTAMP NOT NULL,
    event_info TEXT NOT NULL,
    organizer_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    seats_count INTEGER NOT NULL,
    for_free BOOLEAN NOT NULL DEFAULT TRUE,
    price FLOAT,
    create_date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS event_participants (
    event_id UUID NOT NULL REFERENCES events (event_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (event_id, user_id)
);

CREATE TABLE IF NOT EXISTS bookings (
    booking_id UUID PRIMARY KEY,
    event_id UUID NOT NULL REFERENCES events (event_id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    UNIQUE (event_id, user_id)
);