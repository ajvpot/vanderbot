--! Previous: sha1:61c2131138e47edbed3fc17a7d89e030420f73cf
--! Hash: sha1:8b80e29a8fb44f1fb0b609d5dafbb24d216e86b9

-- Enter migration here
DROP TABLE IF EXISTS message CASCADE;

CREATE TABLE message
(
    blob       jsonb NOT NULL,
    created_at TEXT  NOT NULL
        GENERATED ALWAYS AS ((blob -> 'timestamp')) STORED,
    edited_at  TEXT
               GENERATED ALWAYS AS ((blob -> 'edited_timestamp')) STORED,
    message_id TEXT  NOT NULL
        GENERATED ALWAYS AS (blob -> 'id') STORED,
    channel_id TEXT  NOT NULL
        GENERATED ALWAYS AS (blob -> 'channel_id') STORED,
    guild_id   TEXT
               GENERATED ALWAYS AS (blob -> 'guild_id') STORED
);

alter table message
    add constraint message_unique
        unique (message_id, edited_at);
