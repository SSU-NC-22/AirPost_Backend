package model

type Delivery struct {
	ID			int		`json:"id" gorm:"primaryKey"`
	SrcName		string	`json:"src_name" gorm:"type:varchar(32);not null"`
	SrcPhone	string	`json:"src_phone" gorm:"type:varchar(32);not null"`
	DestName	string	`json:"dest_name" gorm:"type:varchar(32);not null"`
	DestPhone	string	`json:"dest_phone" gorm:"type:varchar(32);not null"`
	NodeID		int		`json:"node_id" gorm:"not null"`
	Node		Node	`json:"node" gorm:"foreignKey:NodeID"`
}

func (Delivery) TableName() string {
	return "deliveries"
}
