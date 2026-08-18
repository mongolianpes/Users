package service

import (
	"database/sql"
	"fmt"
	"testing"
	"users/crypto"
)

type updMessagesInfo struct {
	usersIDs []int
}

func updateTestUsers(db *sql.DB) ([]int, error) {
	if _, err := db.Exec("DELETE FROM users"); err != nil {
		return []int{}, err
	}

	hashedPass, err := crypto.HashString("1234567890")
	if err != nil {
		return []int{}, err
	}

	rowsUserIDs, err := db.Query(`INSERT INTO users (login, name, password) VALUES
		(1, '1', $1),
		(2, '2', $1),
		(3, '3', $1)
		RETURNING user_id`, hashedPass)
	if err != nil {
		return []int{}, err
	}

	var updUsers []int
	var userID1 int
	var userID2 int
	var userID3 int
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID1)
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID2)
	rowsUserIDs.Next()
	rowsUserIDs.Scan(&userID3)
	updUsers = append(updUsers, userID1)
	updUsers = append(updUsers, userID2)
	updUsers = append(updUsers, userID3)

	var userIDAvatar int
	err = db.QueryRow(`INSERT INTO users (login, name, password, avatar_path) VALUES
	('withavatar', 'avatarka', $1, 'path-to-avatar') RETURNING user_id`, hashedPass).Scan(&userIDAvatar)
	if err != nil {
		return []int{}, err
	}
	updUsers = append(updUsers, userIDAvatar)

	return updUsers, err
}

func connectTestToDB() (*sql.DB, error) {
	host := "localhost"
	port := "5432"
	user := "postgres"
	password := "123"
	dbname := "project_farm"

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	err = db.Ping()
	if err != nil {
		return db, err
	}

	return db, nil
}

func TestGetUserInfoFromDBByLogin(t *testing.T) {
	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	updUsers, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	resp, err := GetUserInfoFromDBByLogin(db, "1")
	if err != nil {
		t.Error(err)
	}

	if resp.Name != "1" {
		t.Error("Вернул не ожидаемое имя")
	}

	if resp.UserID != int32(updUsers[0]) {
		t.Error("Вернул не ожидаемый userID")
	}

	if resp.AvatarPath != "" {
		t.Error("Вернул не ожидаемый путь к аватару")
	}

	if resp.Login != "1" {
		t.Error("Вернул не ожидаемый login")
	}

	respWithAvatar, err := GetUserInfoFromDBByLogin(db, "withavatar")
	if err != nil {
		t.Error(err)
	}

	if respWithAvatar.Name != "avatarka" {
		t.Error(err)
	}

	if respWithAvatar.UserID != int32(updUsers[3]) {
		t.Error(err)
	}

	if respWithAvatar.AvatarPath != "path-to-avatar" {
		t.Error(err)
	}
}

func TestGetUserInfoFromDBByID(t *testing.T) {
	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	updUsers, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	resp, err := GetUserInfoFromDBByID(db, updUsers[0])
	if err != nil {
		t.Error(err.Error()+" userID ", updUsers[0])
	}

	if resp.Name != "1" {
		t.Error("Вернул не ожидаемое имя")
	}

	if resp.UserID != int32(updUsers[0]) {
		t.Error("Вернул не ожидаемый userID")
	}

	if resp.AvatarPath != "" {
		t.Error("Вернул не ожидаемый путь к аватару")
	}

	if resp.Login != "1" {
		t.Error("Вернул не ожидаемый login")
	}

	respWithAvatar, err := GetUserInfoFromDBByID(db, updUsers[3])
	if err != nil {
		t.Error(err)
	}

	if respWithAvatar.Name != "avatarka" {
		t.Error(err)
	}

	if respWithAvatar.UserID != int32(updUsers[3]) {
		t.Error(err)
	}

	if respWithAvatar.AvatarPath != "path-to-avatar" {
		t.Error(err)
	}
}

func TestGetUserInfoForAuthFromDB(t *testing.T) {
	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	updUsers, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	resp, currentPass, err := GetUserInfoForAuthFromDB(db, "2")
	if err != nil {
		t.Error(err)
	}

	if resp.UserID != int32(updUsers[1]) {
		t.Error("Вернул не ожидаемый userID")
	}

	if resp.UserName != "2" {
		t.Error("Вернул не ожидаемый username")
	}

	if !crypto.VerifyHash("1234567890", currentPass) {
		t.Error("Вернул не вреный пароль")
	}

	if resp.AvatarPath != "" {
		t.Error("Вернул не ожидаемый avatar path")
	}
}

func TestRegister(t *testing.T) {
	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	_, err = updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	hashedPass, err := crypto.HashString("1234567890")
	if err != nil {
		t.Error(err)
	}

	userID, err := RegisterUserInDB(db, "regtest", "regtest", hashedPass)
	if err != nil {
		t.Error(err)
	}
	if userID == 0 {
		t.Error("Вернул userID 0")
	}

	userID, err = RegisterUserInDB(db, "1", "1", hashedPass)
	if err == nil {
		t.Error(err)
	}
}

func TestAddUserAvatarToDB(t *testing.T) {
	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	updUsers, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	if err := AddUserAvatarToDB(db, updUsers[0], "avatar path"); err != nil {
		t.Error(err)
	}

	var avatarPathInDB string
	if err := db.QueryRow("SELECT avatar_path FROM users WHERE user_id = $1", updUsers[0]).Scan(&avatarPathInDB); err != nil {
		t.Error(err)
	}

	if avatarPathInDB != "avatar path" {
		t.Error("Вернул не верный avatar path")
	}
}

func TestDeleteUserInDB(t *testing.T) {
	db, err := connectTestToDB()
	if err != nil {
		t.Error(err)
	}

	updUsers, err := updateTestUsers(db)
	if err != nil {
		t.Error(err)
	}

	if err := DeleteUserInDB(db, updUsers[0]); err != nil {
		t.Error(err)
	}

	var login string
	if err := db.QueryRow("SELECT login FROM users WHERE user_id = $1", updUsers[0]).Scan(&login); err == nil {
		t.Error("Пользователь не был удален из БД")
	}
}
