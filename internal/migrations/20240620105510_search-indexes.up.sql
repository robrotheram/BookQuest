-- Up Migration: Create Indexes

-- Index on links.id
CREATE INDEX idx_links_id ON links(id);

-- Index on favourite_links.link_id and favourite_links.user_id
CREATE INDEX idx_favourite_links_link_id_user_id ON favourite_links(link_id, user_id);

-- Index on user_to_links.link_id
CREATE INDEX idx_user_to_links_link_id ON user_to_links(link_id);

-- Index on team_links.link_id
CREATE INDEX idx_team_links_link_id ON team_links(link_id);

-- Index on teams.id
CREATE INDEX idx_teams_id ON teams(id);

-- Index on teams.visability
CREATE INDEX idx_teams_visability ON teams(visability);

-- Index on user_to_teams.user_id
CREATE INDEX idx_user_to_teams_user_id ON user_to_teams(user_id);

-- Index on links.sharing
CREATE INDEX idx_links_sharing ON links(sharing);

-- Composite index on user_to_links.user_id and user_to_links.link_id
CREATE INDEX idx_user_to_links_user_id_link_id ON user_to_links(user_id, link_id);

-- Composite index on team_links.team_id and team_links.link_id
CREATE INDEX idx_team_links_team_id_link_id ON team_links(team_id, link_id);
