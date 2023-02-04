--! Previous: sha1:987277198fce74d95e52964f8f48fd1eaff9906b
--! Hash: sha1:5f3b075fe57913e84f1a016c9848c866db5667ed

-- Enter migration here
ALTER TABLE message
    DROP CONSTRAINT IF EXISTS message_unique;
DROP INDEX IF EXISTS message_unique;

ALTER TABLE message
    DROP COLUMN IF EXISTS is_delete;

ALTER TABLE message
    ADD CONSTRAINT message_unique
        UNIQUE (message_id, edited_at);
