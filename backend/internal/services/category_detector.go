package services

import (
	"strings"
)

type CategoryDetector struct{}

func NewCategoryDetector() *CategoryDetector {
	return &CategoryDetector{}
}

func (cd *CategoryDetector) DetectCategory(message string) string {
	messageLower := strings.ToLower(message)

	if cd.isElectronics(messageLower) {
		return "electronics"
	}

	if cd.isGenericModel(messageLower) {
		return "generic_model"
	}

	if cd.isClothing(messageLower) {
		return "clothing"
	}

	if cd.isFurniture(messageLower) {
		return "furniture"
	}

	if cd.isKitchen(messageLower) {
		return "kitchen"
	}

	if cd.isSports(messageLower) {
		return "sports"
	}

	if cd.isTools(messageLower) {
		return "tools"
	}

	if cd.isDecor(messageLower) {
		return "decor"
	}

	if cd.isTextiles(messageLower) {
		return "textiles"
	}

	return ""
}

func (cd *CategoryDetector) isElectronics(msg string) bool {
	keywords := []string{
		"iphone", "ipad", "macbook", "apple watch", "airpods",
		"samsung galaxy", "galaxy", "pixel", "google pixel",
		"xiaomi", "redmi", "oneplus", "realme", "oppo",
		"laptop", "notebook", "computer", "pc",
		"tablet", "phone", "smartphone", "smartwatch",
		"headphones", "earbuds", "speaker", "bluetooth",
		"tv", "television", "monitor", "display",
		"camera", "canon", "nikon", "sony alpha",
		"playstation", "ps5", "xbox", "nintendo switch",
		"console", "gaming",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isGenericModel(msg string) bool {
	keywords := []string{
		"battery", "batteries", "aa", "aaa", "18650", "cr2032",
		"tire", "tires", "tyre", "195/65r15", "205/55r16",
		"bulb", "lamp", "e27", "e14", "gu10", "led bulb",
		"paint", "ral", "ral 9016", "ncs",
		"arduino", "raspberry pi", "esp32",
		"filter", "oil filter", "air filter", "water filter",
		"cartridge", "ink", "toner",
		"spark plug", "brake pad", "timing belt",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isClothing(msg string) bool {
	keywords := []string{
		"jacket", "coat", "hoodie", "sweater",
		"shirt", "t-shirt", "polo", "blouse",
		"pants", "jeans", "trousers", "shorts",
		"dress", "skirt", "suit",
		"shoes", "boots", "sneakers", "sandals",
		"hat", "cap", "scarf", "gloves",
		"underwear", "socks", "belt", "tie",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isFurniture(msg string) bool {
	keywords := []string{
		"sofa", "couch", "armchair", "chair",
		"table", "desk", "dining table", "coffee table",
		"bed", "mattress", "wardrobe", "closet",
		"shelf", "bookshelf", "cabinet", "drawer",
		"stool", "bench", "ottoman",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isKitchen(msg string) bool {
	keywords := []string{
		"pan", "pot", "frying pan", "saucepan",
		"knife", "cutting board", "chopping board",
		"blender", "mixer", "food processor",
		"coffee maker", "espresso", "kettle",
		"toaster", "microwave", "oven",
		"dishwasher", "refrigerator", "fridge",
		"plate", "bowl", "cup", "mug", "glass",
		"cutlery", "fork", "spoon",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isSports(msg string) bool {
	keywords := []string{
		"treadmill", "exercise bike", "elliptical",
		"dumbbell", "barbell", "weight", "kettlebell",
		"yoga mat", "fitness mat", "resistance band",
		"skates", "rollerblades", "skateboard",
		"bicycle", "bike", "mountain bike",
		"tennis racket", "badminton", "basketball",
		"football", "soccer ball", "volleyball",
		"fitness tracker", "running shoes", "gym bag",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isTools(msg string) bool {
	keywords := []string{
		"drill", "screwdriver", "hammer", "saw",
		"wrench", "pliers", "level", "tape measure",
		"power tool", "cordless", "electric",
		"ladder", "toolbox", "workbench",
		"sander", "grinder", "router",
		"nail gun", "stapler", "compressor",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isDecor(msg string) bool {
	keywords := []string{
		"lamp", "floor lamp", "table lamp", "desk lamp",
		"mirror", "wall mirror", "standing mirror",
		"vase", "flower pot", "planter",
		"picture frame", "wall art", "painting",
		"candle", "candleholder", "lantern",
		"rug", "carpet", "doormat",
		"curtain", "blinds", "shade",
		"clock", "wall clock", "alarm clock",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}

func (cd *CategoryDetector) isTextiles(msg string) bool {
	keywords := []string{
		"pillow", "cushion", "throw pillow",
		"blanket", "throw", "duvet", "comforter",
		"bedding", "bed sheet", "fitted sheet",
		"pillowcase", "duvet cover",
		"towel", "bath towel", "hand towel",
		"bathrobe", "bathroom mat",
		"tablecloth", "napkin", "placemat",
	}

	for _, keyword := range keywords {
		if strings.Contains(msg, keyword) {
			return true
		}
	}

	return false
}
