--! Previous: sha1:04822d89c4546d224e57504cccd9ad0f000cdfa9
--! Hash: sha1:987277198fce74d95e52964f8f48fd1eaff9906b

-- Enter migration here
DROP TABLE IF EXISTS message CASCADE;

CREATE TABLE message
(
    blob       jsonb NOT NULL,
    created_at TEXT  NOT NULL
        GENERATED ALWAYS AS ((blob ->> 'timestamp')) STORED,
    edited_at  TEXT
               GENERATED ALWAYS AS ((blob ->> 'edited_timestamp')) STORED,
    message_id TEXT  NOT NULL
        GENERATED ALWAYS AS (blob ->> 'id') STORED,
    channel_id TEXT  NOT NULL
        GENERATED ALWAYS AS (blob ->> 'channel_id') STORED,
    guild_id   TEXT
               GENERATED ALWAYS AS (blob ->> 'guild_id') STORED,
    is_delete  BOOL  NOT NULL
);

ALTER TABLE message
    ADD CONSTRAINT message_unique
        UNIQUE (message_id, edited_at, is_delete);


ALTER TABLE presence
    DROP CONSTRAINT IF EXISTS presence_guild_id_blob_unique;
DROP INDEX IF EXISTS presence_guild_id_blob_unique;
ALTER TABLE presence
    ADD CONSTRAINT presence_guild_id_blob_unique
        UNIQUE (guild_id, blob);



ALTER TABLE PRESENCE
    DROP COLUMN IF EXISTS user_id;
ALTER TABLE presence
    ADD user_id TEXT GENERATED ALWAYS AS (blob -> 'user' ->> 'id') STORED;
