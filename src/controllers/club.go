package controllers

import (
	"context"
	"find-club-graphql/constant"
	"find-club-graphql/models"
	"github.com/graphql-go/graphql"
	gh "github.com/graphql-go/handler"
	"net/http"
)

var (
	handler *gh.Handler

	query = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"club": &graphql.Field{
				Name:        "club",
				Type:        clubType,
				Args:        clubArgs,
				Resolve:     getClub,
				Description: "获取一个社团信息",
			},
			"clubs": &graphql.Field{
				Name:        "clubs",
				Type:        graphql.NewList(clubType),
				Args:        clubsArgs,
				Resolve:     getClubs,
				Description: "获取多个社团信息",
			},
		},
		IsTypeOf:    nil,
		Description: "",
	})
	clubType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Club",
		Description: "Club",
		Interfaces:  nil,
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type:        graphql.ID,
				Description: "社团ID",
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if club, ok := p.Source.(models.Club); ok {
						return club.ID.Hex(), nil
					}
					return nil, constant.ErrorEmpty
				},
			},
			"name": &graphql.Field{
				Type:        graphql.NewNonNull(graphql.String),
				Description: "社团名称",
			},
			"introduction": &graphql.Field{
				Type:        graphql.String,
				Description: "社团简介",
			},
			"detailInfo": &graphql.Field{
				Type:        graphql.String,
				Description: "社团详细信息",
			},
			"category": &graphql.Field{
				Type:        categoryType,
				Description: "社团分类",
			},
			"clubType": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "社团类型",
			},
			"tag": &graphql.Field{
				Type:        graphql.NewList(graphql.String),
				Description: "社团标签",
			},
			"searchCount": &graphql.Field{
				Type:        graphql.Int,
				Description: "搜索次数",
			},
		},
	})
	clubArgs = graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.ID),
			Description: "社团ID",
		},
	}
	clubsArgs = graphql.FieldConfigArgument{
		"category": &graphql.ArgumentConfig{
			Type:        categoryType,
			Description: "社团分类：校级组织或兴趣社团",
		},
		"clubType": &graphql.ArgumentConfig{
			Type:        graphql.NewList(graphql.String),
			Description: `"社团类型，如： "",""`,
		},
		"tag": &graphql.ArgumentConfig{
			Type:        graphql.NewList(graphql.String),
			Description: `社团TAG，如：  "志愿服务","外联","管理","宣传","新媒体"等`,
		},
		"keyWord": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "关键字，按社团名关键字搜索",
		},
		"hot": &graphql.ArgumentConfig{
			Type:         graphql.Boolean,
			DefaultValue: false,
			Description:  "获得热搜社团",
		},
	}

	categoryType = graphql.NewEnum(graphql.EnumConfig{
		Name:        "Category",
		Description: "club分类",
		Values: graphql.EnumValueConfigMap{
			"school": &graphql.EnumValueConfig{
				Value:       constant.CategorySchool,
				Description: "校级组织",
			},
			"student": &graphql.EnumValueConfig{
				Value:       constant.CategoryStudent,
				Description: "兴趣社团",
			},
		},
	})

	mutation = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Mutation",
		Description: "增删查改等操作",
		Fields: graphql.Fields{
			"addClub": &graphql.Field{
				Name:        "addClub",
				Type:        clubType,
				Args:        addClubArgs,
				Resolve:     addClub,
				Description: "新增社团信息",
			},
		},
	})
	addClubArgs = graphql.FieldConfigArgument{
		"name": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(graphql.String),
			Description: "社团名称",
		},
		"introduction": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "简介",
		},
		"detailInfo": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "详细信息",
		},
		"category": &graphql.ArgumentConfig{
			Type:        graphql.NewNonNull(categoryType),
			Description: "社团分类：校级组织或兴趣社团",
		},
		"clubType": &graphql.ArgumentConfig{
			Type:        graphql.NewList(graphql.String),
			Description: `"社团类型，如： "",""`,
		},
		"tag": &graphql.ArgumentConfig{
			Type:        graphql.NewList(graphql.String),
			Description: `社团TAG，如：  "志愿服务","外联","管理","宣传","新媒体"等`,
		},
	}
)

func init() {
	schemaConfig := graphql.SchemaConfig{
		Query:    query,
		Mutation: mutation,
	}
	schema, _ := graphql.NewSchema(schemaConfig)
	handler = gh.New(&gh.Config{
		Schema:           &schema,
		Pretty:           false,
		GraphiQL:         false,
		Playground:       false,
		RootObjectFn:     nil,
		ResultCallbackFn: nil,
		FormatErrorFn:    nil,
	})
}
func Graphql(w http.ResponseWriter, r *http.Request) {
	handler.ContextHandler(context.Background(), w, r)
}

func getClub(p graphql.ResolveParams) (interface{}, error) {
	id, ok := p.Args["id"].(string)
	if !ok {
		return nil, constant.ErrorEmpty
	}

	return models.GetClub(id)

}

func getClubs(p graphql.ResolveParams) (interface{}, error) {
	if keyWord, ok := p.Args["keyWord"].(string); ok {
		return models.GetClubsByKeyword(keyWord)
	}
	if hot, ok := p.Args["hot"].(bool); ok && hot {
		return models.GetHotClubs(6)
	}

	params := make(map[string]interface{})

	if category, ok := p.Args["category"].(string); ok {
		params["category"] = category
	}
	if clubType, ok := p.Args["clubType"].([]interface{}); ok {
		params["club_type"] = clubType
	}
	if tag, ok := p.Args["tag"].(string); ok {
		params["tag"] = tag
	}

	return models.GetClubs(params)
}
func addClub(p graphql.ResolveParams) (interface{}, error) {
	club := models.Club{}
	if name, ok := p.Args["name"].(string); ok {
		club.Name = name
	}
	if introduction, ok := p.Args["introduction"].(string); ok {
		club.Introduction = introduction
	}
	if detailInfo, ok := p.Args["detailInfo"].([]interface{}); ok {
		for _, v := range detailInfo {
			club.DetailInfo = append(club.DetailInfo, v.(string))
		}
	}
	if category, ok := p.Args["category"].(string); ok {
		club.Category = category
	}
	if clubType, ok := p.Args["clubType"].([]interface{}); ok {
		for _, v := range clubType {
			club.ClubType = append(club.ClubType, v.(string))
		}
	}
	if tag, ok := p.Args["tag"].([]interface{}); ok {
		for _, v := range tag {
			club.Tag = append(club.Tag, v.(string))
		}
	}

	return models.AddClub(club)

}
