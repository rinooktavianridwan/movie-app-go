package utils

import "errors"

// Eentity not found errors
var (
	ErrMovieNotFound        = errors.New("movie not found")
	ErrGenreNotFound        = errors.New("genre not found")
	ErrStudioNotFound       = errors.New("studio not found")
	ErrUserNotFound         = errors.New("user not found")
	ErrFacilityNotFound     = errors.New("facility not found")
	ErrScheduleNotFound     = errors.New("schedule not found")
	ErrPromoNotFound        = errors.New("promo not found")
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrTicketNotFound       = errors.New("ticket not found")
	ErrNotificationNotFound = errors.New("notification not found")
)

var (
	ErrInvalidGenreIDs       = errors.New("some genre_ids are invalid")
	ErrInvalidFacilityIDs    = errors.New("some facility_ids are invalid")
	ErrInvalidMovieIDs       = errors.New("some movie_ids are invalid")
	ErrInvalidStudioIDs      = errors.New("studio_id is invalid")
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrFacilityAlreadyExists = errors.New("facility name already exists")
	ErrGenreAlreadyExists    = errors.New("genre name already exists")
	ErrInvalidFileType       = errors.New("invalid file type")
	ErrInvalidTimeRange      = errors.New("start_time must be before end_time")
	ErrFileSizeExceeded      = errors.New("file size exceeded maximum limit")
	ErrInvalidDateFormat     = errors.New("invalid date format, use YYYY-MM-DD")
	ErrInvalidDateRange      = errors.New("end date must be after start date")
	ErrInvalidYear           = errors.New("invalid year format")
	ErrPastDate              = errors.New("date cannot be in the past")
)

// Business logic errors
var (
	ErrScheduleConflict   = errors.New("schedule conflict: studio is already booked")
	ErrScheduleHasTickets = errors.New("cannot delete schedule: tickets already exist")
	ErrStudioHasSchedules = errors.New("cannot delete studio: schedules still exist")
	ErrMovieHasSchedules  = errors.New("cannot delete movie: schedules still exist")
	ErrFacilityInUse      = errors.New("cannot delete facility: still used by studios")
	ErrGenreInUse         = errors.New("cannot delete genre: still used by movies")
	ErrPromoInUse         = errors.New("cannot delete promo: still being used")
	ErrPromoCodeExists    = errors.New("promo code already exists")
)

// Transaction & booking errors
var (
	ErrTransactionAlreadyPaid  = errors.New("transaction already paid")
	ErrTransactionExpired      = errors.New("transaction has expired")
	ErrInsufficientSeats       = errors.New("insufficient available seats")
	ErrSeatAlreadyBooked       = errors.New("one or more seats already booked")
	ErrPaymentProcessingFailed = errors.New("payment processing failed")
	ErrTicketAlreadyScanned    = errors.New("ticket already scanned")
	ErrTicketNotPaid           = errors.New("ticket not paid yet")
	ErrTicketCancelled         = errors.New("cancelled ticket cannot be used")
	ErrUnauthorizedAccess      = errors.New("unauthorized access")
)

// Data not found errors
var (
	ErrNoReportData = errors.New("no data found for the specified period")
)
