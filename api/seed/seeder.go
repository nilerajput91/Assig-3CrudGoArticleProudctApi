package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/nilerajput91/Assig-3CrudGoArticleProudctApi/api/models"
)

var users = []models.User{
	models.User{
		Nickname: "neel rajput",
		Email:    "nilerajput91@gmail.com",
		Password: "pass",
	},
	models.User{
		Nickname: "kam",
		Email:    "nilesh.rajput@gmail.com",
		Password: "pass",
	},
}

var articles = []models.Article{
	models.Article{
		Title:   "java",
		Content: "basic of  java",
	},
	models.Article{
		Title:   "python",
		Content: "rest api in py",
	},
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Article{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Article{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Article{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i, _ := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		articles[i].AuthorID = users[i].ID

		err = db.Debug().Model(&models.Article{}).Create(&articles[i]).Error
		if err != nil {
			log.Fatalf("cannot seed posts table: %v", err)
		}
	}
}
