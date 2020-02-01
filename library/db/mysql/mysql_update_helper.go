package mysql

import (
	"github.com/galaxy-book/common/core/consts"
	"github.com/galaxy-book/common/core/logger"
	"github.com/galaxy-book/common/core/util/strs"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func Update(obj Domain) error {
	conn, err := GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return err
	}
	err = conn.Collection(obj.TableName()).UpdateReturning(obj)
	if err != nil {
		return err
	}
	return nil
}

func TransUpdate(tx sqlbuilder.Tx, obj Domain) error {
	err := tx.Collection(obj.TableName()).UpdateReturning(obj)
	if err != nil {
		return err
	}
	return nil
}

func UpdateSmart(table string, id int64, upd Upd) error {
	_, err := UpdateSmartWithCond(table, db.Cond{
		consts.TcId: id,
	}, upd)

	return err
}

func UpdateSmartWithCond(table string, cond db.Cond, upd Upd) (int64, error) {
	conn, err := GetConnect()
	defer func() {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.GetDefaultLogger().Info(strs.ObjectToString(err))
			}
		}
	}()
	if err != nil {
		return 0, err
	}
	res, err := conn.Update(table).Set(upd).Where(cond).Exec()
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return row, nil
}

func TransUpdateSmart(tx sqlbuilder.Tx, table string, id int64, upd Upd) error {
	_, err := TransUpdateSmartWithCond(tx, table, db.Cond{
		consts.TcId: id,
	}, upd)
	if err != nil {
		log.Error(strs.ObjectToString(err))
	}
	return err
}

func TransUpdateSmartWithCond(tx sqlbuilder.Tx, table string, cond db.Cond, upd Upd) (int64, error) {
	res, err := tx.Update(table).Set(upd).Where(cond).Exec()
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return 0, err
	}
	row, err := res.RowsAffected()
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return 0, err
	}

	return row, nil
}
