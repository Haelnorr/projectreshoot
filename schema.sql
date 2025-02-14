PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS jwtblacklist (
jti TEXT PRIMARY KEY CHECK(jti GLOB '[0-9a-fA-F-]*'),
exp INTEGER NOT NULL
) STRICT;
CREATE TABLE IF NOT EXISTS "users" (
    id INTEGER PRIMARY KEY AUTOINCREMENT, 
    username TEXT NOT NULL UNIQUE, 
    password_hash TEXT DEFAULT "", 
    created_at INTEGER DEFAULT (unixepoch())
) STRICT;
CREATE TRIGGER cleanup_expired_tokens
AFTER INSERT ON jwtblacklist
BEGIN
DELETE FROM jwtblacklist WHERE exp < strftime('%s', 'now');
END;
COMMIT;
