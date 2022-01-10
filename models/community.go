package models

import "time"

type Community struct {
	ID   int64  `json:"id" db:"community_id"`
	Name string `json:"name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64  `json:"id" db:"community_id"`
	Name         string `json:"name" db:"community_name"`
	Introduction string `json:"introduction,omitempty" db:"introduction"` // 详情，omitemoty当字段为空时，则不展示
	CreateTime   time.Time `json:"create_time" db:"create_time"`
}
