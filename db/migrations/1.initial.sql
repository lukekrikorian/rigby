CREATE TABLE IF NOT EXISTS posts (
  id TEXT PRIMARY KEY NOT NULL,
  userid TEXT NOT NULL,
  author TEXT NOT NULL,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  gamerrage INTEGER DEFAULT 0,
  votes INTEGER DEFAULT 0,
  created TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS users (
  id TEXT PRIMARY KEY NOT NULL,
  username TEXT NOT NULL,
  password TEXT NOT NULL,
  created TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS comments (
  id TEXT PRIMARY KEY NOT NULL,
  postid TEXT NOT NULL,
  userid TEXT NOT NULL,
  author TEXT NOT NULL,
  body TEXT NOT NULL,
  created TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS replies (
  id TEXT NOT NULL,
  parentid TEXT NOT NULL,
  userid TEXT NOT NULL,
  author TEXT NOT NULL,
  body TEXT NOT NULL,
  created TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS votes (
  userid TEXT NOT NULL,
  postid TEXT NOT NULL,
  PRIMARY KEY (userid, postid)
);

CREATE TABLE IF NOT EXISTS sessions (
  token TEXT NOT NULL,
  userid TEXT NOT NULL
);
