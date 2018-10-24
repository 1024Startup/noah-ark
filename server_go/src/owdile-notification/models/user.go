package models

import (
	"github.com/astaxie/beego/orm"
)

var (
	UserList map[string]*User
)

func init() {
	//UserList = make(map[string]*User)
	//u := User{"user_11111", "astaxie", "11111", Profile{"male", 20, "Singapore", "astaxie@gmail.com"}}
	//UserList["user_11111"] = &u
}

type User struct {
	Id       int    `orm:"Column(Id);pk"`
	Username string `orm:"size(32);default();column(username)"`
	Password string `orm:"size(32);default();column(password)"`
	Test     string `orm:"size(32);default();column(test)"`
}

func GetAllUsers() interface{} {
	var info User

	err := orm.NewOrm().QueryTable(GetTable("user")).Filter("Id", 1).One(&info)
	if err != nil {
		panic(err)
	}
	//fmt.Printf(info)

	return info.Id
}

/*func AddUser(u User) string {
	u.Id = "user_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	UserList[u.Id] = &u
	return u.Id
}

func GetUser(uid string) (u *User, err error) {
	if u, ok := UserList[uid]; ok {
		return u, nil
	}
	return nil, errors.New("User not exists")
}

func GetAllUsers() map[string]*User {
	return UserList
}

func UpdateUser(uid string, uu *User) (a *User, err error) {
	if u, ok := UserList[uid]; ok {
		if uu.Username != "" {
			u.Username = uu.Username
		}
		if uu.Password != "" {
			u.Password = uu.Password
		}
		if uu.Profile.Age != 0 {
			u.Profile.Age = uu.Profile.Age
		}
		if uu.Profile.Address != "" {
			u.Profile.Address = uu.Profile.Address
		}
		if uu.Profile.Gender != "" {
			u.Profile.Gender = uu.Profile.Gender
		}
		if uu.Profile.Email != "" {
			u.Profile.Email = uu.Profile.Email
		}
		return u, nil
	}
	return nil, errors.New("User Not Exist")
}

func Login(username, password string) bool {
	for _, u := range UserList {
		if u.Username == username && u.Password == password {
			return true
		}
	}
	return false
}

func DeleteUser(uid string) {
	delete(UserList, uid)
}
*/
