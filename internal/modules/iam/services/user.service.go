package services

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"movie-app-go/internal/models"
	"movie-app-go/internal/modules/iam/repositories"
	"movie-app-go/internal/modules/iam/requests"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	excelExport "movie-app-go/internal/utils/export/excel"
	excelUtil "movie-app-go/internal/utils/import/excel"
	opt "movie-app-go/internal/utils/import/excel/options"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	UserRepo *repositories.UserRepository
	AuthRepo *repositories.AuthRepository
	RoleRepo *repositories.RoleRepository
}

func NewUserService(userRepo *repositories.UserRepository, authRepo *repositories.AuthRepository, roleRepo *repositories.RoleRepository) *UserService {
	return &UserService{
		UserRepo: userRepo,
		AuthRepo: authRepo,
		RoleRepo: roleRepo,
	}
}

func (s *UserService) GetAllPaginated(page, perPage int) (repository.PaginationResult[models.User], error) {
	return s.UserRepo.GetAllPaginated(page, perPage)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	user, err := s.UserRepo.GetByID(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, utils.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (s *UserService) Update(id uint, req *requests.UserUpdateRequest) error {
	user, err := s.GetByID(id)
	if err != nil {
		return err
	}

	if req.Email != user.Email {
		exists, err := s.UserRepo.ExistsByEmailExceptID(req.Email, id)
		if err != nil {
			return err
		}
		if exists {
			return utils.ErrEmailAlreadyExists
		}
	}

	user.Name = req.Name
	user.Email = req.Email
	if req.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashed)
	}
	if req.RoleID != nil {
		user.RoleID = req.RoleID
	}

	return s.UserRepo.Update(user)
}

func (s *UserService) Delete(id uint) error {
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	return s.UserRepo.Delete(id)
}

func (s *UserService) UpdateAvatar(userID uint, file *multipart.FileHeader) error {
	user, err := s.GetByID(userID)
	if err != nil {
		return err
	}

	if user.Avatar != nil && *user.Avatar != "" {
		utils.DeleteFile(*user.Avatar)
	}

	avatarPath, err := utils.SaveFile(file, "uploads/avatars", "image", 5)
	if err != nil {
		return err
	}

	relativePath := strings.TrimPrefix(avatarPath, "./")
	user.Avatar = &relativePath

	if err := s.UserRepo.Update(user); err != nil {
		utils.DeleteFile(avatarPath)
		return err
	}

	return nil
}

func (s *UserService) ImportExcelSingleSheet(fileHeader *multipart.FileHeader, sheetName *string) error {
	headerRow := 1
	templatePath := filepath.Join("internal", "modules", "iam", "templates", "User.xlsx")

	f, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file: %w", err)
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("gagal membaca file upload: %w", err)
	}

	sheetForHeader := ""
    if sheetName != nil {
        sheetForHeader = *sheetName
    }
	headersSlice, err := excelUtil.GetHeaders(bytes.NewReader(buf), sheetForHeader, headerRow)
	if err != nil {
		return fmt.Errorf("gagal ambil header dari file: %w", err)
	}

	rows, err := excelUtil.ParseSheetToMapsWithHeader(bytes.NewReader(buf),
		opt.WithHeaderRow(headerRow),
		opt.WithAllSheet(false),
		opt.WithSheetName(sheetForHeader),
		opt.WithValidationPath(templatePath),
	)
	if err != nil {
		return fmt.Errorf("gagal parse file: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("tidak ada data untuk diimport")
	}

	var toInsert []models.User
	for _, r := range rows {
		name := strings.TrimSpace(r[headersSlice[0]])
		email := ""
		if len(headersSlice) > 1 {
			email = strings.TrimSpace(r[headersSlice[1]])
		}
		password := ""
		if len(headersSlice) > 2 {
			password = strings.TrimSpace(r[headersSlice[2]])
		}
		roleName := ""
		if len(headersSlice) > 3 {
			roleName = strings.TrimSpace(r[headersSlice[3]])
		}

		if name == "" || password == "" || email == "" || roleName == "" {
			continue
		}

		if existing, _ := s.AuthRepo.GetUserByEmail(email); existing != nil {
			continue
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("gagal hash password: %w", err)
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashed),
		}

		if roleName != "" && s.RoleRepo != nil {
			if role, err := s.RoleRepo.GetByName(roleName); err == nil && role != nil {
				user.RoleID = &role.ID
			}
		}

		toInsert = append(toInsert, user)
	}

	if len(toInsert) == 0 {
		return fmt.Errorf("tidak ada record valid untuk disimpan")
	}

	if err := s.UserRepo.CreateBatch(toInsert); err != nil {
		return fmt.Errorf("gagal insert users: %w", err)
	}
	return nil
}

