package quizlet

// Base represents Quizlet's base URL
const Base = "https://quizlet.com"

// Constant Quizlet API endpoints
const (
	EndpointLogin    = "/login"
	EndpointSessions = "/sessions"
)

// Dynamic Quizlet API endpoints
var (
	EndpointStudy          = func(id string) string { return "/" + id }
	EndpointStudyMode      = func(id, mode string) string { return EndpointStudy(id) + "/" + mode }
	EndpointStudyHighScore = func(id, mode string) string { return EndpointStudyMode(id, mode) + "/highscores" }
	EndpointStudyEndGame   = func(id, mode string) string { return EndpointStudyMode(id, mode) + "/endgame" }
)

// StudyMode represents the ID of a study mode
type StudyMode int

// Study modes
const (
	StudyModeWrite StudyMode = 1
	StudyModeSpell StudyMode = 8
)
