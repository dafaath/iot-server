package database

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/dafaath/iot-server/configs"
	"github.com/dafaath/iot-server/internal/entities"
	"github.com/dafaath/iot-server/internal/helper"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SQLType int

const (
	TABLE SQLType = iota
	DROP
	ADMIN
	HARDWARE
	NODE
	SENSOR
	CHANNEL
)

func hashPassword(ctx context.Context, password string) (hashedPassword string, err error) {
	hasher := sha256.New()
	hasher.Write([]byte(password))
	passwordHashBytes := hasher.Sum(nil)
	passwordHashString := hex.EncodeToString(passwordHashBytes)
	return passwordHashString, nil
}

func openSqlFile(sqlType SQLType) string {
	var path string
	sqlFolderPath := filepath.Join("internal", "database", "sql")
	switch sqlType {
	case TABLE:
		path = filepath.Join(sqlFolderPath, "table.sql")
	case DROP:
		path = filepath.Join(sqlFolderPath, "drop.sql")
	case HARDWARE:
		path = filepath.Join(sqlFolderPath, "hardware.sql")
	case NODE:
		path = filepath.Join(sqlFolderPath, "node.sql")
	case SENSOR:
		path = filepath.Join(sqlFolderPath, "sensor.sql")
	case CHANNEL:
		path = filepath.Join(sqlFolderPath, "channel.sql")
	default:
		panic("There is no sqltype for this code")
	}

	inp, err := ioutil.ReadFile(path)
	helper.PanicIfError(err)

	sqlStatement := string(inp)
	return sqlStatement
}

func createAdminData(tx pgx.Tx, config *configs.Config) error {
	log.Println("Creating admin data")
	hashedPassword, err := hashPassword(context.Background(), config.Account.AdminPassword)
	if err != nil {
		return err
	}

	admin := entities.User{
		Email:    config.Account.AdminEmail,
		Username: config.Account.AdminUsername,
		Password: hashedPassword,
		Status:   true,
		Token:    uuid.New().String(),
		IsAdmin:  true,
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO user_person (email, username, password, status, token, isAdmin) VALUES ($1, $2, $3, $4, $5, $6)`,
		admin.Email, admin.Username, admin.Password, admin.Status, admin.Token, admin.IsAdmin)
	log.Println("Finish creating admin data")
	return err
}

func createUserData(tx pgx.Tx, config *configs.Config) error {
	log.Println("Creating user data")
	hashedPassword, err := hashPassword(context.Background(), config.Account.UserPassword)
	if err != nil {
		return err
	}

	user := entities.User{
		Email:    config.Account.UserEmail,
		Username: config.Account.UserUsername,
		Password: hashedPassword,
		Status:   true,
		Token:    uuid.New().String(),
		IsAdmin:  false,
	}
	_, err = tx.Exec(context.Background(), `INSERT INTO user_person (email, username, password, status, token, isAdmin) VALUES ($1, $2, $3, $4, $5, $6)`,
		user.Email, user.Username, user.Password, user.Status, user.Token, user.IsAdmin)
	log.Println("Finish creating user data")
	return err
}

func createHardware(tx pgx.Tx) error {
	log.Println("Creating hardware")
	sqlStatement := openSqlFile(HARDWARE)
	_, err := tx.Exec(context.Background(), sqlStatement)
	return err
}

func createNode(tx pgx.Tx) error {
	log.Println("Creating node")
	sqlStatement := openSqlFile(NODE)
	_, err := tx.Exec(context.Background(), sqlStatement)
	return err
}

func createSensor(tx pgx.Tx) error {
	log.Println("Creating sensor")
	sqlStatement := openSqlFile(SENSOR)
	_, err := tx.Exec(context.Background(), sqlStatement)
	return err
}

func createChannel(tx pgx.Tx) error {
	log.Println("Creating channel")
	sqlStatement := openSqlFile(CHANNEL)
	_, err := tx.Exec(context.Background(), sqlStatement)
	return err
}

func createMockData(tx pgx.Tx, config *configs.Config) error {
	log.Println("Creating mock data")
	err := createAdminData(tx, config)
	if err != nil {
		return err
	}

	err = createUserData(tx, config)
	if err != nil {
		return err
	}

	err = createHardware(tx)
	if err != nil {
		return err
	}
	err = createNode(tx)
	if err != nil {
		return err
	}
	err = createSensor(tx)
	if err != nil {
		return err
	}
	err = createChannel(tx)
	if err != nil {
		return err
	}

	log.Println("Finish creating mock data")
	return nil
}

func createTable(tx pgx.Tx) error {
	log.Println("Creating table")
	sqlStatement := openSqlFile(TABLE)
	_, err := tx.Exec(context.Background(), sqlStatement)
	return err
}

func DropTable() {
	db, err := GetConnection()
	helper.PanicIfError(err)

	log.Println("Drop table")
	sqlStatement := openSqlFile(DROP)
	_, err = db.Exec(context.Background(), sqlStatement)
	log.Println("Finish drop table")

	helper.PanicIfError(err)
}

func CreateTableAndMockData() {
	db, err := GetConnection()
	helper.PanicIfError(err)
	config := configs.GetConfig()

	tx, err := db.Begin(context.Background())
	helper.PanicIfError(err)

	err = createTable(tx)
	helper.PanicIfError(err)

	err = createMockData(tx, config)
	helper.PanicIfError(err)

	defer func() {
		recoverErr := recover()
		if recoverErr != nil {
			log.Println("[ERROR]", recoverErr)
			errorRollback := tx.Rollback(context.Background())
			if errorRollback != nil {
				log.Println("Error in rollback", errorRollback)
			} else {
				log.Println("Rollback success")
			}
		} else {
			errorCommit := tx.Commit(context.Background())
			if errorCommit != nil {
				log.Println("Error in commit", errorCommit)
			} else {
				log.Println("Commit success")
			}
		}
	}()
}
