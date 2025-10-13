package repository

import (
	"database/sql"
	"errors"
	"log"

	"bytes"

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
	// todo! разобраться

	// 	var exists bool
	// const query = "SELECT EXISTS(SELECT 1 FROM users WHERE login=$1)"

	// err := DB.Get(&exists, query, login)

	var placeholder int
	query := "SELECT 1 FROM users WHERE login = ? LIMIT 1"

	err := DB.Get(&placeholder, query, login)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func CheckPassExists(login string, passwordHash []byte) (bool, error) {
	var user = User{}

	// Выбираем только те поля, которые нам нужны для проверки.
	// В данном случае - только хэш пароля.
	const query = "SELECT password_hash FROM users WHERE login = ? LIMIT 1"

	err := DB.Get(&user, query, login)

	log.Printf("Checking user '%s', err: %v", login, err)

	if err != nil {

		return false, err
	}

	// 4. СРАВНИВАЕМ ХЭШИ ПАРОЛЕЙ
	// Пользователь найден, теперь 'user.Password' содержит хэш из базы.
	// Слайсы байт ([]byte) нужно сравнивать с помощью bytes.Equal.
	// Простое сравнение `user.Password == passwordHash` не будет работать правильно.
	if bytes.Equal(user.Password, passwordHash) {
		return true, nil
	}

	// Хэши не совпали
	return false, nil
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
