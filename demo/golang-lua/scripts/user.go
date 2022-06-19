package scripts

import "fmt"

type User struct {
	Name string
	Age  int
}

func (u *User) String() string {
	if u == nil {
		return ""
	}
	return fmt.Sprintf("name:%v, age:%v", u.Name, u.Age)
}

func (u *User) SetName(name string) {
	if u == nil {
		return
	}
	u.Name = name
}
