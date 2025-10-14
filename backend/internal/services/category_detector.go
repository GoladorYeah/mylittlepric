package services

type CategoryDetector struct{}

func NewCategoryDetector() *CategoryDetector {
	return &CategoryDetector{}
}

func (cd *CategoryDetector) DetectCategory(message string) string {
	return ""
}
