package repository

import (
	"time"

	"github.com/iMohamedSheta/xapp/pkg/logger"
	"github.com/iMohamedSheta/xqb"
)

type TodoRepository struct{}

func (r *TodoRepository) PaginatedUserTodos(userId int64, perPage int, page int, withMeta bool) ([]map[string]any, map[string]any, error) {

	userTodos, meta, err := xqb.Table("todos").Where("user_id", "=", userId).Paginate(perPage, page, withMeta)

	if err != nil {
		return nil, nil, err
	}

	return userTodos, meta, nil
}

func (r *TodoRepository) Find(userId int64, itemId int64) (map[string]any, error) {

	item, err := xqb.Table("todos").Where("id", "=", itemId).First()

	if err != nil {
		return nil, err
	}

	return item, nil
}

func (r *TodoRepository) Create(userId int64, items []map[string]any) (int64, error) {

	createdItemId, err := xqb.Table("todos").InsertGetId(items)

	if err == nil {
		return 0, err
	}

	return createdItemId, nil
}

func (r *TodoRepository) Update(userId int64, itemId int64, updatedFields map[string]any) (int64, error) {

	updatedFields["updated_at"] = time.Now()

	rowsAffected, err := xqb.Table("todos").Where("id", "=", itemId).Where("user_id", "=", userId).Update(updatedFields)
	if err != nil {
		logger.Log().Error(err.Error())
		return 0, err
	}

	return rowsAffected, nil
}

func (r *TodoRepository) Delete(userId int64, itemId int64) (int64, error) {

	rowsAffected, err := xqb.Table("todos").Where("user_id", "=", userId).Where("id", "=", itemId).Delete()

	if err != nil {
		logger.Log().Error(err.Error())
		return 0, err
	}

	return rowsAffected, nil
}
