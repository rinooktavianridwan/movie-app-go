DROP INDEX IF EXISTS idx_schedules_studio_time;

DROP INDEX IF EXISTS idx_schedules_movie_id;
DROP INDEX IF EXISTS idx_schedules_studio_id;
DROP INDEX IF EXISTS idx_schedules_date;
DROP INDEX IF EXISTS idx_schedules_start_time;

ALTER TABLE schedules DROP CONSTRAINT IF EXISTS schedules_movie_id_fkey;
ALTER TABLE schedules DROP CONSTRAINT IF EXISTS schedules_studio_id_fkey;

ALTER TABLE schedules ADD CONSTRAINT schedules_movie_id_fkey 
    FOREIGN KEY (movie_id) REFERENCES movies(id);

ALTER TABLE schedules ADD CONSTRAINT schedules_studio_id_fkey 
    FOREIGN KEY (studio_id) REFERENCES studios(id);