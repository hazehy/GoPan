package helper

import (
	"log"
	"strings"

	"gopan/gopan/models"

	"xorm.io/xorm"
)

func AddAuditLog(engine *xorm.Engine, actorIdentity, actorName string, actorRole int, action, targetType, targetIdentity, detail string) {
	if engine == nil || strings.TrimSpace(action) == "" {
		return
	}

	trimmedActorIdentity := strings.TrimSpace(actorIdentity)
	trimmedActorName := strings.TrimSpace(actorName)
	resolvedActorRole := actorRole
	if trimmedActorIdentity != "" && trimmedActorName == "" {
		user := new(models.User)
		has, err := engine.Where("identity = ?", trimmedActorIdentity).Get(user)
		if err == nil && has {
			trimmedActorName = strings.TrimSpace(user.Name)
			if resolvedActorRole <= 0 {
				resolvedActorRole = user.Role
			}
		}
	}

	entry := &models.AuditLog{
		Identity:       GenerateUUID(),
		ActorIdentity:  trimmedActorIdentity,
		ActorName:      trimmedActorName,
		ActorRole:      resolvedActorRole,
		Action:         strings.TrimSpace(action),
		TargetType:     strings.TrimSpace(targetType),
		TargetIdentity: strings.TrimSpace(targetIdentity),
		Detail:         strings.TrimSpace(detail),
	}

	if _, err := engine.Insert(entry); err != nil {
		log.Printf("写入审计日志失败: %v", err)
	}
}
