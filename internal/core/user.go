package core

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/VladislavSCV/internal/models"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func GetAllUsers(db *sql.DB, cache *redis.Client) ([]models.User, error) {
	ctx := context.Background()
	cacheKey := "users:all"

	log.Println("Начало выполнения функции GetAllUsers")

	// Проверяем кеш
	log.Println("Проверка кеша для ключа:", cacheKey)
	cachedUsers, err := cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		log.Println("Данные найдены в кеше")
		var users []models.User
		if err := json.Unmarshal(cachedUsers, &users); err == nil {
			log.Println("Успешно распакованы данные из кеша")
			return users, nil
		}
		log.Println("Ошибка при распаковке данных из кеша:", err)
	} else {
		log.Println("Данные не найдены в кеше или ошибка при получении:", err)
	}

	// Если данных нет в кеше, запрашиваем из базы данных
	log.Println("Запрос данных из базы данных")
	var users []models.User

	if db == nil {
		log.Println("Ошибка: соединение с базой данных не установлено (db == nil)")
		return nil, fmt.Errorf("database connection is nil")
	}

	// Проверка соединения с базой данных
	log.Println("Проверка соединения с базой данных")
	err = db.Ping()
	if err != nil {
		log.Println("Ошибка при проверке соединения с базой данных:", err)
		return nil, fmt.Errorf("database connection is closed: %v", err)
	}

	// Выполнение SQL-запроса
	log.Println("Выполнение SQL-запроса для получения списка пользователей")
	rows, err := db.Query(`
        SELECT u.id, u.first_name, u.middle_name, u.last_name, r.value AS role, g.name AS group_name, u.login, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        LEFT JOIN groups g ON u.group_id = g.id
    `)
	if err != nil {
		log.Println("Ошибка при выполнении SQL-запроса:", err)
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	log.Println("Обработка результатов SQL-запроса")
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Role, &user.Group, &user.Login, &user.CreatedAt, &user.UpdatedAt); err != nil {
			log.Println("Ошибка при сканировании строки результата:", err)
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		log.Println("Ошибка при итерации по результатам запроса:", err)
		return nil, fmt.Errorf("error iterating over users: %v", err)
	}

	// Сохраняем данные в кеше
	log.Println("Сохранение данных в кеше")
	usersJSON, err := json.Marshal(users)
	if err == nil {
		log.Println("Успешно сериализованы данные для кеша")
		if err := cache.Set(ctx, cacheKey, usersJSON, 24*time.Hour).Err(); err != nil {
			log.Println("Ошибка при сохранении данных в кеше:", err)
		} else {
			log.Println("Данные успешно сохранены в кеше")
		}
	} else {
		log.Println("Ошибка при сериализации данных для кеша:", err)
	}

	log.Println("Успешное завершение функции GetAllUsers")
	return users, nil
}

func GetStudents(db *sql.DB, cache *redis.Client) ([]models.User, error) {
	ctx := context.Background()
	cacheKey := "users:students"

	// Проверяем кеш
	cachedStudents, err := cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var students []models.User
		if err := json.Unmarshal(cachedStudents, &students); err == nil {
			return students, nil
		}
	}

	// Если данных нет в кеше, запрашиваем из базы данных
	var students []models.User

	rows, err := db.Query(`
        SELECT u.id, u.first_name, u.middle_name, u.last_name, r.value AS role, g.name AS group_name, u.login, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        LEFT JOIN groups g ON u.group_id = g.id
        WHERE r.value = 'student'
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch students: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var student models.User
		if err := rows.Scan(&student.ID, &student.FirstName, &student.MiddleName, &student.LastName, &student.Role, &student.Group, &student.Login, &student.CreatedAt, &student.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan student: %v", err)
		}
		students = append(students, student)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over students: %v", err)
	}

	// Сохраняем данные в кеше
	studentsJSON, err := json.Marshal(students)
	if err == nil {
		cache.Set(ctx, cacheKey, studentsJSON, 24*time.Hour) // Кешируем на 24 часа
	}

	return students, nil
}

func GetTeachers(db *sql.DB, cache *redis.Client) ([]models.User, error) {
	ctx := context.Background()
	cacheKey := "users:teachers"

	// Проверяем кеш
	cachedTeachers, err := cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var teachers []models.User
		if err := json.Unmarshal(cachedTeachers, &teachers); err == nil {
			return teachers, nil
		}
	}

	// Если данных нет в кеше, запрашиваем из базы данных
	var teachers []models.User

	rows, err := db.Query(`
        SELECT u.id, u.first_name, u.middle_name, u.last_name, r.value AS role, g.name AS group_name, u.login, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        LEFT JOIN groups g ON u.group_id = g.id
        WHERE r.value = 'teacher'
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch teachers: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var teacher models.User
		if err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.MiddleName, &teacher.LastName, &teacher.Role, &teacher.Group, &teacher.Login, &teacher.CreatedAt, &teacher.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan teacher: %v", err)
		}
		teachers = append(teachers, teacher)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over teachers: %v", err)
	}

	// Сохраняем данные в кеше
	teachersJSON, err := json.Marshal(teachers)
	if err == nil {
		cache.Set(ctx, cacheKey, teachersJSON, 24*time.Hour) // Кешируем на 24 часа
	}

	return teachers, nil
}

func GetUserByID(db *sql.DB, cache *redis.Client, userID int) (*models.User, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Проверяем кеш
	cachedUser, err := cache.Get(ctx, cacheKey).Bytes()
	if err == nil {
		var user models.User
		if err := json.Unmarshal(cachedUser, &user); err == nil {
			return &user, nil
		}
	}

	// Если данных нет в кеше, запрашиваем из базы данных
	var user models.User

	err = db.QueryRow(`
        SELECT u.id, u.first_name, u.middle_name, u.last_name, r.value AS role, g.name AS group_name, u.login, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        LEFT JOIN groups g ON u.group_id = g.id
        WHERE u.id = $1
    `, userID).Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Role, &user.Group, &user.Login, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}

	// Сохраняем данные в кеше
	userJSON, err := json.Marshal(user)
	if err == nil {
		cache.Set(ctx, cacheKey, userJSON, 24*time.Hour) // Кешируем на 24 часа
	}

	return &user, nil
}

func UpdateUser(db *sql.DB, cache *redis.Client, updates map[string]interface{}) error {
	ctx := context.Background()
	userID := updates["id"].(int)
	cacheKey := fmt.Sprintf("user:%d", userID)

	// Формирование запроса
	query := "UPDATE users SET "
	var args []interface{}
	i := 1

	for k, v := range updates {
		if k == "id" {
			continue
		}
		query += fmt.Sprintf("%s = $%d, ", k, i) // Используем $ для подстановки
		args = append(args, v)
		i++
	}

	// Удаляем последнюю запятую
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, userID)

	// Выполнение запроса
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Инвалидируем кеш
	cache.Del(ctx, cacheKey)
	cache.Del(ctx, "users:all") // Инвалидируем кеш списка всех пользователей

	return nil
}

func DeleteUser(db *sql.DB, cache *redis.Client, userID int) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("user:%d", userID)

	_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Инвалидируем кеш
	cache.Del(ctx, cacheKey)
	cache.Del(ctx, "users:all") // Инвалидируем кеш списка всех пользователей

	return nil
}
