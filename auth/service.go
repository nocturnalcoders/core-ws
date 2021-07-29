package auth

//service generate token
import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

//tentukan dahulu servicenya mau ngapain
//1. bagaimana cara membuat token (Generate)
//2. Melakukan validasi token

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
}

//Secret key ->syarat -> sebuah token dianggap valid jika dia dibuat dengan secret key yang sama dengan yang terdaftar di server
var SECRET_KEY = []byte("EKOST_s3cr3T_k3Y")

func NewService() *jwtService {
	return &jwtService{}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	//Variabel dama JWT.IO Payload disebut juga dengan claim
	//Melakukan passing dengan User ID sebagai key, dan valuenya adalah user id yang dipassing melalui func generate token
	claim := jwt.MapClaims{}
	claim["user_id"] = userID

	//Membuat Token
	//menggunakan HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}

	return signedToken, nil

}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	//perlu di parse dlu tokennya
	//hal yang dilakukan
	//1. memasukan token
	//2. parameter adalah sebuah func dimana dia mengembalikan interace dan error
	//3. jika method sudah benar dengan algo HS256 dia akan dicek dan mecocokan method
	//4. kemudian jika benar dia akan mengembalikan secret key
	//5. supaya encoded token apakah benar di buat dengan secret key
	//6. jika beda dia tidak akan valid
	//7. jika token berhasil divalidasi maka akan masuk

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return token, err
	}
	return token, nil
}
