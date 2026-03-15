-- Active: 1771768570334@@127.0.0.1@5432@postgres
-- 001_create_tables.sql
-- Run this first

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id         UUID        PRIMARY KEY DEFAULT uuid_generate_v4(),
    name       TEXT        NOT NULL,
    email      TEXT        NOT NULL UNIQUE,
    gender     TEXT        NOT NULL CHECK (gender IN ('male', 'female', 'other')),
    birth_date DATE        NOT NULL
);

CREATE TABLE IF NOT EXISTS user_friends (
    user_id   UUID REFERENCES users(id) ON DELETE CASCADE,
    friend_id UUID REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, friend_id),
    -- Prevent self-friendship
    CONSTRAINT no_self_friendship CHECK (user_id <> friend_id)
);
