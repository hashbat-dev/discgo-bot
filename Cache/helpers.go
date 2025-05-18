package cache

func GetUserIDFromCorrelationID(correlationId string) string {
	return ActiveInteractions[correlationId].StartInteraction.Member.User.ID
}
