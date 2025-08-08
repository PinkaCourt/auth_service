package repository

// Шаг 2: Репозиторий для работы с БД
// Создай internal/repository/user_repo.go
// Перенеси туда:
// Логику подключения к БД из connectBD() → Конструктор NewUserRepository()
// SQL-запросы (CREATE TABLE, INSERT)
// Метод CreateUser(login, hashedPassword)
// Проверку на дубликаты логина
