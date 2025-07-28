package validation

type FileValidator interface {
	IsReachable(url string) bool
}
