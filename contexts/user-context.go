package contexts

type UserContext struct {
	RouteContext

	UserID    string
	UserRefID string
}
