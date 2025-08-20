package seed

import (
	"math/rand"
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedMovieGenres(db *gorm.DB, movies []models.Movie, genres []models.Genre) error {
	var movieGenres []models.MovieGenre

	// Mapping manual untuk beberapa film populer
	movieGenreMapping := map[string][]string{
		"Avengers: Endgame":        {"Action", "Adventure", "Superhero"},
		"Spider-Man: No Way Home":  {"Action", "Adventure", "Superhero"},
		"The Dark Knight":          {"Action", "Crime", "Drama", "Superhero"},
		"Inception":                {"Action", "Sci-Fi", "Thriller"},
		"The Shawshank Redemption": {"Drama"},
		"Interstellar":             {"Adventure", "Drama", "Sci-Fi"},
		"The Godfather":            {"Crime", "Drama"},
		"Pulp Fiction":             {"Crime", "Drama", "Thriller"},
		"Forrest Gump":             {"Drama", "Romance"},
		"The Matrix":               {"Action", "Sci-Fi"},
		"Goodfellas":               {"Biography", "Crime", "Drama"},
		"The Lord of the Rings: The Fellowship of the Ring": {"Adventure", "Drama", "Fantasy"},
		"Star Wars: A New Hope":                             {"Adventure", "Fantasy", "Sci-Fi"},
		"Fight Club":                                        {"Drama"},
		"The Lion King":                                     {"Animation", "Adventure", "Drama", "Family", "Musical"},
		"Toy Story":                                         {"Animation", "Adventure", "Comedy", "Family"},
		"Jurassic Park":                                     {"Adventure", "Sci-Fi", "Thriller"},
		"Titanic":                                           {"Drama", "Romance"},
		"The Silence of the Lambs":                          {"Crime", "Drama", "Horror", "Thriller"},
		"Saving Private Ryan":                               {"Drama", "War"},
		"Schindler's List":                                  {"Biography", "Drama", "War"},
		"La La Land":                                        {"Comedy", "Drama", "Musical", "Romance"},
		"Parasite":                                          {"Comedy", "Drama", "Thriller"},
		"Joker":                                             {"Crime", "Drama", "Thriller"},
		"Black Panther":                                     {"Action", "Adventure", "Superhero"},
		"Frozen":                                            {"Animation", "Adventure", "Comedy", "Family", "Musical"},
		"Finding Nemo":                                      {"Animation", "Adventure", "Comedy", "Family"},
		"The Incredibles":                                   {"Animation", "Action", "Adventure", "Comedy", "Family", "Superhero"},
		"WALL-E":                                            {"Animation", "Adventure", "Family", "Sci-Fi"},
		"Up":                                                {"Animation", "Adventure", "Comedy", "Drama", "Family"},
		"Inside Out":                                        {"Animation", "Adventure", "Comedy", "Drama", "Family"},
	}

	// Buat map untuk genre berdasarkan nama
	genreMap := make(map[string]uint)
	for _, genre := range genres {
		genreMap[genre.Name] = genre.ID
	}

	// Assign genre ke movie berdasarkan mapping
	for _, movie := range movies {
		if genreNames, exists := movieGenreMapping[movie.Title]; exists {
			for _, genreName := range genreNames {
				if genreID, genreExists := genreMap[genreName]; genreExists {
					movieGenres = append(movieGenres, models.MovieGenre{
						MovieID: movie.ID,
						GenreID: genreID,
					})
				}
			}
		} else {
			// Untuk movie yang tidak ada di mapping, assign random 1-3 genre
			numGenres := rand.Intn(3) + 1
			usedGenres := make(map[uint]bool)

			for i := 0; i < numGenres && len(usedGenres) < len(genres); i++ {
				randomGenre := genres[rand.Intn(len(genres))]
				if !usedGenres[randomGenre.ID] {
					movieGenres = append(movieGenres, models.MovieGenre{
						MovieID: movie.ID,
						GenreID: randomGenre.ID,
					})
					usedGenres[randomGenre.ID] = true
				}
			}
		}
	}

	if len(movieGenres) > 0 {
		return db.Create(&movieGenres).Error
	}
	return nil
}
