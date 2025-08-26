CREATE TABLE promo_movies (
    id SERIAL PRIMARY KEY,
    promo_id INTEGER NOT NULL REFERENCES promos(id) ON DELETE CASCADE,
    movie_id INTEGER NOT NULL REFERENCES movies(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(promo_id, movie_id)
);

CREATE INDEX idx_promo_movies_promo_id ON promo_movies(promo_id);
CREATE INDEX idx_promo_movies_movie_id ON promo_movies(movie_id);