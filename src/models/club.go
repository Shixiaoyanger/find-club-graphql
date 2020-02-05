package models

import (
	"find-club-graphql/constant"
	"find-club-graphql/dao"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

//Club is a struct that describes club
type Club struct {
	ID           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name         string        `bson:"name" json:"name"`
	Introduction string        `bson:"introduction" json:"introduction"`
	DetailInfo   []string      `bson:"detail_info" json:"detailInfo"`
	Category     string        `bson:"category" json:"category"`
	ClubType     []string      `bson:"club_type" json:"clubType"`
	Tag          []string      `bson:"tag" json:"tag"`
	SearchCount  int           `bson:"search_count" json:"searchCount"`
}

func AddClub(club Club) (Club, error) {
	club.ID = bson.NewObjectId()
	return addClub(club)
}

//GetAllClubs return all clubs in db.
func GetAllClubs() ([]Club, error) {
	query := bson.M{}

	return getClubs(query)
}
func GetClub(id string) (Club, error) {
	if !bson.IsObjectIdHex(id) {
		return Club{}, constant.ErrorParamsWrong
	}
	query := bson.M{
		"_id": bson.ObjectIdHex(id),
	}
	return getClub(query)
}

func GetClubsByClubType(typ string) ([]Club, error) {
	query := bson.M{
		"club_type": typ,
	}
	query["pp"] = "p"
	fmt.Println(query)
	return getClubs(query)
}

func GetClubsByKeyword(keyword string) ([]Club, error) {

	query := bson.M{
		"name": bson.M{
			"$regex": bson.RegEx{
				Pattern: keyword,
			},
		},
	}
	update := bson.M{
		"$inc": bson.M{
			"search_count": 1,
		},
	}

	go addSearchCount(query, update)
	return getClubs(query)
}

func GetClubs(params map[string]interface{}) ([]Club, error) {
	query := bson.M{}
	for k, v := range params {
		query[k] = v
	}
	return getClubs(query)
}

func GetHotClubs(limit int) ([]Club, error) {
	sort := []string{"-search_count"}
	return getClubs2(nil, nil, 0, limit, sort)

}
func getClub(query bson.M) (Club, error) {
	cntrl := dao.NewCloneMgoDBCntlr()
	defer cntrl.Close()

	var club Club
	table := cntrl.GetTable(constant.TableClub)
	err := table.Find(query).One(&club)
	return club, err
}

func getClubs(query bson.M) ([]Club, error) {
	cntrl := dao.NewCloneMgoDBCntlr()
	defer cntrl.Close()

	var clubs []Club
	table := cntrl.GetTable(constant.TableClub)
	err := table.Find(query).All(&clubs)
	return clubs, err
}

func getClubs2(query interface{}, selector interface{}, skip int, limit int, sort []string) ([]Club, error) {
	cntrl := dao.NewCloneMgoDBCntlr()
	defer cntrl.Close()

	var clubs []Club
	q := cntrl.GetTable(constant.TableClub).Find(query).Select(selector)

	if skip != 0 {
		q = q.Skip(skip)
	}
	if limit != 0 {
		q = q.Limit(limit)
	}
	if sort != nil {
		q = q.Sort(sort...)
	}
	err := q.All(&clubs)

	return clubs, err
}
func addClub(club Club) (Club, error) {
	cntrl := dao.NewCloneMgoDBCntlr()
	defer cntrl.Close()

	table := cntrl.GetTable(constant.TableClub)
	err := table.Insert(club)
	return club, err
}
func addSearchCount(selector, update bson.M) {
	cntrl := dao.NewCloneMgoDBCntlr()
	defer cntrl.Close()

	table := cntrl.GetTable(constant.TableClub)
	err := table.Update(selector, update)
	fmt.Println(err)
}
