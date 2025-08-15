package seed

import (
    "gorm.io/gorm"
    "log"
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

    // Seed Users (jika ingin sekalian)
    if err := SeedUsers(db); err != nil {
        log.Println("SeedUsers error:", err)
        return err
    }

    return nil
}