-- Enter migration here
create table presence
(
    user_id  BIGINT,
    guild_id BIGINT,
    blob     jsonb
);

