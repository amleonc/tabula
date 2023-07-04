package internal

import (
	"github.com/amleonc/tabula/internal/services/comment"
	"github.com/amleonc/tabula/internal/services/thread"
	"github.com/amleonc/tabula/internal/services/topic"
	"github.com/amleonc/tabula/internal/services/user"
)

var (
	CommentService = comment.GetService()
	ThreadService  = thread.GetService()
	TopicService   = topic.GetService()
	UserService    = user.NewService()
)
