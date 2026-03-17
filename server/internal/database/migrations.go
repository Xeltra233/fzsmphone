package database

type migration struct {
	Version int
	Name    string
	SQL     string
}

// migrations are applied in order. Add new migrations at the end.
var migrations = []migration{
	{
		Version: 1,
		Name:    "create_users",
		SQL: `
			CREATE TABLE IF NOT EXISTS users (
				id            BIGSERIAL PRIMARY KEY,
				discord_id    TEXT UNIQUE NOT NULL,
				username      TEXT NOT NULL,
				display_name  TEXT NOT NULL DEFAULT '',
				avatar_url    TEXT NOT NULL DEFAULT '',
				email         TEXT NOT NULL DEFAULT '',
				role          TEXT NOT NULL DEFAULT 'user',
				settings      JSONB NOT NULL DEFAULT '{}',
				created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_users_discord_id ON users(discord_id);
		`,
	},
	{
		Version: 2,
		Name:    "create_characters",
		SQL: `
			CREATE TABLE IF NOT EXISTS characters (
				id            BIGSERIAL PRIMARY KEY,
				user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				name          TEXT NOT NULL,
				avatar_url    TEXT NOT NULL DEFAULT '',
				description   TEXT NOT NULL DEFAULT '',
				personality   TEXT NOT NULL DEFAULT '',
				system_prompt TEXT NOT NULL DEFAULT '',
				greeting      TEXT NOT NULL DEFAULT '',
				is_public     BOOLEAN NOT NULL DEFAULT false,
				tags          TEXT[] NOT NULL DEFAULT '{}',
				extra         JSONB NOT NULL DEFAULT '{}',
				created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_characters_user_id ON characters(user_id);
		`,
	},
	{
		Version: 3,
		Name:    "create_conversations",
		SQL: `
			CREATE TABLE IF NOT EXISTS conversations (
				id            BIGSERIAL PRIMARY KEY,
				user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				character_id  BIGINT NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
				title         TEXT NOT NULL DEFAULT '',
				is_group      BOOLEAN NOT NULL DEFAULT false,
				last_message  TEXT NOT NULL DEFAULT '',
				last_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_conversations_user_id ON conversations(user_id);
		`,
	},
	{
		Version: 4,
		Name:    "create_messages",
		SQL: `
			CREATE TABLE IF NOT EXISTS messages (
				id              BIGSERIAL PRIMARY KEY,
				conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
				role            TEXT NOT NULL DEFAULT 'user',
				content         TEXT NOT NULL DEFAULT '',
				msg_type        TEXT NOT NULL DEFAULT 'text',
				extra           JSONB NOT NULL DEFAULT '{}',
				created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_messages_conversation ON messages(conversation_id, created_at);
		`,
	},
	{
		Version: 5,
		Name:    "create_posts",
		SQL: `
			CREATE TABLE IF NOT EXISTS posts (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				board       TEXT NOT NULL DEFAULT 'general',
				title       TEXT NOT NULL,
				content     TEXT NOT NULL DEFAULT '',
				images      TEXT[] NOT NULL DEFAULT '{}',
				likes       INTEGER NOT NULL DEFAULT 0,
				comments    INTEGER NOT NULL DEFAULT 0,
				pinned      BOOLEAN NOT NULL DEFAULT false,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_posts_board ON posts(board, created_at DESC);
			CREATE INDEX IF NOT EXISTS idx_posts_user ON posts(user_id);
		`,
	},
	{
		Version: 6,
		Name:    "create_comments",
		SQL: `
			CREATE TABLE IF NOT EXISTS comments (
				id          BIGSERIAL PRIMARY KEY,
				post_id     BIGINT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				parent_id   BIGINT REFERENCES comments(id) ON DELETE CASCADE,
				content     TEXT NOT NULL,
				likes       INTEGER NOT NULL DEFAULT 0,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_comments_post ON comments(post_id, created_at);
		`,
	},
	{
		Version: 7,
		Name:    "create_weibo",
		SQL: `
			CREATE TABLE IF NOT EXISTS weibo_posts (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				content     TEXT NOT NULL,
				images      TEXT[] NOT NULL DEFAULT '{}',
				likes       INTEGER NOT NULL DEFAULT 0,
				reposts     INTEGER NOT NULL DEFAULT 0,
				comments    INTEGER NOT NULL DEFAULT 0,
				is_hot      BOOLEAN NOT NULL DEFAULT false,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_weibo_created ON weibo_posts(created_at DESC);
		`,
	},
	{
		Version: 8,
		Name:    "create_moments",
		SQL: `
			CREATE TABLE IF NOT EXISTS moments (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				content     TEXT NOT NULL DEFAULT '',
				images      TEXT[] NOT NULL DEFAULT '{}',
				location    TEXT NOT NULL DEFAULT '',
				visibility  TEXT NOT NULL DEFAULT 'public',
				likes       INTEGER NOT NULL DEFAULT 0,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_moments_user ON moments(user_id, created_at DESC);
		`,
	},
	{
		Version: 9,
		Name:    "create_diaries",
		SQL: `
			CREATE TABLE IF NOT EXISTS diaries (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				title       TEXT NOT NULL,
				content     TEXT NOT NULL DEFAULT '',
				mood        TEXT NOT NULL DEFAULT '',
				weather     TEXT NOT NULL DEFAULT '',
				tags        TEXT[] NOT NULL DEFAULT '{}',
				is_private  BOOLEAN NOT NULL DEFAULT true,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_diaries_user ON diaries(user_id, created_at DESC);
		`,
	},
	{
		Version: 10,
		Name:    "create_takeaway",
		SQL: `
			CREATE TABLE IF NOT EXISTS restaurants (
				id          BIGSERIAL PRIMARY KEY,
				name        TEXT NOT NULL,
				category    TEXT NOT NULL DEFAULT '',
				image_url   TEXT NOT NULL DEFAULT '',
				rating      NUMERIC(2,1) NOT NULL DEFAULT 5.0,
				delivery_time TEXT NOT NULL DEFAULT '30-45分钟',
				min_order   NUMERIC(10,2) NOT NULL DEFAULT 0,
				items       JSONB NOT NULL DEFAULT '[]',
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS orders (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				restaurant_id BIGINT NOT NULL REFERENCES restaurants(id),
				items       JSONB NOT NULL DEFAULT '[]',
				total       NUMERIC(10,2) NOT NULL DEFAULT 0,
				status      TEXT NOT NULL DEFAULT 'pending',
				address     TEXT NOT NULL DEFAULT '',
				note        TEXT NOT NULL DEFAULT '',
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_orders_user ON orders(user_id, created_at DESC);
		`,
	},
	{
		Version: 11,
		Name:    "create_music",
		SQL: `
			CREATE TABLE IF NOT EXISTS playlists (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				name        TEXT NOT NULL,
				cover_url   TEXT NOT NULL DEFAULT '',
				is_public   BOOLEAN NOT NULL DEFAULT false,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS songs (
				id          BIGSERIAL PRIMARY KEY,
				title       TEXT NOT NULL,
				artist      TEXT NOT NULL DEFAULT '',
				album       TEXT NOT NULL DEFAULT '',
				cover_url   TEXT NOT NULL DEFAULT '',
				url         TEXT NOT NULL DEFAULT '',
				duration    INTEGER NOT NULL DEFAULT 0,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);

			CREATE TABLE IF NOT EXISTS playlist_songs (
				playlist_id BIGINT NOT NULL REFERENCES playlists(id) ON DELETE CASCADE,
				song_id     BIGINT NOT NULL REFERENCES songs(id) ON DELETE CASCADE,
				position    INTEGER NOT NULL DEFAULT 0,
				PRIMARY KEY (playlist_id, song_id)
			);
		`,
	},
	{
		Version: 12,
		Name:    "create_user_personas",
		SQL: `
			CREATE TABLE IF NOT EXISTS user_personas (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				name        TEXT NOT NULL,
				description TEXT NOT NULL DEFAULT '',
				avatar_url  TEXT NOT NULL DEFAULT '',
				is_default  BOOLEAN NOT NULL DEFAULT false,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_personas_user ON user_personas(user_id);
		`,
	},
	{
		Version: 13,
		Name:    "create_worldbook",
		SQL: `
			CREATE TABLE IF NOT EXISTS worldbook_entries (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				key         TEXT NOT NULL,
				content     TEXT NOT NULL DEFAULT '',
				keywords    TEXT[] NOT NULL DEFAULT '{}',
				is_enabled  BOOLEAN NOT NULL DEFAULT true,
				priority    INTEGER NOT NULL DEFAULT 0,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_worldbook_user ON worldbook_entries(user_id);
		`,
	},
	{
		Version: 14,
		Name:    "create_settings",
		SQL: `
			CREATE TABLE IF NOT EXISTS app_settings (
				id          BIGSERIAL PRIMARY KEY,
				key         TEXT UNIQUE NOT NULL,
				value       JSONB NOT NULL DEFAULT '{}',
				updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
		`,
	},
	{
		Version: 15,
		Name:    "create_presets",
		SQL: `
			CREATE TABLE IF NOT EXISTS presets (
				id              BIGSERIAL PRIMARY KEY,
				user_id         BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				name            TEXT NOT NULL,
				emoji           TEXT NOT NULL DEFAULT '🎭',
				category        TEXT NOT NULL DEFAULT '',
				description     TEXT NOT NULL DEFAULT '',
				content         TEXT NOT NULL DEFAULT '',
				prefill         TEXT NOT NULL DEFAULT '',
				enable_prefill  BOOLEAN NOT NULL DEFAULT false,
				is_builtin      BOOLEAN NOT NULL DEFAULT false,
				gradient        TEXT NOT NULL DEFAULT 'linear-gradient(135deg, #667eea, #764ba2)',
				created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_presets_user ON presets(user_id);
		`,
	},
	{
		Version: 16,
		Name:    "create_call_records",
		SQL: `
			CREATE TABLE IF NOT EXISTS call_records (
				id            BIGSERIAL PRIMARY KEY,
				user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				character_id  TEXT NOT NULL DEFAULT '',
				name          TEXT NOT NULL,
				number        TEXT NOT NULL DEFAULT '',
				avatar        TEXT NOT NULL DEFAULT '',
				type          TEXT NOT NULL DEFAULT 'outgoing',
				call_type     TEXT NOT NULL DEFAULT 'voice',
				duration      TEXT NOT NULL DEFAULT '',
				created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_call_records_user ON call_records(user_id, created_at DESC);
		`,
	},
	{
		Version: 17,
		Name:    "create_sms",
		SQL: `
			CREATE TABLE IF NOT EXISTS sms_threads (
				id            BIGSERIAL PRIMARY KEY,
				user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				recipient     TEXT NOT NULL,
				number        TEXT NOT NULL DEFAULT '',
				character_id  TEXT NOT NULL DEFAULT '',
				avatar        TEXT NOT NULL DEFAULT '',
				last_content  TEXT NOT NULL DEFAULT '',
				last_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
				created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_sms_threads_user ON sms_threads(user_id, last_at DESC);

			CREATE TABLE IF NOT EXISTS sms_messages (
				id          BIGSERIAL PRIMARY KEY,
				thread_id   BIGINT NOT NULL REFERENCES sms_threads(id) ON DELETE CASCADE,
				role        TEXT NOT NULL DEFAULT 'user',
				content     TEXT NOT NULL,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_sms_messages_thread ON sms_messages(thread_id, created_at);
		`,
	},
	{
		Version: 18,
		Name:    "create_wallet",
		SQL: `
			CREATE TABLE IF NOT EXISTS wallets (
				id       BIGSERIAL PRIMARY KEY,
				user_id  BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				balance  NUMERIC(12,2) NOT NULL DEFAULT 10000.00
			);

			CREATE TABLE IF NOT EXISTS wallet_transactions (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				type        TEXT NOT NULL,
				amount      NUMERIC(12,2) NOT NULL,
				description TEXT NOT NULL DEFAULT '',
				target      TEXT NOT NULL DEFAULT '',
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_wallet_tx_user ON wallet_transactions(user_id, created_at DESC);
		`,
	},
	{
		Version: 19,
		Name:    "create_game_records",
		SQL: `
			CREATE TABLE IF NOT EXISTS game_records (
				id          BIGSERIAL PRIMARY KEY,
				user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
				game        TEXT NOT NULL,
				detail      TEXT NOT NULL DEFAULT '',
				amount      NUMERIC(12,2) NOT NULL DEFAULT 0,
				win         BOOLEAN NOT NULL DEFAULT false,
				created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			CREATE INDEX IF NOT EXISTS idx_game_records_user ON game_records(user_id, created_at DESC);
		`,
	},
	{
		Version: 20,
		Name:    "add_user_ban_fields",
		SQL: `
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_banned BOOLEAN NOT NULL DEFAULT false;
ALTER TABLE users ADD COLUMN IF NOT EXISTS ban_reason TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS banned_at TIMESTAMPTZ;
`,
	},
	{
		Version: 21,
		Name:    "add_password_and_super_admin",
		SQL: `
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS is_super_admin BOOLEAN NOT NULL DEFAULT false;
`,
	},
	{
		Version: 22,
		Name:    "create_user_api_settings",
		SQL: `
CREATE TABLE IF NOT EXISTS user_api_settings (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	api_key TEXT NOT NULL DEFAULT '',
	api_url TEXT NOT NULL DEFAULT '',
	model TEXT NOT NULL DEFAULT '',
	is_global BOOLEAN NOT NULL DEFAULT false,
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_user_api_settings_user ON user_api_settings(user_id);
`,
	},
	{
		Version: 23,
		Name:    "allow_null_discord_id",
		SQL: `
ALTER TABLE users ALTER COLUMN discord_id DROP NOT NULL;
`,
	},
	{
		Version: 24,
		Name:    "add_credits_system",
		SQL: `
ALTER TABLE users ADD COLUMN IF NOT EXISTS credits INTEGER NOT NULL DEFAULT 1000;
ALTER TABLE users ADD COLUMN IF NOT EXISTS total_tokens INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_signin_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS signin_streak INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS last_active_at TIMESTAMPTZ;
ALTER TABLE users ADD COLUMN IF NOT EXISTS invite_code VARCHAR(20) NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS invited_by BIGINT REFERENCES users(id);
ALTER TABLE users ADD COLUMN IF NOT EXISTS invite_rewards_claimed INTEGER NOT NULL DEFAULT 0;
ALTER TABLE users ADD COLUMN IF NOT EXISTS model_display_names JSONB NOT NULL DEFAULT '{}';

CREATE TABLE IF NOT EXISTS invite_rewards (
	id BIGSERIAL PRIMARY KEY,
	inviter_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	invited_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	reward_credits INTEGER NOT NULL DEFAULT 0,
	reward_type VARCHAR(20) NOT NULL DEFAULT 'signup',
	created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS signin_records (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  credits_earned INTEGER NOT NULL DEFAULT 0,
  streak_bonus INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS usage_records (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  tokens_used INTEGER NOT NULL DEFAULT 0,
  model VARCHAR(100) NOT NULL DEFAULT '',
  credits_cost INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coupon_codes (
  id BIGSERIAL PRIMARY KEY,
  code VARCHAR(20) NOT NULL UNIQUE,
  credits INTEGER NOT NULL DEFAULT 0,
  max_uses INTEGER NOT NULL DEFAULT 1,
  current_uses INTEGER NOT NULL DEFAULT 0,
  expires_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS coupon_redeem_records (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  coupon_id BIGINT NOT NULL REFERENCES coupon_codes(id),
  credits_earned INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS system_logs (
  id BIGSERIAL PRIMARY KEY,
  level VARCHAR(20) NOT NULL DEFAULT 'INFO',
  action VARCHAR(100) NOT NULL,
  user_id BIGINT,
  username VARCHAR(100),
  details JSONB,
  ip_address VARCHAR(50),
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_system_logs_action ON system_logs(action);
CREATE INDEX IF NOT EXISTS idx_system_logs_user ON system_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_system_logs_created ON system_logs(created_at DESC);
`,
	},
	{
		Version: 26,
		Name:    "add_default_settings",
		SQL: `
INSERT INTO app_settings (key, value, updated_at) VALUES 
  ('default_credits', '1000', NOW()),
  ('signin_daily_credits', '10', NOW()),
  ('signin_streak_bonus', '5', NOW()),
  ('invite_reward_credits', '100', NOW()),
  ('signin_enabled', 'true', NOW()),
  ('invite_enabled', 'true', NOW())
ON CONFLICT (key) DO NOTHING;
`,
	},
}
