package bitrix

// Profile is the current Bitrix24 user.
type Profile struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	AvatarURL string `json:"avatarUrl"`
	PortalURL string `json:"portalUrl"`
	PageURL   string `json:"pageUrl"`
}

// TimeManStatus is the workday state from timeman.status.
type TimeManStatus string

const (
	TimeManUnknown TimeManStatus = ""
	TimeManOpened  TimeManStatus = "opened"
	TimeManPaused  TimeManStatus = "paused"
	TimeManClosed  TimeManStatus = "closed"
	TimeManExpired TimeManStatus = "expired"
)

// TimeMan is the current workday snapshot.
type TimeMan struct {
	Status    TimeManStatus `json:"status"`
	TimeStart string        `json:"timeStart"`
	Duration  int           `json:"duration"`
	TimeLeaks int           `json:"timeLeaks"`
}
