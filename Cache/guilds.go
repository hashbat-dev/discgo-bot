package cache

type Guild struct {
	DbID      int
	DiscordID string
	Name      string
}

type GuildPermissions struct {
	CommandType  int
	RequiredRole string
}

const (
	CommandTypeAdmin   = iota
	CommandTypeBang    = iota
	CommandTypeTrigger = iota
)
