package service

import (
	"database/sql"
	"errors"

	"github.com/lib/pq"

	pb "users/proto"
)

func GetUserInfoFromDBByLogin(db *sql.DB, userLogin string) (*pb.GetUserInfoResponse, error) {
	result := &pb.GetUserInfoResponse{}
	if err := db.QueryRow("SELECT name, user_id, avatar_path FROM users WHERE login = $1", userLogin).Scan(&result.Name, &result.UserID, &result.AvatarPath); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("Не удалось найти данного пользователя")
		default:
			return nil, err
		}
	}

	result.Login = userLogin

	return result, nil
}

func GetUserInfoFromDBByID(db *sql.DB, userID int) (*pb.GetUserInfoResponse, error) {
	result := &pb.GetUserInfoResponse{}
	if err := db.QueryRow("SELECT login, name, avatar_path FROM users WHERE user_id = $1", userID).Scan(&result.Login, &result.Name, &result.AvatarPath); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, errors.New("Не удалось найти данного пользователя")
		default:
			return nil, err
		}
	}

	result.UserID = int32(userID)

	return result, nil
}

func GetUserInfoForAuthFromDB(db *sql.DB, userLogin string) (*pb.AuthResponse, string, error) {
	result := &pb.AuthResponse{}
	var currentPassword string
	if err := db.QueryRow("SELECT user_id, password, name, avatar_path FROM users WHERE login = $1", userLogin).Scan(&result.UserID, &currentPassword, &result.UserName, &result.AvatarPath); err != nil {
		switch err {
		case sql.ErrNoRows:
			return nil, "", errors.New("Неверный логин или пароль")
		default:
			return nil, "", err
		}
	}

	return result, currentPassword, nil
}

func RegisterUserInDB(db *sql.DB, login, name, hashedPassword string) (int, error) {
	var userID int
	if err := db.QueryRow("INSERT INTO users (login, name, password) VALUES ($1, $2, $3) RETURNING user_id", login, name, hashedPassword).Scan(&userID); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505":
				return 0, errors.New("Пользователь с таким логин уже существует")
			default:
				return 0, errors.New("Неизвестная ошибка, попробуйте позже")
			}
		} else {
			return 0, errors.New("Неизвестная ошибка, попробуйте позже")
		}
	}

	return userID, nil
}

func AddUserAvatarToDB(db *sql.DB, userID int, avatarPath string) error {
	result, err := db.Exec("UPDATE users SET avatar_path = $1 WHERE user_id = $2", avatarPath, userID)
	if rows, _ := result.RowsAffected(); rows == 0 {
		return errors.New("Нет пользователя с данным id")
	}
	if err != nil {
		return err
	}

	return nil
}

func DeleteUserInDB(db *sql.DB, userID int) error {
	if _, err := db.Exec("DELETE FROM users WHERE user_id = $1", userID); err != nil {
		return err
	}

	return nil
}

func getUserIDByLoginFromDB(db *sql.DB, login string) (int, error) {
	var userID int
	if err := db.QueryRow("SELECT user_id FROM users WHERE login = $1", login).Scan(&userID); err != nil {
		return userID, err
	}

	return userID, nil
}
