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
}
