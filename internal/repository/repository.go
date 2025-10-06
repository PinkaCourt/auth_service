package repository

import (
	"database/sql"
	"errors"
	"log"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

type User struct {
	ID       int    `db:"id"`
	Name     string `db:"login"`
	Password []byte `db:"password_hash"`
}

var DB *sqlx.DB

func ConnectDB() {
	var err error
	DB, err = sqlx.Open("sqlite", "test.db")
	if err != nil {
		log.Fatalf("Не удалось открыть подключение: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе: %v", err)
	}
	log.Println("Подключено к SQLite!")
}

func CreateUsersTable() {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		login TEXT NOT NULL UNIQUE,
		password_hash BLOB NOT NULL
	);`
	// schema := `DROP TABLE IF EXISTS users`
	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("Не удалось создать таблицу 'users': %v", err)
		return
	}
	log.Println("Таблица 'users' успешно проверена/создана.")
}

// для удаления кривой таблицы
func DeleteUsersTable() {
	schema := `DROP TABLE IF EXISTS users`
	_, err := DB.Exec(schema)
	if err != nil {
		log.Fatalf("Не удалось удалить таблицу 'users': %v", err)
		return
	}
	log.Println("Таблица 'users' успешно удалена")
}

func CheckUserExists(login string) (bool, error) {
	// todo разобраться
	var placeholder int
	query := "SELECT 1 FROM users WHERE login = ? LIMIT 1"

	err := DB.Get(&placeholder, query, login)
	if err != nil {
		if err == sql.ErrNoRows {
			// Это не ошибка, а нормальный результат - пользователя нет
			return false, nil
		}

		return false, err
	}

	return true, nil
}

// RegisterUser регистрирует нового пользователя в БД
func RegisterUser(login string, passwordHash []byte) error {
	exists, err := CheckUserExists(login)
	if err != nil {
		return err
	}
	if exists {
		println("Пользователь с таким логином уже существует")
		// log.New("Пользователь с таким логином уже существует", nil)
		return errors.New("пользователь с таким логином уже существует")
	}

	// Создаем пользователя и вставляем в БД
	newUser := User{
		Name:     login,
		Password: passwordHash,
	}
	// todo это мой код
	// _, err := DB.NamedExec("INSERT INTO users (login, password) VALUES (:login, :password)", newUser)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// Обратите внимание на поля в INSERT и в структуре. Они должны совпадать.
	// В таблице: login, password_hash. В структуре теги `db:"login"` и `db:"password_hash"`
	query := "INSERT INTO users (login, password_hash) VALUES (:login, :password_hash)"
	_, err = DB.NamedExec(query, newUser)
	if err != nil {
		log.Printf("Не удалось зарегистрировать пользователя: %v", err)
		return err
	}

	log.Println("Пользователь", login, "успешно зарегистрирован")
	return nil
}
