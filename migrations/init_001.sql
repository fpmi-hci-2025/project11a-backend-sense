CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ENUMS
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='user_role') THEN
    CREATE TYPE user_role AS ENUM ('reader','user','creator','expert','super');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='publication_type') THEN
    CREATE TYPE publication_type AS ENUM ('quote','post','article');
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname='visibility_type') THEN
    CREATE TYPE visibility_type AS ENUM ('public','community','private');
  END IF;
END$$;

BEGIN;

-- USERS
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  username text NOT NULL UNIQUE,
  icon_url text,
  email text UNIQUE,
  phone text UNIQUE,
  registered_at timestamptz NOT NULL DEFAULT now(),
  description text,
  role user_role NOT NULL DEFAULT 'user',
  statistic jsonb NOT NULL DEFAULT '{}'::jsonb,
  followers_count int NOT NULL DEFAULT 0 CHECK (followers_count >= 0),
  following_count int NOT NULL DEFAULT 0 CHECK (following_count >= 0)
);

-- PUBLICATIONS
CREATE TABLE IF NOT EXISTS publications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  author_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type publication_type NOT NULL,
  content text,
  source text,
  publication_date timestamptz NOT NULL DEFAULT now(),
  visibility visibility_type NOT NULL DEFAULT 'public',
  likes_count int NOT NULL DEFAULT 0 CHECK (likes_count >= 0),
  comments_count int NOT NULL DEFAULT 0 CHECK (comments_count >= 0),
  saved_count int NOT NULL DEFAULT 0 CHECK (saved_count >= 0)
);
CREATE INDEX IF NOT EXISTS idx_publications_author ON publications(author_id);
CREATE INDEX IF NOT EXISTS idx_publications_pubdate ON publications(publication_date);

