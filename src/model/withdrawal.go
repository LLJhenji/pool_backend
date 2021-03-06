package model

import (
	"pool_backend/src/enum"
	"pool_backend/src/global"
	"time"
)

// Withdrawal model
type Withdrawal struct {
	ID            string    `gorm:"column:id;primary_key" json:"id" form:"id"`
	ShareholderID string    `gorm:"column:shareholder_id" json:"shareholder_id" form:"shareholder_id"`
	Amount        float64   `gorm:"column:amount" json:"amount" form:"amount"`
	State         int64     `gorm:"column:state" json:"state" form:"state"`
	Content       string    `gorm:"column:content" json:"content" form:"content"`
	WithLog       string    `gorm:"column:with_log" json:"with_log" form:"with_log"`
	Hash          string    `gorm:"column:hash" json:"hash" form:"hash"`
	CreateAt      time.Time `gorm:"column:create_at" json:"create_at" form:"create_at"`
	EndAt         time.Time `gorm:"column:end_at" json:"end_at" form:"end_at" `
}

//UserWithdrawal model
type UserWithdrawal struct {
	ID                      string    `gorm:"column:id;primary_key" json:"id" form:"id"`
	ShareholderID           string    `gorm:"column:shareholder_id" json:"shareholder_id" form:"shareholder_id"`
	Amount                  float64   `gorm:"column:amount" json:"amount" form:"amount"`
	State                   int64     `gorm:"column:state" json:"state" form:"state"`
	Content                 string    `gorm:"column:content" json:"content" form:"content"`
	WithLog                 string    `gorm:"column:with_log" json:"with_log" form:"with_log"`
	Hash                    string    `gorm:"column:hash" json:"hash" form:"hash"`
	Mobile                  string    `gorm:"column:mobile" json:"mobile" form:"mobile" `
	Income                  float64   `gorm:"column:income" json:"income" form:"income" `
	WithdrawalLimit         float64   `gorm:"column:withdrawal_limit" json:"withdrawal_limit" form:"withdrawal_limit" `
	RecentWithdrawalAccount string    `gorm:"column:recent_withdrawal_account" json:"recent_withdrawal_account" form:"recent_withdrawal_account" `
	IsEnable                bool      `gorm:"column:is_enable" json:"is_enable" form:"is_enable" `
	CreateAt                time.Time `gorm:"column:create_at" json:"create_at" form:"create_at"`
	EndAt                   time.Time `gorm:"column:end_at" json:"end_at" form:"end_at" `
}

//WithdrawalList ????????????????????????
type WithdrawalList struct {
	PageInfo PageInfo         `json:"page_info"`
	Data     []UserWithdrawal `json:"data"`
}

//TableName ??????
func (Withdrawal) TableName() string {
	return "withdrawal"
}

//Update ??????withdrawal
func (withdrawal *Withdrawal) Update(withdrawalParam *Withdrawal) error {
	if result := global.DB.GetDbR().Table(Withdrawal{}.TableName()).Where("id", withdrawalParam.ID).Save(&withdrawalParam); result.Error != nil {
		// ??????db????????????...
		global.Logger.Error("?????? Withdrawal ??????:", result.Error)
		return result.Error
	}
	return nil
}

//Create ??????withdrawal
func (withdrawal *Withdrawal) Create(withdrawalParam *Withdrawal) error {
	if result := global.DB.GetDbR().Table(Withdrawal{}.TableName()).Create(&withdrawalParam); result.Error != nil {
		// ??????db????????????...
		global.Logger.Error("?????? Withdrawal ??????:", result.Error)
		return result.Error
	}
	return nil
}

//GetByID ??????id????????????
func (withdrawal *Withdrawal) GetByID(ID string) (*Withdrawal, error) {
	withdrawalRes := new(Withdrawal)
	if result := global.DB.GetDbR().Table(Withdrawal{}.TableName()).Where("id = ?", ID).First(&withdrawalRes); result.Error != nil {
		// ??????db????????????...
		global.Logger.Error("model Withdrawal GetByID ??????:", result.Error)
		return nil, result.Error
	}
	return withdrawal, nil
}

//GetList ????????????????????????
func (withdrawal *Withdrawal) GetList(num uint64, page uint64) (*WithdrawalList, error) {
	withdrawalList := new(WithdrawalList)
	withdrawalList.PageInfo.Page = page
	withdrawalList.PageInfo.Size = num
	withdrawalList.Data = make([]UserWithdrawal, 0)

	tx := global.DB.GetDbR().Table(withdrawal.TableName() + " AS w").
		Joins("left join shareholder as s on s.id=w.shareholder_id")

	// ????????????
	if err := tx.Count(&withdrawalList.PageInfo.Total).Error; err != nil {
		global.Logger.Error("?????? ???????????? ?????? ???????????? ??????:", err.Error())
		return nil, err
	}

	// ??????
	tx = tx.Order("w.create_at DESC")

	// ??????
	tx = tx.Limit(int(num)).Offset(int(page*num - num))

	// ????????????
	tx = tx.Select("w.*,s.mobile,s.withdrawal_limit,s.is_enable")

	// ????????????
	if err := tx.Find(&withdrawalList.Data).Error; err != nil {
		global.Logger.Error("?????? ???????????? ?????? ??????:", err.Error())
		return nil, err
	}

	withdrawalList.PageInfo.TotalPage = withdrawalList.PageInfo.Total / int64(withdrawalList.PageInfo.Size)
	if withdrawalList.PageInfo.Total%int64(withdrawalList.PageInfo.Size) > 0 {
		withdrawalList.PageInfo.TotalPage++
	}

	return withdrawalList, nil
}

//GetByShareholderID ??????id????????????
func (withdrawal *Withdrawal) GetByShareholderID(ShareholderID string) ([]Withdrawal, error) {
	withdrawalList := make([]Withdrawal, 0)
	if result := global.DB.GetDbR().Table(Withdrawal{}.TableName()).Where("shareholder_id = ?", ShareholderID).Order("create_at DESC").Find(&withdrawalList); result.Error != nil {
		// ??????db????????????...
		global.Logger.Error("model Withdrawal GetByShareholderID ??????:", result.Error)
		return nil, result.Error
	}
	return withdrawalList, nil
}

//WithdrawalSumByShareholder ??????????????????????????????
func (withdrawal *Withdrawal) WithdrawalSumByShareholder(ShareholderID string) (float64, error) {
	var withdrawalSum float64
	var count int64
	tx := global.DB.GetDbR().Table(Withdrawal{}.TableName())
	tx = tx.Where("shareholder_id", ShareholderID).Where("state", enum.WithdrawalSuccess)

	// ????????????
	if err := tx.Count(&count).Error; err != nil {
		return 0.0, err
	}
	if count != 0 {
		//???????????????
		err := tx.Select("SUM(amount) as income").Pluck("amount", &withdrawalSum).Error
		if err != nil {
			global.Logger.Error("model FilPoolDailyIncome Info ??????:", err.Error())
			return 0.0, err
		}
	} else {
		withdrawalSum = 0.0
	}

	return withdrawalSum, nil
}
