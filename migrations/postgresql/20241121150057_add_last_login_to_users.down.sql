-- Migration: add_last_login_to_users
-- Created at: 2024-11-21 15:00:57

-- Write your DOWN migration SQL here
ALTER TABLE users DROP COLUMN last_login;