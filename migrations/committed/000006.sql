--! Previous: sha1:d7c59ac531a7bd905a8b640aba7876f003e795ad
--! Hash: sha1:04822d89c4546d224e57504cccd9ad0f000cdfa9

ALTER TABLE message
    DROP CONSTRAINT IF EXISTS message_unique;
DROP INDEX IF EXISTS message_unique;
ALTER TABLE message
    ADD CONSTRAINT message_unique
        UNIQUE (message_id, edited_at, is_delete);
