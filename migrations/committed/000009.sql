--! Previous: sha1:5f3b075fe57913e84f1a016c9848c866db5667ed
--! Hash: sha1:296b27f87b088c6746db4f1cebf12fd428e80a8d

-- Enter migration here
ALTER TABLE message
    DROP CONSTRAINT IF EXISTS message_unique;
DROP INDEX IF EXISTS message_unique;

ALTER TABLE message
    DROP COLUMN IF EXISTS is_delete;

ALTER TABLE message
    ADD CONSTRAINT message_unique
        UNIQUE (message_id, created_at, edited_at);
