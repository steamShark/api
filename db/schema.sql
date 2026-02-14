PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS websites (
  id TEXT PRIMARY KEY,                    -- size:36 UUID string
  url TEXT NOT NULL,                      -- size:255
  domain TEXT NOT NULL UNIQUE,            -- uniqueIndex
  ssl_certificate INTEGER NOT NULL DEFAULT 0, -- bool: 0/1
  display_name TEXT,                      -- nullable
  tld TEXT NOT NULL,                      -- size:63
  description TEXT,                       -- nullable
  type TEXT NOT NULL DEFAULT 'website',   -- website, tool, extension
  is_not_trusted INTEGER NOT NULL DEFAULT 1,  -- pointer bool in Go, but stored as 0/1
  is_official INTEGER NOT NULL DEFAULT 0,
  steam_login_present INTEGER NOT NULL DEFAULT 0,
  verified INTEGER NOT NULL DEFAULT 0,
  risk_score REAL NOT NULL DEFAULT 0.0,   -- 0..100
  risk_level TEXT NOT NULL DEFAULT 'unknown', -- unknown, none, low, medium, high, critical
  status TEXT NOT NULL DEFAULT 'active',  -- active, inactive, blocked, archived
  notes TEXT,                             -- nullable text
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CHECK (ssl_certificate IN (0,1)),
  CHECK (is_not_trusted IN (0,1)),
  CHECK (is_official IN (0,1)),
  CHECK (steam_login_present IN (0,1)),
  CHECK (verified IN (0,1)),
  CHECK (risk_score >= 0.0 AND risk_score <= 100.0),
  CHECK (type IN ('website','tool','extension')),
  CHECK (risk_level IN ('unknown','none','low','medium','high','critical')),
  CHECK (status IN ('active','inactive','blocked','archived'))
);

-- Helpful index for lookups by URL (optional)
CREATE INDEX IF NOT EXISTS idx_websites_url ON websites(url);

-- Trigger to keep updated_at current (SQLite doesn't do this automatically)
CREATE TRIGGER IF NOT EXISTS trg_websites_updated_at
AFTER UPDATE ON websites
FOR EACH ROW
BEGIN
  UPDATE websites SET updated_at = datetime('now') WHERE id = OLD.id;
END;