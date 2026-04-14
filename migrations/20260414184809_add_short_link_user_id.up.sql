ALTER TABLE short_link ADD COLUMN user_id UUID NULL;

ALTER TABLE short_link DROP CONSTRAINT unique_url;

ALTER TABLE short_link ADD CONSTRAINT unique_user_url UNIQUE (url, user_id);