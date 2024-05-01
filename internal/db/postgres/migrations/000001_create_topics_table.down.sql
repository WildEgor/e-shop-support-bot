-- Удаляем ограничение only_one_active_topic_per_author
ALTER TABLE public.topics DROP CONSTRAINT only_one_active_topic_per_author;

-- Удаляем индекс idx_author_status
DROP INDEX IF EXISTS idx_author_status;

-- Удаляем столбец status и восстанавливаем его как целочисленный тип
ALTER TABLE public.topics DROP COLUMN status;
ALTER TABLE public.topics ADD COLUMN status INT DEFAULT 1;

-- Удаляем столбцы feedback_id, author_tun, support_tid, support_tun
ALTER TABLE public.topics DROP COLUMN feedback_id;
ALTER TABLE public.topics DROP COLUMN author_tun;
ALTER TABLE public.topics DROP COLUMN support_tid;
ALTER TABLE public.topics DROP COLUMN support_tun;

-- Удаляем ограничение valid_status
ALTER TABLE public.topics DROP CONSTRAINT valid_status;

-- Возвращаем столбец id к его первоначальному типу SERIAL
ALTER TABLE public.topics DROP COLUMN id;
ALTER TABLE public.topics ADD COLUMN id SERIAL PRIMARY KEY;

-- Удаляем ограничение msg_tid
ALTER TABLE public.topics DROP CONSTRAINT IF EXISTS topics_msg_tid_fkey;
