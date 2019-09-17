package mysql

import (
	"gitea.bjx.cloud/allstar/common/core/consts"
	"gitea.bjx.cloud/allstar/common/core/errors"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TransX(txFunc func(tx sqlbuilder.Tx) error) error {
	conn, err := GetConnect()
	if err != nil {
		log.Errorf(consts.DBOpenErrorSentence, err)
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	tx, err := conn.NewTx(nil)
	if err != nil {
		log.Errorf(consts.TxOpenErrorSentence, err)
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	defer Close(conn, tx)

	err = txFunc(tx)
	if err != nil {
		log.Error(err)
		Rollback(tx)
		return err
	}

	err = tx.Commit()

	if err != nil {
		log.Errorf("tx.Commit(): %q\n", err)
		return errors.BuildSystemErrorInfo(errors.MysqlOperateError, err)
	}
	return nil
}
