package core

import (
	"database/sql"
	"fmt"
	"github.com/VladislavSCV/internal/models"
)

func GetAllUsers(db *sql.DB) ([]models.User, error) {
	var users []models.User

	if db == nil {
		return nil, fmt.Errorf("database connection is nil")
	}

	rows, err := db.Query(`
        SELECT u.id, u.first_name, u.middle_name, u.last_name, r.value AS role, g.name AS group_name, u.login, u.created_at, u.updated_at
        FROM users u
        JOIN roles r ON u.role_id = r.id
        LEFT JOIN groups g ON u.group_id = g.id
    `)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.MiddleName, &user.LastName, &user.Role, &user.Group, &user.Login, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %v", err)
	}

	return users, nil
}

func GetStudents(db *sql.DB) ([]models.User, error) {
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

	return students, nil
}

func GetTeachers(db *sql.DB) ([]models.User, error) {
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

	return teachers, nil
}

func GetUserByID(db *sql.DB, userID int) (*models.User, error) {
	var user models.User

	err := db.QueryRow(`
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

	return &user, nil
}

func UpdateUser(db *sql.DB, updates map[string]interface{}) error {
	//log.Printf("Updating user with ID: %d, FirstName: %s, Role: %s, Group: %s", user.ID, user.FirstName, user.RoleID, user.Group)

	// Формирование запроса
	query := "UPDATE users SET "
	var args []interface{}
	i := 1

	for k, v := range updates {
		query += fmt.Sprintf("%s = $%d, ", k, i) // Используем $ для подстановки
		args = append(args, v)
		i++
	}

	// Удаляем последнюю запятую
	query = query[:len(query)-2]
	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, updates["id"])

	// Выполнение запроса
	_, err := db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	return nil
}

func DeleteUser(db *sql.DB, userID int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	return nil
}
