package model

import "time"

type Delivery struct {
	ID			int		`json:"id" gorm:"primaryKey"`
	OrderNum	string	`json:"order_num" gorm:"type:varchar(32);not null"`

	SrcName		string	`json:"src_name" gorm:"type:varchar(32);not null"`
	SrcPhone	string	`json:"src_phone" gorm:"type:varchar(32);not null"`
	SrcNodeID	int		`json:"src_node_id" gorm:"not null"`
	SrcNode		Node	`json:"src_node" gorm:"foreignKey:SrcNodeID"`

	DestName	string	`json:"dest_name" gorm:"type:varchar(32);not null"`
	DestPhone	string	`json:"dest_phone" gorm:"type:varchar(32);not null"`
	DestNodeID	int		`json:"dest_node_id" gorm:"not null"`
	DestNode	Node	`json:"dest_node" gorm:"foreignKey:DestNodeID"`
	
	CreatedAt	time.Time `json:"created_at" gorm:"not null"`
}

func (Delivery) TableName() string {
	return "deliveries"
}
