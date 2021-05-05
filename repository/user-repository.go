package repository

import (
	"log"

	"github.com/mop/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//UserRepository is the contract that shows what the userRepository can and has to do
type UserRepository interface {
	SaveUser(user entity.User) entity.User
	UpdateUser(user entity.User) entity.User
	VerifyCredentials(email string, password string) interface{}
	IsDuplicateEmail(email string) (tx *gorm.DB)
	FindByEmail(email string) entity.User
	ProfileUser(userID string) entity.User
	GetAnswers() []entity.MostAnswersDTO
}

type userConnection struct {
	connection *gorm.DB
}

///NewUserRepository creates a new instance of the user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userConnection{
		connection: db,
	}
}

func (db *userConnection) SaveUser(user entity.User) entity.User {
	user.Password = hashAndSalt([]byte(user.Password))
	db.connection.Save(&user)
	return user

}

func (db *userConnection) UpdateUser(user entity.User) entity.User {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser entity.User
		db.connection.Find(&tempUser, user.ID)
		user.Password = tempUser.Password
	}

	db.connection.Save(&user)
	return user
}

func (db *userConnection) VerifyCredentials(email string, password string) interface{} {
	var user entity.User
	res := db.connection.Where("email = ?", email).Take(&user)
	if res.Error == nil {
		return user
	}
	return nil
}

func (db *userConnection) IsDuplicateEmail(email string) (tx *gorm.DB) {
	var user entity.User
	return db.connection.Where("email = ?", email).Take(&user)
}

func (db *userConnection) FindByEmail(email string) entity.User {
	var user entity.User
	db.connection.Where("email = ?", email).Take(&user)
	return user
}

func (db *userConnection) ProfileUser(userID string) entity.User {
	var user entity.User
	db.connection.Preload("Questions").Preload("Questions.User").Find(&user, userID)
	return user
}

func (db *userConnection) GetAnswers() []entity.MostAnswersDTO {

	var result []entity.MostAnswersDTO
	db.connection.Raw("select u.id as UserID, u.firstname as Firstname, u.lastname as Lastname, COUNT(a.user_id) as NumberOfAnswers from users u join answers a on u.id = a.user_id group by a.user_id order by COUNT(a.user_id) desc").Scan(&result)
	return result
}

// select u.*, COUNT(a.user_id)  from users u join answers a on u.id = a.user_id group by a.user_id order by COUNT(a.user_id) desc;

//HashAndSalt method hashes the password of the users
func hashAndSalt(pass []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		panic("Failed to hash the password")
	}
	return string(hash)
}
