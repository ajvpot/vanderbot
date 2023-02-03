--! Previous: sha1:8b80e29a8fb44f1fb0b609d5dafbb24d216e86b9
--! Hash: sha1:f06f5bac7a7f91f58fff61547fe5b26bf7fb323c

-- Enter migration here

alter table presence
    drop constraint if exists presence_blob_unique;

alter table presence
    add constraint presence_blob_unique
        unique (blob);