-- MEDIA
CREATE TABLE IF NOT EXISTS media_assets (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  owner_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  data bytea NOT NULL,
  filename text,
  mime text NOT NULL,
  width integer CHECK (width IS NULL OR width >= 0),
  height integer CHECK (height IS NULL OR height >= 0),
  exif jsonb,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_media_owner ON media_assets(owner_id);

CREATE TABLE IF NOT EXISTS publication_media (
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  media_id uuid NOT NULL REFERENCES media_assets(id) ON DELETE CASCADE,
  ord integer NOT NULL DEFAULT 0 CHECK (ord >= 0),
  PRIMARY KEY (publication_id, media_id),
  CONSTRAINT uq_pub_media_ord UNIQUE (publication_id, ord)
);

-- COMMENTS
CREATE TABLE IF NOT EXISTS comments (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  parent_id uuid REFERENCES comments(id) ON DELETE SET NULL,
  author_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  text text NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  likes_count integer NOT NULL DEFAULT 0 CHECK (likes_count >= 0)
);
CREATE INDEX IF NOT EXISTS idx_comments_pub_created ON comments(publication_id, created_at);
CREATE INDEX IF NOT EXISTS idx_comments_parent ON comments(parent_id);

-- LIKES
CREATE TABLE IF NOT EXISTS publication_likes (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_pub_like UNIQUE (user_id, publication_id)
);
CREATE INDEX IF NOT EXISTS idx_pub_likes_pub ON publication_likes(publication_id);

CREATE TABLE IF NOT EXISTS comment_likes (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  comment_id uuid NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_com_like UNIQUE (user_id, comment_id)
);
CREATE INDEX IF NOT EXISTS idx_com_likes_comment ON comment_likes(comment_id);

-- SAVED
CREATE TABLE IF NOT EXISTS saved_items (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  added_at timestamptz NOT NULL DEFAULT now(),
  note text,
  CONSTRAINT uq_saved UNIQUE (user_id, publication_id)
);
CREATE INDEX IF NOT EXISTS idx_saved_items_user ON saved_items(user_id);
CREATE INDEX IF NOT EXISTS idx_saved_items_pub ON saved_items(publication_id);

-- RECOMMENDATIONS
CREATE TABLE IF NOT EXISTS recommendations (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  algorithm text NOT NULL,
  reason text NOT NULL,
  rank int NOT NULL CHECK (rank >= 0),
  created_at timestamptz NOT NULL DEFAULT now(),
  hidden boolean NOT NULL DEFAULT false
);
CREATE INDEX IF NOT EXISTS idx_recs_user_rank ON recommendations(user_id, rank);
CREATE INDEX IF NOT EXISTS idx_recs_pub ON recommendations(publication_id);

-- VIEWS
CREATE TABLE IF NOT EXISTS publications_views (
  view_uuid uuid PRIMARY KEY DEFAULT gen_random_uuid(),     
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  viewed_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_pub_view_triplet UNIQUE (user_id, publication_id, viewed_at)
);

CREATE INDEX IF NOT EXISTS idx_pub_views_pub_ts ON publications_views(publication_id, viewed_at);
CREATE INDEX IF NOT EXISTS idx_pub_views_user_ts ON publications_views(user_id, viewed_at);

-- FOLLOWERS/FOLLOWING (подписки пользователей)
CREATE TABLE IF NOT EXISTS user_follows (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  follower_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  following_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT uq_user_follows UNIQUE (follower_id, following_id),
  CONSTRAINT chk_no_self_follow CHECK (follower_id != following_id)
);
CREATE INDEX IF NOT EXISTS idx_user_follows_follower ON user_follows(follower_id);
CREATE INDEX IF NOT EXISTS idx_user_follows_following ON user_follows(following_id);

-- TAGS/HASHTAGS (теги для публикаций)
CREATE TABLE IF NOT EXISTS tags (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name text NOT NULL UNIQUE,
  description text,
  usage_count int NOT NULL DEFAULT 0 CHECK (usage_count >= 0),
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
CREATE INDEX IF NOT EXISTS idx_tags_usage ON tags(usage_count DESC);

CREATE TABLE IF NOT EXISTS publication_tags (
  publication_id uuid NOT NULL REFERENCES publications(id) ON DELETE CASCADE,
  tag_id uuid NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (publication_id, tag_id)
);
CREATE INDEX IF NOT EXISTS idx_pub_tags_pub ON publication_tags(publication_id);
CREATE INDEX IF NOT EXISTS idx_pub_tags_tag ON publication_tags(tag_id);

-- NOTIFICATIONS (уведомления)
CREATE TABLE IF NOT EXISTS notifications (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  type text NOT NULL, -- 'like', 'comment', 'follow', 'mention'
  title text NOT NULL,
  message text NOT NULL,
  data jsonb, -- дополнительные данные (ID публикации, комментария и т.д.)
  is_read boolean NOT NULL DEFAULT false,
  created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX IF NOT EXISTS idx_notifications_user ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_read ON notifications(user_id, is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_created ON notifications(created_at DESC);

-- SESSIONS/TOKENS (сессии и токены)
CREATE TABLE IF NOT EXISTS user_sessions (
  id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  token_hash text NOT NULL UNIQUE, -- хеш JWT токена
  refresh_token_hash text, -- хеш refresh токена
  expires_at timestamptz NOT NULL,
  created_at timestamptz NOT NULL DEFAULT now(),
  last_used_at timestamptz NOT NULL DEFAULT now(),
  user_agent text,
  ip_address inet
);
CREATE INDEX IF NOT EXISTS idx_sessions_user ON user_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON user_sessions(token_hash);
CREATE INDEX IF NOT EXISTS idx_sessions_expires ON user_sessions(expires_at);


-- Триггеры для автоматического обновления счетчиков
CREATE OR REPLACE FUNCTION update_user_followers_count()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE users SET followers_count = followers_count + 1 WHERE id = NEW.following_id;
    UPDATE users SET following_count = following_count + 1 WHERE id = NEW.follower_id;
    RETURN NEW;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE users SET followers_count = followers_count - 1 WHERE id = OLD.following_id;
    UPDATE users SET following_count = following_count - 1 WHERE id = OLD.follower_id;
    RETURN OLD;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_followers_count
  AFTER INSERT OR DELETE ON user_follows
  FOR EACH ROW EXECUTE FUNCTION update_user_followers_count();

-- Триггер для обновления счетчика использования тегов
CREATE OR REPLACE FUNCTION update_tag_usage_count()
RETURNS TRIGGER AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    UPDATE tags SET usage_count = usage_count + 1 WHERE id = NEW.tag_id;
    RETURN NEW;
  ELSIF TG_OP = 'DELETE' THEN
    UPDATE tags SET usage_count = usage_count - 1 WHERE id = OLD.tag_id;
    RETURN OLD;
  END IF;
  RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_tag_usage
  AFTER INSERT OR DELETE ON publication_tags
  FOR EACH ROW EXECUTE FUNCTION update_tag_usage_count();

COMMIT;