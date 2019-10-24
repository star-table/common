package mysql

import (
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/logger"
	"gitea.bjx.cloud/allstar/common/core/util/strs"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func SelectById(table string, id interface{}, obj interface{}) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return err
	}
	err = conn.Collection(table).Find(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func TransSelectById(tx sqlbuilder.Tx, table string, id interface{}, obj interface{}) error {
	err := tx.Collection(table).Find(db.Cond{
		consts.TcId:       id,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func SelectCountByCond(table string, cond db.Cond) (uint64, error) {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return 0, err
	}
	unit, err := conn.Collection(table).Find(cond).Count()
	if err != nil {
		return 0, err
	}
	return unit, nil
}

func TransSelectCountByCond(tx sqlbuilder.Tx, table string, cond db.Cond) (uint64, error) {
	unit, err := tx.Collection(table).Find(cond).Count()
	if err != nil {
		return 0, err
	}
	return unit, nil
}

func SelectOneByCond(table string, cond db.Cond, obj interface{}) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return err
	}
	err = conn.Collection(table).Find(cond).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func TransSelectOneByCond(tx sqlbuilder.Tx, table string, cond db.Cond, obj interface{}) error {
	err := tx.Collection(table).Find(cond).One(obj)
	if err != nil {
		return err
	}
	return nil
}

func SelectByQuery(query string, objs interface{}, args ...interface{}) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return err
	}
	var iter sqlbuilder.Iterator = nil
	if len(args) > 0 {
		iter = conn.Iterator(query, args...)
	} else {
		iter = conn.Iterator(query)
	}
	err = iter.All(objs)
	return err
}

func TransSelectByQuery(tx sqlbuilder.Tx, query string, objs interface{}, args ...interface{}) error {
	var iter sqlbuilder.Iterator = nil
	if len(args) > 0 {
		iter = tx.Iterator(query, args...)
	} else {
		iter = tx.Iterator(query)
	}
	err := iter.All(objs)
	return err
}

func SelectAllByCond(table string, cond db.Cond, objs interface{}) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return err
	}
	err = conn.Collection(table).Find(cond).All(objs)
	if err != nil {
		return err
	}
	return nil
}

func TransSelectAllByCond(tx sqlbuilder.Tx, table string, cond db.Cond, objs interface{}) error {
	err := tx.Collection(table).Find(cond).All(objs)
	if err != nil {
		return err
	}
	return nil
}

func SelectAllByCondWithPageAndOrder(table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) (uint64, error) {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return 0, err
	}

	mid := conn.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	count, err := mid.TotalEntries()
	if err != nil {
		return 0, err
	}
	err = mid.All(objs)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func TransSelectAllByCondWithPageAndOrder(tx sqlbuilder.Tx, table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) (uint64, error) {
	mid := tx.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	count, err := mid.TotalEntries()
	if err != nil {
		return 0, err
	}
	err = mid.All(objs)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func SelectAllByCondWithNumAndOrder(table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			log.Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return err
	}

	mid := conn.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	err = mid.All(objs)
	if err != nil {
		return err
	}
	return nil
}

func TransSelectAllByCondWithNumAndOrder(tx sqlbuilder.Tx, table string, cond db.Cond, union *db.Union, page int, size int, order interface{}, objs interface{}) error {
	mid := tx.Collection(table).Find(cond)
	if union != nil {
		mid = mid.And(union)
	}
	if size > 0 && page > 0 {
		mid = mid.Page(uint(page)).Paginate(uint(size))
	}
	if order != nil && order != "" {
		mid = mid.OrderBy(order)
	}
	err := mid.All(objs)
	if err != nil {
		return err
	}
	return nil
}

func IsExistByCond(table string, cond db.Cond) (bool, error) {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(strs.ObjectToString(err))
		}
	}()
	if err != nil {
		return false, err
	}
	exist, err := conn.Collection(table).Find(cond).Exists()
	if err != nil {
		return false, err
	}
	return exist, nil
}
