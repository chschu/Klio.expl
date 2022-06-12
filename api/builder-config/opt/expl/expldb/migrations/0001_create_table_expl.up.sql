CREATE TABLE IF NOT EXISTS expl (
    id SERIAL PRIMARY KEY,               -- unique id; monotonicity is required for some queries
    key TEXT NOT NULL,                   -- entry key in its original notation
    key_normalized TEXT NOT NULL,        -- normalized entry key; main search criterion
    value TEXT NOT NULL,                 -- entry value
    created_by TEXT,                     -- user name who added the entry; NULL for some old entries
    created_at TIMESTAMP WITH TIME ZONE, -- instant when the entry was added; NULL for some old entries
    visible BOOL NOT NULL                -- flag indicating if the entry is visible
)
