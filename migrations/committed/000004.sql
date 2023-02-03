--! Previous: sha1:f06f5bac7a7f91f58fff61547fe5b26bf7fb323c
--! Hash: sha1:875a2400a3f198033fc726257819709a9c79617a

-- Enter migration here
DROP TABLE IF EXISTS presence CASCADE;

CREATE TABLE presence
(
    guild_id   TEXT                      NOT NULL,
    blob       jsonb                     NOT NULL,
    created_at timestamptz DEFAULT NOW() NOT NULL,
    user_id    TEXT                      NOT NULL
        GENERATED ALWAYS AS ((blob -> 'user' -> 'id')) STORED
);

ALTER TABLE presence
    ADD CONSTRAINT presence_guild_id_blob_unique
        UNIQUE (guild_id, blob);
