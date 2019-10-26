package migrate

import (
	"fmt"
	"github.com/go-pg/migrations"
	"github.com/go-pg/pg"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"log"
)

// Migrate runs go-pg migrations
func Migrate(env *models.ExtenderEnvironment) {
	db := pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%d", env.DbHost, env.DbPort),
		User:            env.DbUser,
		Password:        env.DbPassword,
		Database:        env.DbName,
		ApplicationName: env.AppName,
		MinIdleConns:    env.DbMinIdleConns,
		PoolSize:        env.DbPoolSize,
		MaxRetries:      10,
	})
	defer db.Close()

	oldVersion, newVersion, err := migrations.Run(db, "up")
	if err != nil {
		log.Println(err.Error())
	}
	if newVersion != oldVersion {
		fmt.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	} else {
		fmt.Printf("version is %d\n", oldVersion)
	}

	//err := db.RunInTransaction(func(tx *pg.Tx) error {
	//	oldVersion, newVersion, err := migrations.Run(db, "init")
	//	if err != nil {
	//		return err
	//	}
	//	if newVersion != oldVersion {
	//		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	//	} else {
	//		log.Printf("version is %d\n", oldVersion)
	//	}
	//	return nil
	//})
	//
	//log.Println(err)
	//
	//err = db.RunInTransaction(func(tx *pg.Tx) error {
	//	oldVersion, newVersion, err := migrations.Run(db, "up")
	//	if err != nil {
	//		return err
	//	}
	//	if newVersion != oldVersion {
	//		log.Printf("migrated from version %d to %d\n", oldVersion, newVersion)
	//	} else {
	//		log.Printf("version is %d\n", oldVersion)
	//	}
	//	return nil
	//})

	if err != nil {
		log.Println(err)
	}
}
