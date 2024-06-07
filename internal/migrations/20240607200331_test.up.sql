SET statement_timeout = 0;

CREATE EXTENSION IF NOT EXISTS pg_trgm WITH SCHEMA public;

COMMENT ON EXTENSION pg_trgm IS 'text similarity measurement and index searching based on trigrams';


CREATE TABLE IF NOT EXISTS favourite_links (
    link_id uuid NOT NULL,
    user_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS link_meta (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    user_id uuid,
    link_id uuid,
    clicked bigint,
    last_used timestamp with time zone
);

CREATE TABLE IF NOT EXISTS links (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    title character varying,
    description character varying,
    tags character varying,
    icon character varying,
    url character varying,
    updated timestamp with time zone,
    sharing character varying
);

CREATE TABLE IF NOT EXISTS team_links (
    team_id uuid NOT NULL,
    link_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    name character varying,
    description character varying
);

CREATE TABLE IF NOT EXISTS user_to_links (
    link_id uuid NOT NULL,
    user_id uuid NOT NULL
);

CREATE TABLE IF NOT EXISTS user_to_teams (
    user_id uuid NOT NULL,
    team_id uuid NOT NULL,
    permission character varying
);

CREATE TABLE IF NOT EXISTS users (
    id uuid DEFAULT gen_random_uuid() NOT NULL,
    username character varying,
    email character varying
);

ALTER TABLE ONLY favourite_links
    ADD CONSTRAINT favourite_links_pkey PRIMARY KEY (link_id, user_id);

ALTER TABLE ONLY link_meta
    ADD CONSTRAINT link_meta_pkey PRIMARY KEY (id);

ALTER TABLE ONLY links
    ADD CONSTRAINT links_pkey PRIMARY KEY (id);

ALTER TABLE ONLY team_links
    ADD CONSTRAINT team_links_pkey PRIMARY KEY (team_id, link_id);

ALTER TABLE ONLY teams
    ADD CONSTRAINT teams_pkey PRIMARY KEY (id);

ALTER TABLE ONLY user_to_links
    ADD CONSTRAINT user_to_links_pkey PRIMARY KEY (link_id, user_id);

ALTER TABLE ONLY user_to_teams
    ADD CONSTRAINT user_to_teams_pkey PRIMARY KEY (user_id, team_id);

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
