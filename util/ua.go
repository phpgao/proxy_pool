package util

import (
	"github.com/corpix/uarand"
)

func GetRandomUA() string {
	return uarand.GetRandom()
}
