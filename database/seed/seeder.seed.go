package seed

import (
	"log"

	"gorm.io/gorm"
)

func RunAllSeeders(db *gorm.DB) error {
	// Seed Facilities
	facilities, err := SeedFacilities(db)
	if err != nil {
		log.Println("SeedFacilities error:", err)
		return err
	}

	// Seed Studios
	studios, err := SeedStudios(db)
	if err != nil {
		log.Println("SeedStudios error:", err)
		return err
	}

	// Seed FacilityStudios (relasi)
	if err := SeedFacilityStudios(db, studios, facilities); err != nil {
		log.Println("SeedFacilityStudios error:", err)
		return err
	}

	// Seed Users
	if err := SeedUsers(db); err != nil {
		log.Println("SeedUsers error:", err)
		return err
	}

	// Seed Genres
	genres, err := SeedGenres(db)
	if err != nil {
		log.Println("SeedGenres error:", err)
		return err
	}

	// Seed Movies
	movies, err := SeedMovies(db)
	if err != nil {
		log.Println("SeedMovies error:", err)
		return err
	}

	// Seed MovieGenres (relasi)
	if err := SeedMovieGenres(db, movies, genres); err != nil {
		log.Println("SeedMovieGenres error:", err)
		return err
	}

	log.Println("All seeders completed successfully!")
	return nil
}
