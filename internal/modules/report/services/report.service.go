package services

import (
	"movie-app-go/internal/modules/report/repositories"
	"movie-app-go/internal/modules/report/responses"
	"movie-app-go/internal/repository"
	"movie-app-go/internal/utils"
	"time"
)

type ReportService struct {
	ReportRepo *repositories.ReportRepository
}

func NewReportService(reportRepo *repositories.ReportRepository) *ReportService {
	return &ReportService{
		ReportRepo: reportRepo,
	}
}

func (s *ReportService) GetDailySalesReport(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailySalesReport], error) {
	result, err := s.ReportRepo.GetDailySalesData(startDate, endDate, page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.DailySalesReport]{}, err
	}

	if len(result.Data) == 0 {
		return repository.PaginationResult[responses.DailySalesReport]{}, utils.ErrNoReportData
	}

	return result, nil
}

func (s *ReportService) GetMonthlySalesReport(year int, page, perPage int) (repository.PaginationResult[responses.MonthlySalesReport], error) {
	result, err := s.ReportRepo.GetMonthlySalesData(year, page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.MonthlySalesReport]{}, err
	}

	if len(result.Data) == 0 {
		return repository.PaginationResult[responses.MonthlySalesReport]{}, utils.ErrNoReportData
	}

	return result, nil
}

func (s *ReportService) GetDailyReportByMovie(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailyMovieReport], error) {
	result, err := s.ReportRepo.GetDailyMovieReportData(startDate, endDate, page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.DailyMovieReport]{}, err
	}

	if len(result.Data) == 0 {
		return repository.PaginationResult[responses.DailyMovieReport]{}, utils.ErrNoReportData
	}

	return result, nil
}

func (s *ReportService) GetMonthlyReportByMovie(year int, page, perPage int) (repository.PaginationResult[responses.MonthlyMovieReport], error) {
	result, err := s.ReportRepo.GetMonthlyMovieReportData(year, page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.MonthlyMovieReport]{}, err
	}

	if len(result.Data) == 0 {
		return repository.PaginationResult[responses.MonthlyMovieReport]{}, utils.ErrNoReportData
	}

	return result, nil
}

func (s *ReportService) GetDailyReportByStudio(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailyStudioReport], error) {
	result, err := s.ReportRepo.GetDailyStudioReportData(startDate, endDate, page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.DailyStudioReport]{}, err
	}

	if len(result.Data) == 0 {
		return repository.PaginationResult[responses.DailyStudioReport]{}, utils.ErrNoReportData
	}

	return result, nil
}

func (s *ReportService) GetMonthlyReportByStudio(year int, page, perPage int) (repository.PaginationResult[responses.MonthlyStudioReport], error) {
	result, err := s.ReportRepo.GetMonthlyStudioReportData(year, page, perPage)
	if err != nil {
		return repository.PaginationResult[responses.MonthlyStudioReport]{}, err
	}

	if len(result.Data) == 0 {
		return repository.PaginationResult[responses.MonthlyStudioReport]{}, utils.ErrNoReportData
	}

	return result, nil
}
