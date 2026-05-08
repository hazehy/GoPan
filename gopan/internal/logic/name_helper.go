package logic

import (
	"fmt"

	"gopan/gopan/helper"
	"gopan/gopan/models"

	"xorm.io/xorm"
)

func buildAvailableName(engine *xorm.Engine, userIdentity string, parentId int64, rawName string) (string, error) {
	baseName := helper.NormalizeInput(rawName)
	if baseName == "" {
		baseName = "未命名文件"
	}

	count, err := engine.Table("user_repository").Unscoped().
		Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, parentId, baseName).
		Count(new(models.UserRepository))
	if err != nil {
		return "", err
	}
	if count == 0 {
		return baseName, nil
	}

	for index := 1; ; index++ {
		candidate := fmt.Sprintf("%s(%d)", baseName, index)
		candidateCount, candidateErr := engine.Table("user_repository").Unscoped().
			Where("user_identity = ? AND parent_id = ? AND name = ?", userIdentity, parentId, candidate).
			Count(new(models.UserRepository))
		if candidateErr != nil {
			return "", candidateErr
		}
		if candidateCount == 0 {
			return candidate, nil
		}
	}
}
