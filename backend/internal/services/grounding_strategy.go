package services

type GroundingStrategy struct{}

func NewGroundingStrategy(mode string) *GroundingStrategy {
	return &GroundingStrategy{}
}
