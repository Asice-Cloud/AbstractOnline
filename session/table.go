package session

type Table struct {
	UserSession
	AdminSession
	VisitorSession
}

type UserSession struct {
	UserID      int
	UserName    string
	AccessToken string
}

type AdminSession struct {
	ID          int
	AdminName   string
	AccessToken string
}

type VisitorSession struct {
	//session_name string
	ID string
}
