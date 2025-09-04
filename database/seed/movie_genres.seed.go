package seed

import (
	"fmt"
	"log"
	"math/rand"
	"movie-app-go/internal/models"

	"gorm.io/gorm"
)

func SeedMovieGenres(db *gorm.DB, movies []models.Movie, genres []models.Genre) error {

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

	var allMovieGenres []models.MovieGenre
	genreMap := make(map[string]uint)
	for _, genre := range genres {
		genreMap[genre.Name] = genre.ID
	}

	for _, movie := range movies {
		if genreNames, exists := movieGenreMapping[movie.Title]; exists {
			for _, genreName := range genreNames {
				if genreID, genreExists := genreMap[genreName]; genreExists {
					allMovieGenres = append(allMovieGenres, models.MovieGenre{
						MovieID: movie.ID,
						GenreID: genreID,
					})
				}
			}
		} else {
			numGenres := rand.Intn(3) + 1
			selectedGenres := rand.Perm(len(genres))[:numGenres]
			for _, genreIndex := range selectedGenres {
				allMovieGenres = append(allMovieGenres, models.MovieGenre{
					MovieID: movie.ID,
					GenreID: genres[genreIndex].ID,
				})
			}
		}
	}

	var existingMovieGenres []models.MovieGenre
	if err := db.Find(&existingMovieGenres).Error; err != nil {
		return err
	}

	existingMap := make(map[string]bool)
	for _, existing := range existingMovieGenres {
		key := fmt.Sprintf("%d_%d", existing.MovieID, existing.GenreID)
		existingMap[key] = true
	}

	var newMovieGenres []models.MovieGenre
	for _, mg := range allMovieGenres {
		key := fmt.Sprintf("%d_%d", mg.MovieID, mg.GenreID)
		if !existingMap[key] {
			newMovieGenres = append(newMovieGenres, mg)
		}
	}

	if len(newMovieGenres) > 0 {
		if err := db.Create(&newMovieGenres).Error; err != nil {
			return err
		}
		log.Printf("Successfully created %d new movie genre relationships", len(newMovieGenres))
	}

	log.Println("Movie genres seeding completed")
	return nil
}
