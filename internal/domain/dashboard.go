package domain

type DashboardStats struct {
	TotalCustomers     int
	TotalBookings      int
	PendingBookings    int
	ActiveDeliveries   int
	AvailableDrivers   int
	PendingQuotes      int
	NewContactRequests int
	PendingReviews     int
	RevenueToday       float64
}
