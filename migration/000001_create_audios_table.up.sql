CREATE TABLE public.audios
(
    uuid  UUID NOT NULL PRIMARY KEY DEFAULT GEN_RANDOM_UUID() ,
    "group" TEXT NOT NULL ,
    song TEXT NOT NULL ,
    release_date DATE NOT NULL ,
    link TEXT NOT NULL ,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP(3) NOT NULL,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP(3) NOT NULL
);

CREATE INDEX idx_realise_group
    ON public.audios ("group");

CREATE INDEX idx_audios_song_fulltext
    ON public.audios
        USING GIN (to_tsvector('english', song));

CREATE INDEX idx_realise_date
    ON public.audios (release_date);

CREATE INDEX idx_realise_link
    ON public.audios (link);