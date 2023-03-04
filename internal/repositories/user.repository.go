package repositories

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/dafaath/iot-server/v2/configs"
	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"

	"gopkg.in/gomail.v2"
)

type UserRepository struct {
	mailDialer *gomail.Dialer
}

func (u *UserRepository) Create(ctx context.Context, tx helper.Querier, payload entities.UserCreate) (user entities.UserRead, err error) {
	hashedPassword, err := u.hashPassword(context.Background(), payload.Password)
	if err != nil {
		return user, err
	}

	status := false
	isAdmin := false

	user = entities.UserRead{
		Email:    payload.Email,
		Username: payload.Username,
		Status:   status,
		IsAdmin:  isAdmin,
	}
	sqlStatement := `
	INSERT INTO user_person (
		email,
		username,
		password,
		status,
		isadmin  
	)
	VALUES ($1, $2, $3, $4, $5) RETURNING id_user`
	err = tx.QueryRow(ctx, sqlStatement, payload.Email, payload.Username, hashedPassword, status, isAdmin).Scan(&user.IdUser)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (u *UserRepository) GetAll(ctx context.Context, tx helper.Querier) (users []entities.UserRead, err error) {
	users = []entities.UserRead{}
	sqlStatement := `SELECT id_user, email, username,  status,  isadmin FROM user_person`
	rows, err := tx.Query(ctx, sqlStatement)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.UserRead
		err := rows.Scan(
			&user.IdUser,
			&user.Email,
			&user.Username,
			&user.Status,
			&user.IsAdmin,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (u *UserRepository) GetById(ctx context.Context, tx helper.Querier, id int) (user entities.UserRead, err error) {
	sqlStatement := `SELECT id_user, email, username, status,  isadmin FROM user_person WHERE id_user=$1`
	err = tx.QueryRow(ctx, sqlStatement, id).Scan(
		&user.IdUser,
		&user.Email,
		&user.Username,
		&user.Status,
		&user.IsAdmin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, fiber.NewError(404, fmt.Sprintf("User with id %d not found", id))
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepository) UpdatePassword(ctx context.Context, tx helper.Querier, id int, password string) (err error) {
	hashPassword, err := u.hashPassword(ctx, password)
	if err != nil {
		return err
	}

	sqlStatement := `
	UPDATE user_person
	SET password=$1 
	WHERE id_user=$2`
	res, err := tx.Exec(ctx, sqlStatement, hashPassword, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on update user password with id %d", id))
	}
	return nil
}

func (u *UserRepository) UpdateStatus(ctx context.Context, tx helper.Querier, id int, status bool) (err error) {
	sqlStatement := `
	UPDATE user_person 
	set status=$1 
	WHERE id_user=$2`
	res, err := tx.Exec(ctx, sqlStatement, status, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on update user status with id %d", id))
	}
	return nil
}

func (u *UserRepository) Delete(ctx context.Context, tx helper.Querier, id int) (err error) {
	sqlStatement := `DELETE FROM user_person WHERE id_user=$1`
	res, err := tx.Exec(ctx, sqlStatement, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on delete with id %d", id))
	}
	return nil
}

func (u *UserRepository) GetByEmail(ctx context.Context, tx helper.Querier, email string) (user entities.UserRead, err error) {
	sqlStatement := `SELECT id_user, email, username,  status,   isAdmin FROM user_person WHERE email=$1`
	err = tx.QueryRow(ctx, sqlStatement, email).Scan(
		&user.IdUser,
		&user.Email,
		&user.Username,
		&user.Status,
		&user.IsAdmin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, fiber.NewError(404, fmt.Sprintf("User with email %s not found", email))
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepository) GetByUsername(ctx context.Context, tx helper.Querier, username string) (user entities.UserRead, err error) {
	sqlStatement := `SELECT id_user, email, username,  status, isAdmin FROM user_person WHERE username=$1`
	err = tx.QueryRow(ctx, sqlStatement, username).Scan(
		&user.IdUser,
		&user.Email,
		&user.Username,
		&user.Status,
		&user.IsAdmin,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, fiber.NewError(404, fmt.Sprintf("User with username %s not found", username))
		}
		return user, err
	}
	return user, nil
}

func (u *UserRepository) MatchPassword(ctx context.Context, tx helper.Querier, user entities.UserRead, password string) (err error) {
	var userPassword string
	sqlStatement := `SELECT password FROM user_person WHERE id_user=$1`
	err = tx.QueryRow(ctx, sqlStatement, user.IdUser).Scan(
		&userPassword,
	)
	if err != nil {
		return err
	}

	passwordHashString, err := u.hashPassword(ctx, password)
	if err != nil {
		return err
	}

	if userPassword != passwordHashString {
		return fiber.NewError(401, "Wrong password")
	}

	return err
}

func (u *UserRepository) hashPassword(ctx context.Context, password string) (hashedPassword string, err error) {
	hasher := sha256.New()
	_, err = hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}

	passwordHashBytes := hasher.Sum(nil)
	passwordHashString := hex.EncodeToString(passwordHashBytes)
	return passwordHashString, nil
}

func (u *UserRepository) SendEmail(ctx context.Context, to string, subject string, body string) (err error) {
	configs := configs.GetConfig()

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", configs.Mail.SenderName)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	err = u.mailDialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) SendEmailActivation(ctx context.Context, user entities.UserRead) (err error) {
	configs := configs.GetConfig()

	jwtToken, err := helper.SignUserToken(user)
	if err != nil {
		return err
	}

	urlCode := fmt.Sprintf("http://%s:%d/user/activation?token=%s", configs.Server.Host, configs.Server.Port, jwtToken)
	subject := "Registration Email"
	body := fmt.Sprintf(`<html>              
			<head>>
				<title>Activation Message</title>
			</head>
			<body>
			
				<h1>Activation Message</h1>
				<h4>Dear %s</h4>
				<p>We have accepted your registration. Your account is:</p>
				<li>
					<ul> Id User: %d </ul>
					<ul> Username: %s </ul>
				</li>
				<p>Click <a href=%s>here</a> to activate your account</p>
				<p><h5>Thank you</h5></p>     
			</body>
    </html> `, user.Username, user.IdUser, user.Username, urlCode)

	err = u.SendEmail(ctx, user.Email, subject, body)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) SendEmailForgotPassword(ctx context.Context, user entities.UserRead, newPassword string) (err error) {
	subject := "Forgot Password Email"
	body := fmt.Sprintf(`<html>
		  <head>
		  </head>
		  <body>
			<h3>Dear %s. </h3>
			<p>We have accepted your forget password request. Use this password for log in.</p>
			<p><h4>%s</h4></p>
			<p>Thank You</p>
		  </body
		</html>`, user.Username, newPassword)

	err = u.SendEmail(ctx, user.Email, subject, body)
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRepository) SignJWT(ctx context.Context, user entities.UserRead) (token string, err error) {
	return helper.SignUserToken(user)
}

func NewUserRepository(mailDialer *gomail.Dialer) (UserRepository, error) {
	return UserRepository{
		mailDialer: mailDialer,
	}, nil
}
