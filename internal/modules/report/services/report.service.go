package services

import (
	"movie-app-go/internal/enums"
	"movie-app-go/internal/modules/report/responses"
	"movie-app-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type ReportService struct {
	DB *gorm.DB
}

func NewReportService(db *gorm.DB) *ReportService {
	return &ReportService{DB: db}
}

func (s *ReportService) GetDailySalesReport(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailySalesReport], error) {
	dataQuery := s.DB.Table("transactions t").
		Select("DATE(t.created_at) as date, COALESCE(SUM(tk.price), 0) as total_revenue, COUNT(tk.id) as total_tickets").
		Joins("LEFT JOIN tickets tk ON t.id = tk.transaction_id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("DATE(t.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("DATE(t.created_at)").
		Order("date DESC")

	countQuery := s.DB.Table("transactions t").
		Joins("LEFT JOIN tickets tk ON t.id = tk.transaction_id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("DATE(t.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("DATE(t.created_at)")

	return repository.PaginateRaw[responses.DailySalesReport](dataQuery, countQuery, page, perPage)
}

func (s *ReportService) GetMonthlySalesReport(year int, page, perPage int) (repository.PaginationResult[responses.MonthlySalesReport], error) {
	dataQuery := s.DB.Table("transactions t").
		Select("TRIM(TO_CHAR(t.created_at, 'Month')) as month, EXTRACT(YEAR FROM t.created_at) as year, COALESCE(SUM(tk.price), 0) as total_revenue, COUNT(tk.id) as total_tickets").
		Joins("LEFT JOIN tickets tk ON t.id = tk.transaction_id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TRIM(TO_CHAR(t.created_at, 'Month'))").
		Order("EXTRACT(MONTH FROM t.created_at)")

	countQuery := s.DB.Table("transactions t").
		Joins("LEFT JOIN tickets tk ON t.id = tk.transaction_id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TO_CHAR(t.created_at, 'Month')")

	return repository.PaginateRaw[responses.MonthlySalesReport](dataQuery, countQuery, page, perPage)
}

func (s *ReportService) GetDailyReportByMovie(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailyMovieReport], error) {
	dataQuery := s.DB.Table("movies m").
		Select("DATE(t.created_at) as date, m.id as movie_id, m.title as movie_title, COALESCE(SUM(tk.price), 0) as total_revenue, COUNT(tk.id) as total_tickets").
		Joins("LEFT JOIN schedules s ON m.id = s.movie_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("DATE(t.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("DATE(t.created_at), m.id, m.title").
		Having("COUNT(tk.id) > 0").
		Order("date DESC, total_revenue DESC")

	countQuery := s.DB.Table("movies m").
		Joins("LEFT JOIN schedules s ON m.id = s.movie_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("DATE(t.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("DATE(t.created_at), m.id, m.title").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.DailyMovieReport](dataQuery, countQuery, page, perPage)
}

func (s *ReportService) GetDailyReportByStudio(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailyStudioReport], error) {
	dataQuery := s.DB.Table("studios st").
		Select("DATE(t.created_at) as date, st.id as studio_id, st.name as studio_name, COALESCE(SUM(tk.price), 0) as total_revenue, COUNT(tk.id) as total_tickets").
		Joins("LEFT JOIN schedules s ON st.id = s.studio_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("DATE(t.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("DATE(t.created_at), st.id, st.name").
		Having("COUNT(tk.id) > 0").
		Order("date DESC, total_revenue DESC")

	countQuery := s.DB.Table("studios st").
		Joins("LEFT JOIN schedules s ON st.id = s.studio_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("DATE(t.created_at) BETWEEN ? AND ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).
		Group("DATE(t.created_at), st.id, st.name").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.DailyStudioReport](dataQuery, countQuery, page, perPage)
}

func (s *ReportService) GetMonthlyReportByMovie(year int, page, perPage int) (repository.PaginationResult[responses.MonthlyMovieReport], error) {
	dataQuery := s.DB.Table("movies m").
		Select("TRIM(TO_CHAR(t.created_at, 'Month')) as month, EXTRACT(YEAR FROM t.created_at) as year, m.id as movie_id, m.title as movie_title, COALESCE(SUM(tk.price), 0) as total_revenue, COUNT(tk.id) as total_tickets").
		Joins("LEFT JOIN schedules s ON m.id = s.movie_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TRIM(TO_CHAR(t.created_at, 'Month')), m.id, m.title").
		Having("COUNT(tk.id) > 0").
		Order("EXTRACT(MONTH FROM t.created_at), total_revenue DESC")

	countQuery := s.DB.Table("movies m").
		Joins("LEFT JOIN schedules s ON m.id = s.movie_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TO_CHAR(t.created_at, 'Month'), m.id, m.title").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.MonthlyMovieReport](dataQuery, countQuery, page, perPage)
}

func (s *ReportService) GetMonthlyReportByStudio(year int, page, perPage int) (repository.PaginationResult[responses.MonthlyStudioReport], error) {
	dataQuery := s.DB.Table("studios st").
		Select("TRIM(TO_CHAR(t.created_at, 'Month')) as month, EXTRACT(YEAR FROM t.created_at) as year, st.id as studio_id, st.name as studio_name, COALESCE(SUM(tk.price), 0) as total_revenue, COUNT(tk.id) as total_tickets").
		Joins("LEFT JOIN schedules s ON st.id = s.studio_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TRIM(TO_CHAR(t.created_at, 'Month')), st.id, st.name").
		Having("COUNT(tk.id) > 0").
		Order("EXTRACT(MONTH FROM t.created_at), total_revenue DESC")

	countQuery := s.DB.Table("studios st").
		Joins("LEFT JOIN schedules s ON st.id = s.studio_id").
		Joins("LEFT JOIN tickets tk ON s.id = tk.schedule_id").
		Joins("LEFT JOIN transactions t ON tk.transaction_id = t.id").
		Where("t.payment_status = ?", enums.PaymentStatusSuccess).
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TO_CHAR(t.created_at, 'Month'), st.id, st.name").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.MonthlyStudioReport](dataQuery, countQuery, page, perPage)
}
