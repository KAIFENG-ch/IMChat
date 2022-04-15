package dao

import (
	"IMChat/model"
	"IMChat/serialize"
	"IMChat/utils"
	"golang.org/x/crypto/bcrypt"
	"log"
	"mime/multipart"
	"strconv"
)

type UserRegister struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

type UserUpdate struct {
	Username  string `form:"username"`
	Gender    string `form:"gender"`
	Email     string `form:"email"`
	Age       uint   `form:"age"`
	Birthday  string `form:"birthday"`
	Signature string `form:"signature"`
}

type GroupRegister struct {
	Name string `form:"name"`
}

func (u UserRegister) Register() *serialize.Base {
	var user model.User
	hashPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	user = model.User{
		Name:     u.Username,
		Password: string(hashPassword),
	}
	model.DB.Create(&user)
	return &serialize.Base{
		Status: 200,
		Msg:    "ok",
		Data:   user.Name,
	}
}

func (u UserRegister) Login() *serialize.Base {
	var user model.User
	model.DB.Model(&model.User{}).Where("name = ?", u.Username).First(&user)
	if user.ID == 0 {
		return &serialize.Base{
			Status: 200,
			Msg:    "failed",
			Data:   "登录失败！",
		}
	}
	result := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password))
	if result != nil {
		return &serialize.Base{
			Status: 200,
			Msg:    "failed",
			Data:   "密码错误！",
		}
	}
	token, err := utils.CreateToken(user)
	if err != nil {
		panic(err)
	}
	return &serialize.Base{
		Status: 200,
		Msg:    "OK",
		Data: serialize.Login{
			Reply: "登录成功！",
			Token: token,
		},
	}
}

func (u UserUpdate) Update(ID uint, headPhoto *multipart.FileHeader, url string) *serialize.Base {
	model.DB.Model(&model.User{}).Where("id = ?", ID).
		Updates(map[string]interface{}{"name": u.Username, "age": u.Age,
			"gender": u.Gender, "birthday": u.Birthday, "email": u.Email,
			"signature": u.Signature, "head_photo": headPhoto.Filename})
	return &serialize.Base{
		Status: 200,
		Msg:    "ok",
		Data: serialize.Update{
			Reply: "更新成功！",
			Url:   url,
		},
	}
}

func MakeFriends(userID string, friendID string) *serialize.Base {
	fid, _ := strconv.Atoi(friendID)
	uid, _ := strconv.Atoi(userID)
	var user model.User
	var friend model.User
	model.DB.Model(&model.User{}).Where("id = ?", fid).First(&user)
	model.DB.Model(&model.User{}).Where("id = ?", uid).First(&friend)
	err := model.DB.Model(&user).Association("Friends").Append(&friend)
	if err != nil {
		log.Println(err)
	}
	err = model.DB.Model(&friend).Association("Friends").Append(&user)
	if err != nil {
		log.Println(err)
	}
	return &serialize.Base{
		Status: 200,
		Msg:    "ok",
		Data:   "交友成功",
	}
}

func (c *GroupRegister) CreateGroup(userID uint) *serialize.Base {
	var group model.Group
	var user model.User
	model.DB.Model(&model.User{}).Where("id = ?", userID).First(&user)
	group = model.Group{
		Name:  c.Name,
		Users: []model.User{user},
	}
	model.DB.Create(&group)
	return &serialize.Base{
		Status: 200,
		Msg:    "ok",
		Data:   "创建成功！",
	}
}

func PullGroup(userID string, GroupID int) *serialize.Base {
	uid, err := strconv.Atoi(userID)
	if err != nil {
		panic(err)
	}
	var user model.User
	var group model.Group
	model.DB.Model(&model.User{}).Where("id = ?", uid).First(&user)
	model.DB.Model(&model.Group{}).Where("id = ?", GroupID).First(&group)
	err = model.DB.Model(&group).Association("Users").Append(&user)
	if err != nil {
		log.Println(err)
	}
	return &serialize.Base{
		Status: 200,
		Msg:    "ok",
		Data:   "进群成功！",
	}
}

func FindUser(id string) model.User {
	var userInfo model.User
	model.DB.Model(&model.User{}).Where("id = ?", id).First(&userInfo)
	return userInfo
}

func FindFriends(id string) []model.User {
	var userInfo model.User
	var friends []model.User
	model.DB.Model(&model.User{}).Where("id = ?", id).First(&userInfo)
	err := model.DB.Model(&userInfo).Association("Friends").Find(&friends)
	if err != nil {
		log.Println(err)
	}
	return friends
}

func FindGroupUser(groupID int, id string) []model.User {
	var group model.Group
	var users []model.User
	model.DB.Model(&model.Group{}).Where("id = ?", groupID).First(&group)
	err := model.DB.Model(&group).Where("user_id = ?", id).
		Association("Users").Find(&users)
	if err != nil {
		log.Println(err)
	}
	return users
}

func FindGroup(userId string) []model.Group {
	var groups []model.Group
	var user model.User
	model.DB.Model(&model.User{}).Where("id = ?", userId).First(&user)
	err := model.DB.Model(&user).Association("Groups").Find(&groups)
	if err != nil {
		log.Println(err)
	}
	return groups
}

func FindMembers(groupId int) []model.User {
	var group model.Group
	var users []model.User
	model.DB.Model(&model.Group{}).Where("id = ?", groupId).First(&group)
	err := model.DB.Model(&group).Association("Users").Find(&users)
	if err != nil {
		log.Println(err)
	}
	return users
}

func FindOneGroup(groupID int) model.Group {
	var group model.Group
	model.DB.Model(&model.Group{}).Where("id = ?", groupID).First(&group)
	return group
}
