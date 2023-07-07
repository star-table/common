package mysql

import (
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/util/strs"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TransX(txFunc func(tx sqlbuilder.Tx) error) error {
	conn, err := GetConnect()
	if err != nil {
		log.Error(consts.DBOpenErrorSentence + strs.ObjectToString(err))
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Error(consts.TxOpenErrorSentence + strs.ObjectToString(err))
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	defer Close(conn, tx)

	err = txFunc(tx)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		Rollback(tx)
		return err
	}

	err = tx.Commit()

	if err != nil {
		log.Error("tx.Commit(): " + strs.ObjectToString(err))
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	return nil
}
