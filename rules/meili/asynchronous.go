package meili

import (
	"github.com/meilisearch/meilisearch-go"
	"main/common/event/rule"
)

type ClientService struct {
	Client        *meilisearch.ServiceManager
	Rule          *rule.RuleServer
	ChannelClient *rule.ChannelClient
}

func (cs *ClientService) PutNewClient() {}
