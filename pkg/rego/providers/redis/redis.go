package redis

import (
	"fmt"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/open-policy-agent/opa/util"
	rg "github.com/redislabs/redisgraph-go"
)

type Config struct {
	Address  string `json:"address,omitempty"`
	Password string `json:"password,omitempty"`
}

type RedisProvider struct {
	config Config
}

func New(config []byte) (*RedisProvider, error) {
	c := Config{
		Address: ":6379",
	}
	if err := util.Unmarshal(config, &c); err != nil {
		return nil, err
	}
	return &RedisProvider{
		config: c,
	}, nil
}

func (p *RedisProvider) Can(user, action, resource, owner string) (bool, error) {
	fmt.Println("gbac.can", user, "perform", action, "on resource", resource, "in", owner)

	conn, _ := redis.Dial("tcp", p.config.Address)
	defer conn.Close()

	graph := rg.GraphNew("gbac", conn)

	actionsRes, err := graph.Query(
		fmt.Sprintf(`MATCH (action{name:'%v'})<-[:can*0..]-(a:action) RETURN a.name`, action),
	)

	if err != nil {
		return false, err
	}

	actions := []string{}
	for actionsRes.Next() {
		record := actionsRes.Record()
		name, hasName := record.Get("a.name")
		if !hasName {
			return false, fmt.Errorf("Issue getting valid actions. Lookup: %v. Action record has no name: %v", action, record)
		}
		actions = append(actions, `'`+name.(string)+`'`)
	}
	validActions := fmt.Sprintf("[%v]", strings.Join(actions, ","))
	fmt.Println("Actions for action", actions, action, validActions)

	resParts := strings.Split(resource, ":")

	if len(resParts) > 1 {
		directPermQuery := fmt.Sprintf(
			`MATCH (u:user{id:'%v'})-[a:%v]->(%v{id:'%v'}) WHERE a.action IN %v RETURN u.id`,
			user,
			resParts[0],
			resParts[0],
			resParts[1],
			validActions,
		)
		hasDirectPerm, err := graph.Query(
			directPermQuery,
		)

		if err != nil {
			return false, fmt.Errorf("Error while checking direct permissions: %v", err)
		}

		if hasDirectPerm != nil && !hasDirectPerm.Empty() {
			return true, nil
		}
	}

	generalPermQuery := fmt.Sprintf(
		`MATCH (c:container{id:'%v'})-[:parent*0..]->(container)<-[a:all]-(user{id:'%v'}) WHERE a.action in %v RETURN c.id`,
		owner,
		user,
		validActions,
	)
	hasGeneralPerm, err := graph.Query(generalPermQuery)

	if err != nil {
		return false, err
	}

	if hasGeneralPerm != nil && !hasGeneralPerm.Empty() {
		return true, nil
	}

	resourcePermQuery := fmt.Sprintf(
		`MATCH (c:container{id:'%v'})-[:parent*0..]->(container)<-[a:%v]-(user{id:'%v'}) WHERE a.action in %v RETURN c.id`,
		owner,
		resParts[0],
		user,
		validActions,
	)
	hasResourcePerm, err := graph.Query(resourcePermQuery)

	if err != nil {
		return false, err
	}

	if hasResourcePerm != nil && !hasResourcePerm.Empty() {
		return true, nil
	}

	return false, nil
}
