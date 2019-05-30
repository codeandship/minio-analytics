package minioanalytics

// Storage is the main storage interface
type Storage interface {
	Open() error
	Close() error
	ListAnalytics() ([]Analytics, error)
}
