package opksshplugingoogleworkspace

import (
	"maps"
	"time"
)

type (
	Group struct {
		FetchedAt time.Time          `json:"fetched_at"` // fetcher at
		Email     string             `json:"email"`      // group's email
		Members   map[string]*Member `json:"members"`    // user's email => member
	}

	Customer struct {
		CustomerID string            `json:"customer_id"`
		Groups     map[string]*Group `json:"groups"`
	}

	Info struct {
		Customers map[string]*Customer `json:"customers"`
	}
)

func (g *Group) IsExpired(deadline time.Time) bool {
	if g == nil {
		return true
	}
	return g.FetchedAt.Before(deadline)
}

func (g *Group) GetMember(deadline time.Time, userEmail string) *Member {
	if g == nil {
		return nil
	}
	if g.FetchedAt.Before(deadline) {
		return nil
	}
	return g.Members[userEmail]
}

func (g *Group) AddMember(member *Member) {
	if g == nil {
		panic(nil)
	}
	if g.Members == nil {
		g.Members = make(map[string]*Member)
	}
	g.Members[member.Email] = &Member{
		Id:     member.Id,
		Email:  member.Email,
		Status: member.Status,
		Type:   member.Type,
	}
}

func (c *Customer) GetGroup(deadline time.Time, groupEmail string) *Group {
	if c == nil {
		return nil
	}
	group := c.Groups[groupEmail]
	if group == nil {
		return nil
	}
	if group.FetchedAt.Before(deadline) {
		return nil
	}
	return group
}

func (c *Customer) AddGroup(fetchTime time.Time, groupEmail string) *Group {
	if c == nil {
		panic(nil)
	}
	if c.Groups == nil {
		c.Groups = make(map[string]*Group)
	}
	group := c.GetGroup(fetchTime, groupEmail)
	if group != nil && group.FetchedAt.After(fetchTime) {
		// avoid overwrite more fresh data
		return &Group{
			FetchedAt: group.FetchedAt,
			Email:     group.Email,
			Members:   nil,
		}
	}
	if group == nil {
		group = &Group{
			FetchedAt: fetchTime,
			Email:     groupEmail,
		}
		c.Groups[groupEmail] = group
	}
	return group
}

func (i *Info) GetCustomer(customerId string) *Customer {
	if i == nil {
		return nil
	}
	return i.Customers[customerId]
}

func (i *Info) AddCustomer(customerId string) *Customer {
	if i == nil {
		panic(nil)
	}
	if i.Customers == nil {
		i.Customers = make(map[string]*Customer)
	}
	customer := i.GetCustomer(customerId)
	if customer == nil {
		customer = &Customer{
			CustomerID: customerId,
		}
		i.Customers[customerId] = customer
	}
	return customer
}

func (left *Info) Merge(right *Info) {
	if left == nil {
		panic(nil)
	}
	if right == nil {
		return
	}
	maps.Keys(right.Customers)(func(customerId string) bool {
		customer := right.Customers[customerId]
		maps.Keys(customer.Groups)(func(groupEmail string) bool {
			group := customer.Groups[groupEmail]
			maps.Keys(group.Members)(func(memberEmail string) bool {
				member := group.Members[memberEmail]
				left.AddCustomer(customerId).AddGroup(group.FetchedAt, group.Email).AddMember(member)
				return true
			})
			return true
		})
		return true
	})
}
