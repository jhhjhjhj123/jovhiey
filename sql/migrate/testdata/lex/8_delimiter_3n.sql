-- atlas:delimiter \n\n\n

CREATE EXTENSION IF NOT EXISTS unaccent;


CREATE OR REPLACE FUNCTION public.slugify(
    v TEXT
) RETURNS TEXT STRICT IMMUTABLE AS $$
BEGIN
    RETURN trim(BOTH '-' FROM regexp_replace(lower(unaccent(trim(v))), '[^a-z0-9\\-_]+', '-', 'gi'));
END;
LANGUAGE plpgsql;
