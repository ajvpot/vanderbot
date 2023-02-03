--! Previous: sha1:875a2400a3f198033fc726257819709a9c79617a
--! Hash: sha1:d7c59ac531a7bd905a8b640aba7876f003e795ad

-- Enter migration here
ALTER TABLE message
    DROP COLUMN IF EXISTS is_delete;

ALTER TABLE message
    ADD is_delete bool NOT NULL DEFAULT FALSE;
