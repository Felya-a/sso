package auth_service

type UserModel struct {
	ID       int64  `db:"id"`
	Email    string `db:"email"`
	PassHash []byte `db:"password"`
}
