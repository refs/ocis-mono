package service

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"github.com/CiscoM31/godata"
	"github.com/blevesearch/bleve"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes/empty"
	merrors "github.com/micro/go-micro/v2/errors"
	"github.com/refs/ocis-mono/ocis-accounts/pkg/proto/v0"
	"github.com/refs/ocis-mono/ocis-accounts/pkg/provider"
)

func (s Service) indexGroups(path string) (err error) {
	var f *os.File
	if f, err = os.Open(path); err != nil {
		s.log.Error().Err(err).Str("dir", path).Msg("could not open groups folder")
		return
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		s.log.Error().Err(err).Str("dir", path).Msg("could not list groups folder")
		return
	}
	for _, file := range list {
		err = s.indexGroup(file.Name())
		if err != nil {
			s.log.Error().Err(err).Str("file", file.Name()).Msg("could not index account")
		}
	}

	return
}

// accLock mutually exclude readers from writers on group files
var groupLock sync.Mutex

func (s Service) indexGroup(id string) error {
	g := &proto.BleveGroup{
		BleveType: "group",
	}
	if err := s.loadGroup(id, &g.Group); err != nil {
		s.log.Error().Err(err).Str("group", id).Msg("could not load group")
		return err
	}
	s.log.Debug().Interface("group", g).Msg("found group")
	if err := s.index.Index(g.Id, g); err != nil {
		s.log.Error().Err(err).Interface("group", g).Msg("could not index group")
		return err
	}
	return nil
}

func (s Service) loadGroup(id string, g *proto.Group) (err error) {
	path := filepath.Join(s.Config.Server.AccountsDataPath, "groups", id)

	groupLock.Lock()
	defer groupLock.Unlock()
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		return merrors.NotFound(s.id, "could not read group: %v", err.Error())
	}

	if err = json.Unmarshal(data, g); err != nil {
		return merrors.InternalServerError(s.id, "could not unmarshal group: %v", err.Error())
	}

	return
}

func (s Service) writeGroup(g *proto.Group) (err error) {

	// leave only the member id
	s.deflateMembers(g)

	var bytes []byte
	if bytes, err = json.Marshal(g); err != nil {
		return merrors.InternalServerError(s.id, "could not marshal group: %v", err.Error())
	}

	path := filepath.Join(s.Config.Server.AccountsDataPath, "groups", g.Id)

	groupLock.Lock()
	defer groupLock.Unlock()
	if err = ioutil.WriteFile(path, bytes, 0600); err != nil {
		return merrors.InternalServerError(s.id, "could not write group: %v", err.Error())
	}
	return
}

func (s Service) expandMembers(g *proto.Group) {
	if g == nil {
		return
	}
	expanded := []*proto.Account{}
	for i := range g.Members {
		// TODO resolve by name, when a create or update is issued they may not have an id? fall back to searching the group id in the index?
		a := &proto.Account{}
		if err := s.loadAccount(g.Members[i].Id, a); err == nil {
			expanded = append(expanded, a)
		} else {
			// log errors but continue execution for now
			s.log.Error().Err(err).Str("id", g.Members[i].Id).Msg("could not load account")
		}
	}
	g.Members = expanded
}

// deflateMembers replaces the users of a group with an instance that only contains the id
func (s Service) deflateMembers(g *proto.Group) {
	if g == nil {
		return
	}
	deflated := []*proto.Account{}
	for i := range g.Members {
		if g.Members[i].Id != "" {
			deflated = append(deflated, &proto.Account{Id: g.Members[i].Id})
		} else {
			// TODO fetch and use an id when group only has a name but no id
			s.log.Error().Str("id", g.Id).Interface("account", g.Members[i]).Msg("resolving members by name is not implemented yet")
		}
	}
	g.Members = deflated
}

