package seed

import (
    "movie-app-go/internal/models"
    "gorm.io/gorm"
)

func SeedFacilityStudios(db *gorm.DB, studios []models.Studio, facilities []models.Facility) error {
    // Contoh: Studio 1 punya semua fasilitas, Studio 2 hanya 2, Studio 3 hanya 1
    var relasi []models.FacilityStudio
    if len(studios) > 0 && len(facilities) > 0 {
        // Studio 1: semua fasilitas
        for _, f := range facilities {
            relasi = append(relasi, models.FacilityStudio{
                StudioID:   studios[0].ID,
                FacilityID: f.ID,
            })
        }
        // Studio 2: fasilitas 1 dan 2
        if len(facilities) > 1 {
            relasi = append(relasi, models.FacilityStudio{
                StudioID:   studios[1].ID,
                FacilityID: facilities[0].ID,
            }, models.FacilityStudio{
                StudioID:   studios[1].ID,
                FacilityID: facilities[1].ID,
            })
        }
        // Studio 3: fasilitas 3 saja
        if len(facilities) > 2 {
            relasi = append(relasi, models.FacilityStudio{
                StudioID:   studios[2].ID,
                FacilityID: facilities[2].ID,
            })
        }
    }
    if len(relasi) > 0 {
        return db.Create(&relasi).Error
    }
    return nil
}