package service

import (
	"IMChat/model"
	"IMChat/serialize"
	"IMChat/utils"
	"golang.org/x/crypto/bcrypt"
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
		Data:   serialize.Datalist{Item: token, Total: 1},
	}
}

func (u UserUpdate) Update(ID uint) *serialize.Base {
	model.DB.Model(&model.User{}).Where("id = ?", ID).
		Updates(map[string]interface{}{"name": u.Username, "age": u.Age,
			"gender": u.Gender, "birthday": u.Birthday, "email": u.Email,
			"signature": u.Signature})
	return &serialize.Base{
		Status: 200,
		Msg:    "ok",
		Data:   "修改成功！",
	}
}

//func MakeFriends(userID uint, blackID string) *serialize.Base {
//id, err := strconv.Atoi(blackID)
//var User model.User
//if err != nil {
//	panic("数据格式错误！")
//}
//model.DB.Model(&model.User{}).Where("id = ?", userID).Find(&User)
//
//}
