package cmd

import (
	"github.com/jihanlugas/pos/cryption"
	"github.com/jihanlugas/pos/db"
	"github.com/jihanlugas/pos/model"
	"github.com/jihanlugas/pos/utils"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
	"time"
)

var dbCmd = &cobra.Command{
	Use:   "db",
	Short: "Run server",
	Long: `With this command you can
	up : create database table
	down :  drop database table
	seed :  insert data table
	`,
}

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Up table",
	Long:  "Up table",
	Run: func(cmd *cobra.Command, args []string) {
		up()
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Down table",
	Long:  "Down table",
	Run: func(cmd *cobra.Command, args []string) {
		down()
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed data table",
	Long:  "Seed data table",
	Run: func(cmd *cobra.Command, args []string) {
		seed()
	},
}

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Down, up, seed table",
	Long:  "Down, up, seed table",
	Run: func(cmd *cobra.Command, args []string) {
		down()
		up()
		seed()
	},
}

func init() {
	rootCmd.AddCommand(dbCmd)
	dbCmd.AddCommand(upCmd)
	dbCmd.AddCommand(downCmd)
	dbCmd.AddCommand(resetCmd)
	dbCmd.AddCommand(seedCmd)
}

func up() {
	var err error
	conn, closeConn := db.GetConnection()
	defer closeConn()

	// table
	err = conn.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Company{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Usercompany{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Companysettings{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Item{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().AutoMigrate(&model.Itemvariant{})
	if err != nil {
		panic(err)
	}

	// view
	vUser := conn.Model(&model.User{}).
		Select("users.*, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("left join users u1 on u1.id = users.create_by").
		Joins("left join users u2 on u2.id = users.update_by").
		Joins("left join users u3 on u3.id = users.delete_by")

	vCompany := conn.Model(&model.Company{}).
		Select("companies.*, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("left join users u1 on u1.id = companies.create_by").
		Joins("left join users u2 on u2.id = companies.update_by").
		Joins("left join users u3 on u3.id = companies.delete_by")

	vUsercompany := conn.Model(&model.Usercompany{}).
		Select("usercompanies.*, users.fullname as user_name, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join users on users.id = usercompanies.user_id").
		Joins("join companies on companies.id = usercompanies.company_id ").
		Joins("left join users u1 on u1.id = usercompanies.create_by").
		Joins("left join users u2 on u2.id = usercompanies.update_by").
		Joins("left join users u3 on u3.id = usercompanies.delete_by")

	vCompanysetting := conn.Model(&model.Companysettings{}).
		Select("companysettings.*, companies.name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join companies on companies.id = companysettings.id ").
		Joins("left join users u1 on u1.id = companysettings.create_by").
		Joins("left join users u2 on u2.id = companysettings.update_by").
		Joins("left join users u3 on u3.id = companysettings.delete_by")

	vItem := conn.Model(&model.Item{}).
		Select("items.*, companies.name as company_name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join companies on companies.id = items.company_id").
		Joins("left join users u1 on u1.id = items.create_by").
		Joins("left join users u2 on u2.id = items.update_by").
		Joins("left join users u3 on u3.id = items.delete_by")

	vItemvariant := conn.Model(&model.Itemvariant{}).
		Select("itemvariants.*, items.name as item_name, u1.fullname as create_name, u2.fullname as update_name, u3.fullname as delete_name").
		Joins("join items on items.id = itemvariants.item_id").
		Joins("left join users u1 on u1.id = itemvariants.create_by").
		Joins("left join users u2 on u2.id = itemvariants.update_by").
		Joins("left join users u3 on u3.id = itemvariants.delete_by")

	err = conn.Migrator().CreateView("users_view", gorm.ViewOption{
		Replace: true,
		Query:   vUser,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("companies_view", gorm.ViewOption{
		Replace: true,
		Query:   vCompany,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("usercompanies_view", gorm.ViewOption{
		Replace: true,
		Query:   vUsercompany,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("companysettings_view", gorm.ViewOption{
		Replace: true,
		Query:   vCompanysetting,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("items_view", gorm.ViewOption{
		Replace: true,
		Query:   vItem,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().CreateView("itemvariants_view", gorm.ViewOption{
		Replace: true,
		Query:   vItemvariant,
	})
	if err != nil {
		panic(err)
	}

}

func down() {
	var err error
	conn, closeConn := db.GetConnection()
	defer closeConn()

	// view
	err = conn.Migrator().DropView("users_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("companies_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("usercompanies_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("companysettings_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("items_view")
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropView("itemvariants_view")
	if err != nil {
		panic(err)
	}

	// table
	err = conn.Migrator().DropTable(&model.User{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Company{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Usercompany{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Companysettings{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Item{})
	if err != nil {
		panic(err)
	}
	err = conn.Migrator().DropTable(&model.Itemvariant{})
	if err != nil {
		panic(err)
	}

}

func seed() {
	now := time.Now()
	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}

	conn, closeConn := db.GetConnection()
	defer closeConn()

	tx := conn.Begin()

	users := []model.User{
		{RoleID: utils.GetUniqueID(), Email: "jihanlugas2@gmail.com", Username: "jihanlugas", NoHp: "6287770333043", Fullname: "Jihan Lugas", Passwd: password, PassVersion: 1, IsActive: true, LastLoginDt: nil, PhotoID: "", CreateBy: "", CreateDt: now, UpdateBy: "", UpdateDt: now},
	}
	tx.Create(&users)

	companies := []model.Company{
		{Name: "My Corporation", CreateBy: users[0].ID, CreateDt: now, UpdateBy: users[0].ID, UpdateDt: now},
	}
	tx.Create(&companies)

	usercompanies := []model.Usercompany{
		{UserID: users[0].ID, CompanyID: companies[0].ID, IsDefaultCompany: true, IsCreator: true, CreateBy: users[0].ID, CreateDt: now, UpdateBy: users[0].ID, UpdateDt: now},
	}
	tx.Create(&usercompanies)

	companysettings := []model.Companysettings{
		{ID: companies[0].ID, CreateBy: users[0].ID, CreateDt: now, UpdateBy: users[0].ID, UpdateDt: now},
	}
	tx.Create(&companysettings)

	err = tx.Commit().Error
	if err != nil {
		panic(err)
	}

}
