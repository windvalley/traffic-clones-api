package main

import (
	"errors"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type githubTraffic struct {
	gorm.Model
	User      string `gorm:"index:idx_user_reponame_type_timestamp"`
	RepoName  string `gorm:"column:repo_name;index:idx_user_reponame_type_timestamp"`
	Type      string `gorm:"index:idx_user_reponame_type_timestamp"`
	Uniques   int
	Count     int
	Timestamp string `gorm:"index:idx_user_reponame_type_timestamp"`
}

func (g *githubTraffic) TableName() string {
	return "github_traffic"
}

var (
	db     *gorm.DB
	dbfile = "./github_traffic.db"
)

func init() {
	var err error

	db, err = gorm.Open(sqlite.Open(dbfile), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&githubTraffic{}); err != nil {
		panic(err)
	}
}

type clonesTotal struct {
	Count   int
	Uniques int
}

func updateGithubTrafficClones(githubClones []clonesItem, user, repoName string) {
	for _, v := range githubClones {
		var record githubTraffic

		err := db.Where(
			"user = ? and repo_name = ? and type = ? and timestamp = ?",
			user,
			repoName,
			trafficClonesLabel,
			v.Timestamp,
		).First(&record).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			recordNew := &githubTraffic{
				User:      user,
				RepoName:  repoName,
				Type:      trafficClonesLabel,
				Uniques:   v.Uniques,
				Count:     v.Count,
				Timestamp: v.Timestamp,
			}
			if err := db.Create(&recordNew).Error; err != nil {
				log.Printf("create %v failed: %v", *recordNew, err)
			}
			continue
		}

		record.Uniques = v.Uniques
		record.Count = v.Count
		if err := db.Save(&record).Error; err != nil {
			log.Printf("update %v failed: %v", record, err)
		}
	}
}

func getClonesTotal(user, repoName string) *clonesTotal {
	var total clonesTotal

	db.Raw(
		"select sum(count) as count, sum(uniques) as uniques from github_traffic where user=? and repo_name=? and type=?",
		user,
		repoName,
		trafficClonesLabel,
	).Scan(&total)

	return &total
}
