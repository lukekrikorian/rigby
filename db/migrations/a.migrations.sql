CREATE TABLE IF NOT EXISTS migrations (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  performed TIMESTAMP DEFAULT current_timestamp
);