// ListGroups implements the GroupsServiceHandler interface
func (s Service) ListGroups(c context.Context, in *proto.ListGroupsRequest, out *proto.ListGroupsResponse) (err error) {

	// only search for groups
	tq := bleve.NewTermQuery("group")
	tq.SetField("bleve_type")

	query := bleve.NewConjunctionQuery(tq)

	if in.Query != "" {
		// parse the query like an odata filter
		var q *godata.GoDataFilterQuery
		if q, err = godata.ParseFilterString(in.Query); err != nil {
			s.log.Error().Err(err).Msg("could not parse query")
			return merrors.InternalServerError(s.id, "could not parse query: %v", err.Error())
		}

		// convert to bleve query
		bq, err := provider.BuildBleveQuery(q)
		if err != nil {
			s.log.Error().Err(err).Msg("could not build bleve query")
			return merrors.InternalServerError(s.id, "could not build bleve query: %v", err.Error())
		}
		query.AddQuery(bq)
	}

	s.log.Debug().Interface("query", query).Msg("using query")

	searchRequest := bleve.NewSearchRequest(query)
	var searchResult *bleve.SearchResult
	searchResult, err = s.index.Search(searchRequest)
	if err != nil {
		s.log.Error().Err(err).Msg("could not execute bleve search")
		return merrors.InternalServerError(s.id, "could not execute bleve search: %v", err.Error())
	}

	s.log.Debug().Interface("result", searchResult).Msg("result")

	out.Groups = make([]*proto.Group, 0)

	for _, hit := range searchResult.Hits {

		g := &proto.Group{}
		if err = s.loadGroup(hit.ID, g); err != nil {
			s.log.Error().Err(err).Str("group", hit.ID).Msg("could not load group, skipping")
			continue
		}
		s.log.Debug().Interface("group", g).Msg("found group")

		// TODO add accounts if requested
		// if in.FieldMask ...
		s.expandMembers(g)

		out.Groups = append(out.Groups, g)
	}

	return
}

// GetGroup implements the GroupsServiceHandler interface
func (s Service) GetGroup(c context.Context, in *proto.GetGroupRequest, out *proto.Group) (err error) {
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	if err = s.loadGroup(id, out); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not load group")
		return
	}
	s.log.Debug().Interface("group", out).Msg("found group")

	// TODO only add accounts if requested
	// if in.FieldMask ...
	s.expandMembers(out)

	return
}

// CreateGroup implements the GroupsServiceHandler interface
func (s Service) CreateGroup(c context.Context, in *proto.CreateGroupRequest, out *proto.Group) (err error) {
	var id string
	if in.Group == nil {
		return merrors.BadRequest(s.id, "account missing")
	}
	if in.Group.Id == "" {
		in.Group.Id = uuid.Must(uuid.NewV4()).String()
	}

	if id, err = cleanupID(in.Group.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	// extract member id
	s.deflateMembers(in.Group)

	if err = s.writeGroup(in.Group); err != nil {
		s.log.Error().Err(err).Interface("group", in.Group).Msg("could not persist new group")
		return
	}

	if err = s.indexGroup(id); err != nil {
		return merrors.InternalServerError(s.id, "could not index new group: %v", err.Error())
	}

	return
}

// UpdateGroup implements the GroupsServiceHandler interface
func (s Service) UpdateGroup(c context.Context, in *proto.UpdateGroupRequest, out *proto.Group) (err error) {
	return merrors.InternalServerError(s.id, "not implemented")
}

// DeleteGroup implements the GroupsServiceHandler interface
func (s Service) DeleteGroup(c context.Context, in *proto.DeleteGroupRequest, out *empty.Empty) (err error) {
	var id string
	if id, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}
	path := filepath.Join(s.Config.Server.AccountsDataPath, "groups", id)

	g := &proto.Group{}
	if err = s.loadGroup(id, g); err != nil {
		s.log.Error().Err(err).Str("id", id).Msg("could not load account")
		return
	}

	// delete memberof relationship in users
	for i := range g.Members {
		err = s.RemoveMember(c, &proto.RemoveMemberRequest{
			AccountId: g.Members[i].Id,
			GroupId:   id,
		}, g)
		if err != nil {
			s.log.Error().Err(err).Str("groupid", id).Str("accountid", g.Members[i].Id).Msg("could not remove account memberof, skipping")
		}
	}
	if err = os.Remove(path); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove group")
		return merrors.InternalServerError(s.id, "could not remove group: %v", err.Error())
	}

	if err = s.index.Delete(id); err != nil {
		s.log.Error().Err(err).Str("id", id).Str("path", path).Msg("could not remove group from index")
		return merrors.InternalServerError(s.id, "could not remove group from index: %v", err.Error())
	}

	s.log.Info().Str("id", id).Msg("deleted group")
	return
}

