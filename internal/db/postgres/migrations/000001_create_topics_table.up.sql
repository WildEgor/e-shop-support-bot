CREATE TABLE IF NOT EXISTS public.topics
(
    id SERIAL PRIMARY KEY,
    feedback_id UUID DEFAULT NULL,
    author_tid INT DEFAULT 0,
    author_tun VARCHAR(255) DEFAULT '',
    support_tid INT DEFAULT 0,
    support_tun VARCHAR(255) DEFAULT '',
    msg_tid INT DEFAULT 0,
    content VARCHAR(255) DEFAULT '',
    status INT DEFAULT 1,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_author_status ON public.topics (author_tid, status);

ALTER TABLE public.topics ADD CONSTRAINT valid_status CHECK (status IN (1, 0));

ALTER TABLE public.topics ADD CONSTRAINT only_one_active_topic_per_author
    EXCLUDE (author_tid WITH =, status WITH =) WHERE (status = 1);