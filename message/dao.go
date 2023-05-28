package message

import (
	"facechat/db"
)

var dbIns = db.GetDb("lilac")

type CustomUser struct {
	Uuid         string
	Email        string
	Access_token string
}

type Chat struct {
	Uuid  string
	Name  string
	Users []*CustomUser
}

type Test struct {
	Uuid         string
	Access_token string
}

func (u *CustomUser) FindUserByAccessToken() (data CustomUser, err error) {
	dbIns.Raw("select u.uuid, u.access_token from authtoken_token t, user_customuser u where t.user_id = u.uuid and t.key = ?", u.Access_token).Scan(&data)
	return
}

func (c *Chat) FindAllUserInRoom() ([]CustomUser, error) {
	var users []CustomUser
	err := dbIns.Raw(`select uuid, email, access_token from user_customuser where uuid in (
		select user_id from chat_chatmembership where chat_id = ?
	)`, c.Uuid).Find(&users).Error
	return users, err

}

func (u *CustomUser) UpdateConnect() (data CustomUser, err error) {
	dbIns.Raw("update user_customuser set is_connect = true, last_login_at = now() where uuid = ? returning uuid, email", u.Uuid).Scan(&data)
	return
}

func (u *CustomUser) UpdateDisConnect() (data CustomUser, err error) {
	dbIns.Raw("update user_customuser set is_connect = false where uuid = ? returning uuid, email", u.Uuid).Scan(&data)
	return
}
