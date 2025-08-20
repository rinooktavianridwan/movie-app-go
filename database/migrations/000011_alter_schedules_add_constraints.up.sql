ALTER TABLE schedules DROP CONSTRAINT schedules_movie_id_fkey;
ALTER TABLE schedules DROP CONSTRAINT schedules_studio_id_fkey;

ALTER TABLE schedules ADD CONSTRAINT schedules_movie_id_fkey 
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE;

ALTER TABLE schedules ADD CONSTRAINT schedules_studio_id_fkey 
    FOREIGN KEY (studio_id) REFERENCES studios(id) ON DELETE CASCADE;

CREATE INDEX IF NOT EXISTS idx_schedules_movie_id ON schedules(movie_id);
CREATE INDEX IF NOT EXISTS idx_schedules_studio_id ON schedules(studio_id);
CREATE INDEX IF NOT EXISTS idx_schedules_date ON schedules(date);
CREATE INDEX IF NOT EXISTS idx_schedules_start_time ON schedules(start_time);

CREATE UNIQUE INDEX IF NOT EXISTS idx_schedules_studio_time 
    ON schedules(studio_id, start_time, end_time) 
    WHERE deleted_at IS NULL;