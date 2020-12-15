package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

type Article struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (p *Article) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

func (p *Article) Validate() error {

	if p.Title == "" {
		return errors.New("Required Title")
	}

	if p.Content == "" {
		return errors.New("Required Content")
	}
	if p.AuthorID < 1 {
		return errors.New("Required Author")
	}
	return nil
}

func (p *Article) SaveArticle(db *gorm.DB) (*Article, error) {
	var err error
	err = db.Debug().Model(&Article{}).Create(&p).Error
	if err != nil {
		return &Article{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Article{}, err
		}
	}
	return p, nil
}

func (p *Article) FindAllArticles(db *gorm.DB) (*[]Article, error) {
	var err error
	articles := []Article{}
	err = db.Debug().Model(&Article{}).Limit(100).Find(&articles).Error
	if err != nil {
		return &[]Article{}, err
	}
	if len(articles) > 0 {
		for i, _ := range articles {
			err := db.Debug().Model(&User{}).Where("id = ?", articles[i].AuthorID).Take(&articles[i].Author).Error
			if err != nil {
				return &[]Article{}, err
			}
		}
	}
	return &articles, nil
}

func (p *Article) FindArticleByID(db *gorm.DB, pid uint64) (*Article, error) {
	var err error
	err = db.Debug().Model(&Article{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Article{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Article{}, err
		}
	}
	return p, nil
}

func (p *Article) UpdateAArticle(db *gorm.DB) (*Article, error) {

	var err error

	err = db.Debug().Model(&Article{}).Where("id = ?", p.ID).Updates(Article{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Article{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Article{}, err
		}
	}
	return p, nil
}

func (p *Article) DeleteAArticle(db *gorm.DB, pid uint64, uid uint32) (int64, error) {

	db = db.Debug().Model(&Article{}).Where("id = ? and author_id = ?", pid, uid).Take(&Article{}).Delete(&Article{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New(" Article not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
