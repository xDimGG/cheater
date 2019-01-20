package quizlet

// Base represents Quizlet's base URL
const Base = "https://quizlet.com"

// Constant Quizlet API endpoints
const (
	EndpointLogin = "/login"
)

// Dynamic Quizlet API endpoints
var (
	EndpointStudy = func(id string) string { return "/" + id }

	EndpointMatch          = func(id string) string { return EndpointStudy(id) + "/match" }
	EndpointMatchHighScore = func(id string) string { return EndpointMatch(id) + "/highscores" }
)
