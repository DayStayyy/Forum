package structs

import "time"

type User struct {
	ID, Perms                                                            int
	Name, Password, Mail, Image, Nationality, Description, RememberToken string
}

type Post struct {
	ID, UserID     int
	Title, Content string
	Categories     []string
	PostedTime     time.Time
}

type FullPost struct {
	Likes, Dislikes int
	ID, UserID      int
	Title, Content  string
	Categories      []string
	Date, Owner     string
}

type Comment struct {
	ID, UserID, ParentID, PostID int
	Content                      string
	PostedTime                   time.Time
}

type FullComment struct {
	ID, UserID, ParentID, PostID int
	Content, Owner, Date         string
	PostedTime                   time.Time
}

type Category struct {
	ID          int
	Name, Color string
}

type Like struct {
	ID, UserID, PostID, CommentID, IsPositive int
	PostedTime                                time.Time
}

type ConnectInfos struct {
	Name, Password string
}

// RETURN STRUCTURES FOR HANDLER

type HomePageReturn struct {
	Post            FullPost
	Name            string
	Likes, Dislikes int
}

type RegistrationErrors struct {
	Name, Mail, Passwd string
}
