CREATE TABLE public.lyrics
(
    uuid  UUID NOT NULL PRIMARY KEY DEFAULT GEN_RANDOM_UUID() ,
    audio_uuid UUID NOT NULL ,
    "order" INTEGER NOT NULL ,
    text TEXT NOT NULL ,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP(3) NOT NULL ,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP(3) NOT NULL ,
    FOREIGN KEY (audio_uuid) REFERENCES audios(uuid) ON DELETE CASCADE
);

CREATE INDEX idx_lyrics_text_fulltext
    ON public.lyrics
    USING GIN (to_tsvector('english', text));