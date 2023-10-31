package storage

import (
	"context"
	"ilserver/domain"
)

type Storage interface {
	RecordCountInTable(ctx context.Context, name string) (int, error)

	// *** admin ***

	InsertAdmin(ctx context.Context, admin domain.Admin) (int64, error)
	HasAdminByLogin(ctx context.Context, login string) (bool, error)
	HasAdminWithLoginAndPassword(ctx context.Context, login, password string) (bool, error)
	AllAdmins(ctx context.Context) (domain.AdminList, error)
	AdminByLogin(ctx context.Context, login string) (domain.Admin, error)
	UpdateAdminPasswordByLogin(ctx context.Context, login, password string) error
	DeleteAdminByLogin(ctx context.Context, login string) error
	DeleteAdmins(ctx context.Context) error

	// *** topic ***

	InsertTopic(ctx context.Context, topic domain.Topic) (int64, error)
	InsertTopics(ctx context.Context, topics domain.TopicList) error
	AllTopics(ctx context.Context) (domain.TopicList, error)
	Topic(ctx context.Context, idr int) (domain.Topic, error)
	RandomTopic(ctx context.Context, lang int) (domain.Topic, error)
	DeleteTopic(ctx context.Context, idr int) error
	DeleteTopics(ctx context.Context) error

	//...
}
