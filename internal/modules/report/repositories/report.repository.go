package repositories

import (
	"movie-app-go/internal/modules/report/responses"
	"movie-app-go/internal/repository"
	"time"

	"gorm.io/gorm"
)

type ReportRepository struct {
	DB *gorm.DB
}

func NewReportRepository(db *gorm.DB) *ReportRepository {
	return &ReportRepository{DB: db}
}

func (r *ReportRepository) GetDailySalesData(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailySalesReport], error) {
	dataQuery := r.DB.Select(`
        DATE(transactions.created_at) as date,
        COUNT(tickets.id) as total_tickets,
        SUM(transactions.total_amount) as total_sales
    `).
		Table("transactions").
		Joins("JOIN tickets ON tickets.transaction_id = transactions.id").
		Where("transactions.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("transactions.payment_status = ?", "success").
		Group("DATE(transactions.created_at)").
		Order("date DESC")

	countQuery := r.DB.Select("COUNT(DISTINCT DATE(transactions.created_at))").
		Table("transactions").
		Joins("JOIN tickets ON tickets.transaction_id = transactions.id").
		Where("transactions.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("transactions.payment_status = ?", "success")

	return repository.PaginateRaw[responses.DailySalesReport](dataQuery, countQuery, page, perPage)
}

func (r *ReportRepository) GetMonthlySalesData(year int, page, perPage int) (repository.PaginationResult[responses.MonthlySalesReport], error) {
	dataQuery := r.DB.Select(`
        EXTRACT(MONTH FROM transactions.created_at) as month,
        EXTRACT(YEAR FROM transactions.created_at) as year,
        TO_CHAR(transactions.created_at, 'Month') as month_name,
        COUNT(tickets.id) as total_tickets,
        SUM(transactions.total_amount) as total_sales
    `).
		Table("transactions").
		Joins("JOIN tickets ON tickets.transaction_id = transactions.id").
		Where("EXTRACT(YEAR FROM transactions.created_at) = ?", year).
		Where("transactions.payment_status = ?", "success").
		Group("EXTRACT(MONTH FROM transactions.created_at), EXTRACT(YEAR FROM transactions.created_at), TO_CHAR(transactions.created_at, 'Month')").
		Order("month")

	countQuery := r.DB.Select("COUNT(DISTINCT EXTRACT(MONTH FROM transactions.created_at))").
		Table("transactions").
		Joins("JOIN tickets ON tickets.transaction_id = transactions.id").
		Where("EXTRACT(YEAR FROM transactions.created_at) = ?", year).
		Where("transactions.payment_status = ?", "success")

	return repository.PaginateRaw[responses.MonthlySalesReport](dataQuery, countQuery, page, perPage)
}

func (r *ReportRepository) GetDailyMovieReportData(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailyMovieReport], error) {
	dataQuery := r.DB.Select(`
        DATE(t.created_at) as date,
        m.id as movie_id,
        m.title as movie_title,
        COUNT(tk.id) as total_tickets,
        SUM(t.total_amount) as total_sales
    `).
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN movies m ON m.id = s.movie_id").
		Where("t.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("t.payment_status = ?", "success").
		Group("DATE(t.created_at), m.id, m.title").
		Having("COUNT(tk.id) > 0").
		Order("date DESC, total_sales DESC")

	countQuery := r.DB.Select("COUNT(DISTINCT (DATE(t.created_at), m.id))").
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN movies m ON m.id = s.movie_id").
		Where("t.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("t.payment_status = ?", "success").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.DailyMovieReport](dataQuery, countQuery, page, perPage)
}

func (r *ReportRepository) GetMonthlyMovieReportData(year int, page, perPage int) (repository.PaginationResult[responses.MonthlyMovieReport], error) {
	dataQuery := r.DB.Select(`
        EXTRACT(MONTH FROM t.created_at) as month,
        EXTRACT(YEAR FROM t.created_at) as year,
        TO_CHAR(t.created_at, 'Month') as month_name,
        m.id as movie_id,
        m.title as movie_title,
        COUNT(tk.id) as total_tickets,
        SUM(t.total_amount) as total_sales
    `).
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN movies m ON m.id = s.movie_id").
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Where("t.payment_status = ?", "success").
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TO_CHAR(t.created_at, 'Month'), m.id, m.title").
		Having("COUNT(tk.id) > 0").
		Order("month, total_sales DESC")

	countQuery := r.DB.Select("COUNT(DISTINCT (EXTRACT(MONTH FROM t.created_at), m.id))").
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN movies m ON m.id = s.movie_id").
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Where("t.payment_status = ?", "success").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.MonthlyMovieReport](dataQuery, countQuery, page, perPage)
}

func (r *ReportRepository) GetDailyStudioReportData(startDate, endDate time.Time, page, perPage int) (repository.PaginationResult[responses.DailyStudioReport], error) {
	dataQuery := r.DB.Select(`
        DATE(t.created_at) as date,
        st.id as studio_id,
        st.name as studio_name,
        COUNT(tk.id) as total_tickets,
        SUM(t.total_amount) as total_sales
    `).
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN studios st ON st.id = s.studio_id").
		Where("t.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("t.payment_status = ?", "success").
		Group("DATE(t.created_at), st.id, st.name").
		Having("COUNT(tk.id) > 0").
		Order("date DESC, total_sales DESC")

	countQuery := r.DB.Select("COUNT(DISTINCT (DATE(t.created_at), st.id))").
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN studios st ON st.id = s.studio_id").
		Where("t.created_at BETWEEN ? AND ?", startDate, endDate).
		Where("t.payment_status = ?", "success").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.DailyStudioReport](dataQuery, countQuery, page, perPage)
}

func (r *ReportRepository) GetMonthlyStudioReportData(year int, page, perPage int) (repository.PaginationResult[responses.MonthlyStudioReport], error) {
	dataQuery := r.DB.Select(`
        EXTRACT(MONTH FROM t.created_at) as month,
        EXTRACT(YEAR FROM t.created_at) as year,
        TO_CHAR(t.created_at, 'Month') as month_name,
        st.id as studio_id,
        st.name as studio_name,
        COUNT(tk.id) as total_tickets,
        SUM(t.total_amount) as total_sales
    `).
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN studios st ON st.id = s.studio_id").
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Where("t.payment_status = ?", "success").
		Group("EXTRACT(MONTH FROM t.created_at), EXTRACT(YEAR FROM t.created_at), TO_CHAR(t.created_at, 'Month'), st.id, st.name").
		Having("COUNT(tk.id) > 0").
		Order("month, total_sales DESC")

	countQuery := r.DB.Select("COUNT(DISTINCT (EXTRACT(MONTH FROM t.created_at), st.id))").
		Table("transactions t").
		Joins("JOIN tickets tk ON tk.transaction_id = t.id").
		Joins("JOIN schedules s ON s.id = tk.schedule_id").
		Joins("JOIN studios st ON st.id = s.studio_id").
		Where("EXTRACT(YEAR FROM t.created_at) = ?", year).
		Where("t.payment_status = ?", "success").
		Having("COUNT(tk.id) > 0")

	return repository.PaginateRaw[responses.MonthlyStudioReport](dataQuery, countQuery, page, perPage)
}
