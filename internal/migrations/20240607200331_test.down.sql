ALTER TABLE ONLY public.favourite_links
    DROP CONSTRAINT favourite_links_pkey;

ALTER TABLE ONLY public.link_meta
    DROP CONSTRAINT link_meta_pkey;

ALTER TABLE ONLY public.links
    DROP CONSTRAINT links_pkey;

ALTER TABLE ONLY public.team_links
    DROP CONSTRAINT team_links_pkey;

ALTER TABLE ONLY public.teams
    DROP CONSTRAINT teams_pkey;

ALTER TABLE ONLY public.user_to_links
    DROP CONSTRAINT user_to_links_pkey;

ALTER TABLE ONLY public.user_to_teams
    DROP CONSTRAINT user_to_teams_pkey;

ALTER TABLE ONLY public.users
    DROP CONSTRAINT users_pkey;

-- Drop tables in the reverse order of their creation
DROP TABLE IF EXISTS public.favourite_links;
DROP TABLE IF EXISTS public.link_meta;
DROP TABLE IF EXISTS public.links;
DROP TABLE IF EXISTS public.team_links;
DROP TABLE IF EXISTS public.teams;
DROP TABLE IF EXISTS public.user_to_links;
DROP TABLE IF EXISTS public.user_to_teams;
DROP TABLE IF EXISTS public.users;