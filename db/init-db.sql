CREATE DATABASE greenlight;

\c greenlight;

CREATE ROLE greenlight WITH LOGIN PASSWORD 'pa55word';
CREATE EXTENSION IF NOT EXISTS citext;

\c postgres;

GRANT CREATE ON DATABASE greenlight TO greenlight;

ALTER DATABASE greenlight OWNER TO greenlight;
