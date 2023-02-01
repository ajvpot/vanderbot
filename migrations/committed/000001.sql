--! Previous: -
--! Hash: sha1:d7dc968634be4da40d717fddcef4b62031e6ceb3

-- Enter migration here
DROP TABLE IF EXISTS presence CASCADE;

CREATE TABLE presence
(
    user_id  BIGINT,
    guild_id BIGINT,
    blob     jsonb
);
