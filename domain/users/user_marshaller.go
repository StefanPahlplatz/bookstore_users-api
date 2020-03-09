package users

import "encoding/json"

type PublicUser struct {
	Id        int64  `json:"id"`
	CreatedOn string `json:"created_on"`
	Status    string `json:"status"`
}

type PrivateUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	CreatedOn string `json:"created_on"`
	Status    string `json:"status"`
}

func (users Users) Marshall(isPublic bool) interface{} {
	result := make([]interface{}, len(users))
	for index, user := range users {
		result[index] = user.Marshall(isPublic)
	}
	return result
}

func (user *User) Marshall(isPublic bool) interface{} {
	if isPublic {
		return PublicUser{
			Id:        user.Id,
			CreatedOn: user.CreatedOn,
			Status:    user.Status,
		}
	}
	userJson, _ := json.Marshal(user)
	var privateUser PrivateUser
	json.Unmarshal(userJson, &privateUser)
	return privateUser
}
