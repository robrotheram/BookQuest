-- Down Migration: Drop Indexes

-- Drop index on links.id
DROP INDEX IF EXISTS idx_links_id;

-- Drop index on favourite_links.link_id and favourite_links.user_id
DROP INDEX IF EXISTS idx_favourite_links_link_id_user_id;

-- Drop index on user_to_links.link_id
DROP INDEX IF EXISTS idx_user_to_links_link_id;

-- Drop index on team_links.link_id
DROP INDEX IF EXISTS idx_team_links_link_id;

-- Drop index on teams.id
DROP INDEX IF EXISTS idx_teams_id;

-- Drop index on teams.visability
DROP INDEX IF EXISTS idx_teams_visability;

-- Drop index on user_to_teams.user_id
DROP INDEX IF EXISTS idx_user_to_teams_user_id;

-- Drop index on links.sharing
DROP INDEX IF EXISTS idx_links_sharing;

-- Drop composite index on user_to_links.user_id and user_to_links.link_id
DROP INDEX IF EXISTS idx_user_to_links_user_id_link_id;

-- Drop composite index on team_links.team_id and team_links.link_id
DROP INDEX IF EXISTS idx_team_links_team_id_link_id;
