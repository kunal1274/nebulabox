-- Rollback for initial schema migration
-- Drops all tables in reverse dependency order

DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS networks;
DROP TABLE IF EXISTS tenants;
DROP TABLE IF EXISTS teams;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS templates;
DROP TABLE IF EXISTS container_groups;
DROP TABLE IF EXISTS nodes;
DROP TABLE IF EXISTS deployments;
DROP TABLE IF EXISTS snapshots;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS invites;
DROP TABLE IF EXISTS workspace_members;
DROP TABLE IF EXISTS workspaces;
DROP TABLE IF EXISTS images;
DROP TABLE IF EXISTS containers;

