package mysql

import (
	"gitea.bjx.cloud/allstar/common/core/logger"
	"upper.io/db.v3/lib/sqlbuilder"
)

func Insert(obj Domain) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(err)
		}
	}()
	if err != nil {
		return err
	}
	_, err = conn.Collection(obj.TableName()).Insert(obj)
	if err != nil {
		return err
	}
	return nil
}

func TransInsert(tx sqlbuilder.Tx, obj Domain) error {
	_, err := tx.Collection(obj.TableName()).Insert(obj)
	if err != nil {
		return err
	}
	return nil
}

func TransBatchInsert(tx sqlbuilder.Tx, obj Domain, objs []interface{}) error {

	//a := objs.([]interface{})

	batch := tx.InsertInto(obj.TableName()).Batch(len(objs))
	go func() {
		defer batch.Done()
		for i := range objs {
			batch.Values(objs[i])
		}
	}()
	err := batch.Wait()
	if err != nil {
		return err
	}

	return nil
}

func BatchInsert(obj Domain, objs []interface{}) error {
	conn, err := GetConnect()
	defer func() {
		if err := conn.Close(); err != nil {
			logger.GetDefaultLogger().Info(err)
		}
	}()
	if err != nil {
		log.Error(err)
		return err
	}

	batch := conn.InsertInto(obj.TableName()).Batch(len(objs))
	go func() {
		defer batch.Done()
		for i := range objs {
			batch.Values(objs[i])
		}
	}()
	err = batch.Wait()
	if err != nil {
		return err
	}

	return nil
}

//func BatchDone(pos []interface{}, batch *sqlbuilder.BatchInserter) {
//		defer batch.Done()
//		for i := range pos {
//			batch.Values(pos[i])
//		}
//}
