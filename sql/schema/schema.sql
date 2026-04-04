-- Users table
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Messages table (active messages)
CREATE TABLE IF NOT EXISTS messages (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  text VARCHAR(180) NOT NULL,
  latitude DECIMAL(10, 8) NOT NULL,
  longitude DECIMAL(11, 8) NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  expires_at TIMESTAMP DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours') NOT NULL
);

-- Archived messages table (expired messages for audit)
CREATE TABLE IF NOT EXISTS archived_messages (
  id SERIAL PRIMARY KEY,
  original_message_id INTEGER NOT NULL,
  user_id INTEGER NOT NULL,
  text VARCHAR(180) NOT NULL,
  latitude DECIMAL(10, 8) NOT NULL,
  longitude DECIMAL(11, 8) NOT NULL,
  created_at TIMESTAMP NOT NULL,
  expired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  archived_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for users table
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Indexes for messages table
CREATE INDEX IF NOT EXISTS idx_messages_user_id ON messages(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_expires_at ON messages(expires_at);
CREATE INDEX IF NOT EXISTS idx_messages_location ON messages (latitude, longitude);

-- Indexes for archived_messages table
CREATE INDEX IF NOT EXISTS idx_archived_messages_user_id ON archived_messages(user_id);
CREATE INDEX IF NOT EXISTS idx_archived_messages_created_at ON archived_messages(created_at);

-- Function to archive expired messages
CREATE OR REPLACE FUNCTION archive_expired_messages()
RETURNS void AS $$
BEGIN
  INSERT INTO archived_messages (original_message_id, user_id, text, latitude, longitude, created_at, expired_at)
  SELECT id, user_id, text, latitude, longitude, created_at, CURRENT_TIMESTAMP
  FROM messages
  WHERE expires_at <= CURRENT_TIMESTAMP;

  DELETE FROM messages WHERE expires_at <= CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;
