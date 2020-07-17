package main

import (
	"github.com/Knetic/govaluate"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/graphcms/autobahn/casbin/testing/sql"
	"log"
	"strings"
)


func main() {
	e, err := casbin.NewEnforcer("./model.conf", "./policy.csv")
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	ok, err := e.Enforce("test", "/post/jk", struct {
		Owner string
		Age int
		Foobar int
		Foobaz int
	}{
		//Owner: "test",
		//Age: 20,
	}, "read")

	if err != nil {
		// handle err
		log.Fatal(err)
	}
	log.Println(ok)

	rvals := []interface{} {
		"test",
		"/post/test",
		struct {
			Owner string
			Age int
			Foobar int
			Foobaz int
		}{
			//Owner: "test",
			//Age: 20,
		},
		"read",
	}

	ok, err = e.EnforceWithMatcher("g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act", rvals...)

	if err != nil {
		// handle err
		log.Fatal(err)
	}
	log.Println(ok)

	permissions, err := e.GetImplicitPermissionsForUser(rvals[0].(string))
	if err != nil {
		// handle err
		log.Fatal(err)
	}
	log.Println(permissions)
	where := []string{}
	for _, p := range permissions {
		if util.KeyMatch(rvals[1].(string), p[1]) {
			log.Println(p, p[2])
			exp, err := govaluate.NewEvaluableExpression(util.RemoveComments(util.EscapeAssertion(p[2])))
			if err != nil {
				// handle err
				log.Fatal(err)
			}
			expTokens := exp.Tokens()
			for _, token := range expTokens {
				log.Printf("%s %s %T", token.Kind, token.Value,  token.Value)
			}
			q, err := sql.EvaluableExpression{
				EvaluableExpression: exp,
				Sub: rvals[0].(string),
			}.ToSQLQuery()
			if err != nil {
				// handle err
				log.Fatal(err)
			}
			where = append(where, q)
		}
	}

	log.Println(strings.Join(where, " AND "))

}
