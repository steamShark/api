-- TABLE
CREATE TABLE IF NOT EXISTS websites (
  id text PRIMARY KEY,                    -- UUID string stored as text (or change to uuid)
  url text NOT NULL,
  domain text NOT NULL UNIQUE,
  ssl_certificate boolean NOT NULL DEFAULT false,
  display_name text,
  tld text NOT NULL,
  description text,
  type text NOT NULL DEFAULT 'website',
  is_not_trusted boolean NOT NULL DEFAULT true,
  is_official boolean NOT NULL DEFAULT false,
  steam_login_present boolean NOT NULL DEFAULT false,
  verified boolean NOT NULL DEFAULT false,
  risk_score double precision NOT NULL DEFAULT 0.0,
  risk_level text NOT NULL DEFAULT 'unknown',
  status text NOT NULL DEFAULT 'active',
  notes text,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),

  CHECK (risk_score >= 0.0 AND risk_score <= 100.0),
  CHECK (type IN ('website','tool','extension')),
  CHECK (risk_level IN ('unknown','none','low','medium','high','critical')),
  CHECK (status IN ('active','inactive','blocked','archived'))
);

-- INDEX
CREATE INDEX IF NOT EXISTS idx_websites_url ON websites(url);

-- UPDATED_AT trigger function
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger AS $$
BEGIN
  NEW.updated_at := now();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- TRIGGER
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1
    FROM pg_trigger
    WHERE tgname = 'trg_websites_updated_at'
  ) THEN
    CREATE TRIGGER trg_websites_updated_at
    BEFORE UPDATE ON websites
    FOR EACH ROW
    EXECUTE FUNCTION set_updated_at();
  END IF;
END;
$$;