func (s *UserService) ImportExcelMultiSheet(fileHeader *multipart.FileHeader) error {
	headerRow := 1
	templatePath := filepath.Join("internal", "modules", "iam", "templates", "User.xlsx")

	f, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("gagal membuka file: %w", err)
	}
	defer f.Close()
	buf, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("gagal membaca file upload: %w", err)
	}

	rows, err := excelUtil.ParseSheetToMapsWithHeader(bytes.NewReader(buf),
		opt.WithHeaderRow(headerRow),
		opt.WithAllSheet(true),
		opt.WithValidationPath(templatePath),
	)
	if err != nil {
		return fmt.Errorf("gagal parse file: %w", err)
	}

	headersSlice, err := excelUtil.GetHeaders(bytes.NewReader(buf), "", headerRow)
	if err != nil {
		return fmt.Errorf("gagal ambil header dari file: %w", err)
	}

	if len(rows) == 0 {
		return fmt.Errorf("tidak ada data untuk diimport")
	}

	var toInsert []models.User
	for _, r := range rows {
		name := strings.TrimSpace(r[headersSlice[0]])
		email := ""
		if len(headersSlice) > 1 {
			email = strings.TrimSpace(r[headersSlice[1]])
		}
		password := ""
		if len(headersSlice) > 2 {
			password = strings.TrimSpace(r[headersSlice[2]])
		}
		roleName := ""
		if len(headersSlice) > 3 {
			roleName = strings.TrimSpace(r[headersSlice[3]])
		}

		if name == "" || password == "" || email == "" || roleName == "" {
			continue
		}

		if existing, _ := s.AuthRepo.GetUserByEmail(email); existing != nil {
			return fmt.Errorf(`email %s sudah terdaftar`, email)
		}

		hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("gagal hash password: %w", err)
		}

		user := models.User{
			Name:     name,
			Email:    email,
			Password: string(hashed),
		}

		if roleName != "" && s.RoleRepo != nil {
			if role, err := s.RoleRepo.GetByName(roleName); err == nil && role != nil {
				user.RoleID = &role.ID
			}
		}
		toInsert = append(toInsert, user)
	}

	if len(toInsert) == 0 {
		return fmt.Errorf("tidak ada record valid untuk disimpan")
	}

	if err := s.UserRepo.CreateBatch(toInsert); err != nil {
		return fmt.Errorf("gagal insert users: %w", err)
	}
	return nil
}

func (s *UserService) ExportUsersExcel() ([]byte, error) {
	perPage := 50
	page := 1
	var records []map[string]string

	for {
		result, err := s.UserRepo.GetAllPaginated(page, perPage)
		if err != nil {
			return nil, err
		}

		for _, u := range result.Data {
			roleName := ""
			if u.Role != nil {
				roleName = u.Role.Name
			}
			rec := map[string]string{
				"Fullname": u.Name,
				"Email":    u.Email,
				"Roles":    roleName,
			}
			records = append(records, rec)
		}

		if len(result.Data) < perPage {
			break
		}
		page++
	}

	headers := []string{"Fullname", "Email", "Roles"}
	return excelExport.GenerateFromMaps("Users", headers, records)
}
