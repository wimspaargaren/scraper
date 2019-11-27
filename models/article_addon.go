package models

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/goadesign/goa"
	"github.com/jinzhu/gorm"
)

// UpdatePDFFound modifies a single record.
func (m *ArticleDB) UpdatePDFFound(ctx context.Context, model *Article) error {
	defer goa.MeasureSince([]string{"goa", "db", "article", "update"}, time.Now())
	model.GotPdf = true
	obj, err := m.Get(ctx, model.ID)
	if err != nil {
		goa.LogError(ctx, "error updating Article", "error", err.Error())
		return err
	}
	err = m.Db.Model(obj).Updates(model).Error

	return err
}

// ListNoDOI returns an array of Article
func (m *ArticleDB) ListNoDOI(ctx context.Context) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "ListNoDOI"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("processed = 1 and (doi = '' OR doi IS NULL)").Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// ListDOILinks returns an array of Article
func (m *ArticleDB) ListDOILinks(ctx context.Context) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "ListDOILinks"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("doi LIKE '%http%'").Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// Count number of total articles
type Count struct {
	Count int
}

// ListArticles returns an array of Article
func (m *ArticleDB) ListArticles(ctx context.Context, statuses []Status, page int, search *string) ([]*Article, int, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "ListArticles"}, time.Now())
	additionalSearch := ""
	if search != nil {
		temp := strings.ToLower(*search)
		additionalSearch = " AND (lower(title) LIKE '%" + temp + "%' OR lower(doi) LIKE '%" + temp + "%' OR year = " + temp + ")"
	}
	var count Count
	err := m.Db.Table(m.TableName()).Select("COUNT(id)").Where("processed IN (?)"+additionalSearch, statuses).
		Find(&count).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	limit := 20
	var objs []*Article
	err = m.Db.Table(m.TableName()).Where("processed IN (?)"+additionalSearch, statuses).
		Limit(limit).
		Offset(page * limit).
		Order("doi").
		Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}

	return objs, count.Count, nil
}

// ListOnStatus returns an array of Article
func (m *ArticleDB) ListOnStatus(ctx context.Context, status Status) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "ListArticles"}, time.Now())

	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("processed = ?", status).
		Order("doi").
		Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

// ListOnStatus returns an array of Article
func (m *ArticleDB) ListOnDoi(ctx context.Context, doi string) ([]*Article, error) {
	defer goa.MeasureSince([]string{"goa", "db", "article", "ListArticles"}, time.Now())
	doiLowered := strings.ToLower(doi)
	var objs []*Article
	err := m.Db.Table(m.TableName()).Where("lower(doi) = ?", doiLowered).
		Order("url").
		Find(&objs).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return objs, nil
}

type Keywords struct {
	List []string `json:"list"`
}

// AddKeywords add keywords to an article
func (a *Article) AddKeywords(keywords Keywords) error {
	b, err := json.Marshal(keywords)
	if err != nil {
		return err
	}
	a.Keywords = b
	return nil

}

// AddKeywords add keywords to an article
func (a *Article) GetKeywords() (Keywords, error) {
	var res Keywords
	err := json.Unmarshal(a.Keywords, &res)
	return res, err
}
