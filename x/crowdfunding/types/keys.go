package types

const (
	// ModuleName is the name of the crowdfunding module
	ModuleName = "crowdfunding"

	// StoreKey is the string store key for the crowdfunding module
	StoreKey = ModuleName

	// RouterKey is the message route for the crowdfunding module
	RouterKey = ModuleName

	// QuerierRoute is the querier route for the crowdfunding module
	QuerierRoute = ModuleName
)

// Key prefixes for store
var (
	// IdeaKey is the prefix for idea store
	IdeaKey = []byte{0x01}

	// CampaignKey is the prefix for campaign store
	CampaignKey = []byte{0x02}

	// ContributionKey is the prefix for contribution store
	ContributionKey = []byte{0x03}

	// BackerKey is the prefix for backer store
	BackerKey = []byte{0x04}
)

// GetIdeaKey returns the key for an idea by ID
func GetIdeaKey(id string) []byte {
	return append(IdeaKey, []byte(id)...)
}

// GetCampaignKey returns the key for a campaign by ID
func GetCampaignKey(id string) []byte {
	return append(CampaignKey, []byte(id)...)
}

// GetContributionKey returns the key for a contribution by ID
func GetContributionKey(id string) []byte {
	return append(ContributionKey, []byte(id)...)
}

// GetBackerKey returns the key for a backer by address
func GetBackerKey(address string) []byte {
	return append(BackerKey, []byte(address)...)
}
