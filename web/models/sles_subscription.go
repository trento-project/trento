package models

type SlesSubscription struct {
	ID                 string
	Version            string
	Type               string
	Arch               string
	Status             string
	StartsAt           string
	ExpiresAt          string
	SubscriptionStatus string
}

type PremiumData struct {
	IsPremium     bool
	Sles4SapCount int
}
