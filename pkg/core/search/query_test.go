package search

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type ApplicationQuery struct {
	Id       string    `search:"type:icontains;column:id;table:receipt" form:"id"`
	Domain   string    `search:"type:icontains;column:domain;table:receipt" form:"domain"`
	Version  string    `search:"type:exact;column:version;table:receipt" form:"version"`
	Status   []int     `search:"type:in;column:status;table:receipt" form:"status"`
	Start    time.Time `search:"type:gte;column:created_at;table:receipt" form:"start"`
	End      time.Time `search:"type:lte;column:created_at;table:receipt" form:"end"`
	TestJoin `search:"type:left;on:id:receipt_id;table:receipt_goods;join:receipts"`
	NotNeed  string `search:"-"`
	ApplicationOrder
}

type ApplicationOrder struct {
	IdOrder string `search:"type:order;column:id;table:receipt" form:"id_order"`
}

type TestJoin struct {
	PaymentAccount string `search:"type:icontains;column:payment_account;table:receipts" form:"payment_account"`
}

func TestResolveSearchQuery(t *testing.T) {
	// Only pass t into top-level Convey calls
	Convey("Given some integer with a starting value", t, func() {
		d := ApplicationQuery{
			Id:               "aaa",
			Domain:           "bbb",
			Version:          "ccc",
			Status:           []int{1, 2},
			Start:            time.Now().Add(-8 * time.Hour),
			End:              time.Now(),
			ApplicationOrder: ApplicationOrder{IdOrder: "desc"},
			TestJoin:         TestJoin{PaymentAccount: "1212"},
		}
		condition := &GormCondition{
			GormPublic: GormPublic{},
			Join:       make([]*GormJoin, 0),
		}
		ResolveSearchQuery("mysql", d, condition)
		fmt.Println(condition)
	})
}
