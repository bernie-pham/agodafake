package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/bernie-pham/agodafake/internal/model"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	Errorlog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan model.MailData
}
