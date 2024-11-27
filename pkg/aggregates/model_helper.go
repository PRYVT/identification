package aggregates

import (
	models "github.com/L4B0MB4/PRYVT/identification/pkg/models/query"
)

func GetUserModelFromAggregate(userAggregate *UserAggregate) *models.UserInfo {
	return &models.UserInfo{
		ID:           userAggregate.AggregateId,
		DisplayName:  userAggregate.DisplayName,
		Name:         userAggregate.Name,
		Email:        userAggregate.Email,
		ChangeDate:   userAggregate.ChangeDate,
		PasswordHash: userAggregate.PasswordHash,
	}
}
