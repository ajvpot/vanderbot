--! Previous: -
--! Hash: sha1:61c2131138e47edbed3fc17a7d89e030420f73cf

-- Enter migration here
DROP TABLE IF EXISTS presence CASCADE;

CREATE TABLE presence
(
    user_id  BIGINT NOT NULL,
    guild_id BIGINT NOT NULL,
    blob     jsonb  NOT NULL
);
