package workout

import "github.com/gin-gonic/gin"

type Service interface {
}

type service struct {
}

func (s *service) CreateWorkout(c *gin.Context) {}