// AddMember implements the GroupsServiceHandler interface
func (s Service) AddMember(c context.Context, in *proto.AddMemberRequest, out *proto.Group) (err error) {

	// cleanup ids
	var groupID string
	if groupID, err = cleanupID(in.GroupId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	var accountID string
	if accountID, err = cleanupID(in.AccountId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	// load structs
	a := &proto.Account{}
	if err = s.loadAccount(accountID, a); err != nil {
		s.log.Error().Err(err).Str("id", accountID).Msg("could not load account")
		return
	}

	g := &proto.Group{}
	if err = s.loadGroup(groupID, g); err != nil {
		s.log.Error().Err(err).Str("id", groupID).Msg("could not load group")
		return
	}

	// check if we need to add the account to the group
	alreadyRelated := false
	for i := range g.Members {
		if g.Members[i].Id == a.Id {
			alreadyRelated = true
		}
	}
	if !alreadyRelated {
		g.Members = append(g.Members, a)
	}

	// check if we need to add the group to the account
	alreadyRelated = false
	for i := range a.MemberOf {
		if a.MemberOf[i].Id == g.Id {
			alreadyRelated = true
			break
		}
	}
	if !alreadyRelated {
		a.MemberOf = append(a.MemberOf, g)
	}

	if err = s.writeAccount(a); err != nil {
		s.log.Error().Err(err).Interface("account", a).Msg("could not persist account")
		return
	}
	if err = s.writeGroup(g); err != nil {
		s.log.Error().Err(err).Interface("group", g).Msg("could not persist group")
		return
	}
	// FIXME update index!
	// TODO rollback changes when only one of them failed?
	// TODO store relation in another file?
	// TODO return error if they are already related?
	return nil
}

// RemoveMember implements the GroupsServiceHandler interface
func (s Service) RemoveMember(c context.Context, in *proto.RemoveMemberRequest, out *proto.Group) (err error) {

	// cleanup ids
	var groupID string
	if groupID, err = cleanupID(in.GroupId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	var accountID string
	if accountID, err = cleanupID(in.AccountId); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up account id: %v", err.Error())
	}

	// load structs
	a := &proto.Account{}
	if err = s.loadAccount(accountID, a); err != nil {
		s.log.Error().Err(err).Str("id", accountID).Msg("could not load account")
		return
	}

	g := &proto.Group{}
	if err = s.loadGroup(groupID, g); err != nil {
		s.log.Error().Err(err).Str("id", groupID).Msg("could not load group")
		return
	}

	//remove the account from the group if it exists
	newMembers := []*proto.Account{}
	for i := range g.Members {
		if g.Members[i].Id != a.Id {
			newMembers = append(newMembers, g.Members[i])
		}
	}
	g.Members = newMembers

	// remove the group from the account if it exists
	newGroups := []*proto.Group{}
	for i := range a.MemberOf {
		if a.MemberOf[i].Id != g.Id {
			newGroups = append(newGroups, a.MemberOf[i])
		}
	}
	a.MemberOf = newGroups

	if err = s.writeAccount(a); err != nil {
		s.log.Error().Err(err).Interface("account", a).Msg("could not persist account")
		return
	}
	if err = s.writeGroup(g); err != nil {
		s.log.Error().Err(err).Interface("group", g).Msg("could not persist group")
		return
	}
	// FIXME update index!
	// TODO rollback changes when only one of them failed?
	// TODO store relation in another file?
	// TODO return error if they are not related?
	return nil
}

// ListMembers implements the GroupsServiceHandler interface
func (s Service) ListMembers(c context.Context, in *proto.ListMembersRequest, out *proto.ListMembersResponse) (err error) {

	// cleanup ids
	var groupID string
	if groupID, err = cleanupID(in.Id); err != nil {
		return merrors.InternalServerError(s.id, "could not clean up group id: %v", err.Error())
	}

	g := &proto.Group{}
	if err = s.loadGroup(groupID, g); err != nil {
		s.log.Error().Err(err).Str("id", groupID).Msg("could not load group")
		return
	}

	// TODO only expand accounts if requested
	// if in.FieldMask ...
	s.expandMembers(g)

	out.Members = g.Members

	return
}
