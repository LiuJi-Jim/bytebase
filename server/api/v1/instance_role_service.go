package v1

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/bytebase/bytebase/api"
	"github.com/bytebase/bytebase/common"
	"github.com/bytebase/bytebase/plugin/db"
	v1pb "github.com/bytebase/bytebase/proto/generated-go/v1"
	"github.com/bytebase/bytebase/server/component/dbfactory"
	"github.com/bytebase/bytebase/store"
)

// InstanceRoleService implements the database role service.
type InstanceRoleService struct {
	v1pb.UnimplementedInstanceRoleServiceServer
	store     *store.Store
	dbFactory *dbfactory.DBFactory
}

// NewInstanceRoleService creates a new InstanceRoleService.
func NewInstanceRoleService(store *store.Store, dbFactory *dbfactory.DBFactory) *InstanceRoleService {
	return &InstanceRoleService{
		store:     store,
		dbFactory: dbFactory,
	}
}

// GetRole gets an role.
func (s *InstanceRoleService) GetRole(ctx context.Context, request *v1pb.GetRoleRequest) (*v1pb.DatabaseRole, error) {
	environmentID, instanceID, roleName, err := getEnvironmentInstanceRoleID(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	environment, instance, err := s.getEnvironmentAndInstance(ctx, environmentID, instanceID)
	if err != nil {
		return nil, err
	}

	role, err := func() (*db.DatabaseRoleMessage, error) {
		driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, convertToLegacyInstance(instance, environment), "" /* database name */)
		if err != nil {
			return nil, err
		}
		defer driver.Close(ctx)

		role, err := driver.FindRole(ctx, roleName)
		if err != nil {
			return nil, err
		}

		return role, nil
	}()
	if err != nil {
		if common.ErrorCode(err) == common.NotFound {
			return nil, status.Errorf(codes.NotFound, "role %s not found", request.Name)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return convertToRole(role, instance), nil
}

// ListRoles lists all roles in an instance.
func (s *InstanceRoleService) ListRoles(ctx context.Context, request *v1pb.ListRolesRequest) (*v1pb.ListRolesResponse, error) {
	environmentID, instanceID, err := getEnvironmentInstanceID(request.Parent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	environment, instance, err := s.getEnvironmentAndInstance(ctx, environmentID, instanceID)
	if err != nil {
		return nil, err
	}

	roleList, err := func() ([]*db.DatabaseRoleMessage, error) {
		driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, convertToLegacyInstance(instance, environment), "" /* database name */)
		if err != nil {
			return nil, err
		}
		defer driver.Close(ctx)

		roleList, err := driver.ListRole(ctx)
		if err != nil {
			return nil, err
		}

		return roleList, nil
	}()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	response := &v1pb.ListRolesResponse{}
	for _, role := range roleList {
		response.Roles = append(response.Roles, convertToRole(role, instance))
	}
	return response, nil
}

// CreateRole creates an role.
func (s *InstanceRoleService) CreateRole(ctx context.Context, request *v1pb.CreateRoleRequest) (*v1pb.DatabaseRole, error) {
	if request.Role == nil {
		return nil, status.Errorf(codes.InvalidArgument, "role must be set")
	}
	environmentID, instanceID, err := getEnvironmentInstanceID(request.Parent)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	environment, instance, err := s.getEnvironmentAndInstance(ctx, environmentID, instanceID)
	if err != nil {
		return nil, err
	}

	roleUpsert := &db.DatabaseRoleUpsertMessage{
		Name:            request.Role.Title,
		Password:        request.Role.Password,
		ConnectionLimit: request.Role.ConnectionLimit,
		ValidUntil:      request.Role.ValidUntil,
		Attribute: &db.DatabaseRoleAttributeMessage{
			SuperUser:   request.Role.Attribute.SuperUser,
			NoInherit:   request.Role.Attribute.NoInherit,
			CreateRole:  request.Role.Attribute.CreateRole,
			CreateDb:    request.Role.Attribute.CreateDb,
			CanLogin:    request.Role.Attribute.CanLogin,
			Replication: request.Role.Attribute.Replication,
			BypassRls:   request.Role.Attribute.BypassRls,
		},
	}
	if err := validateRole(roleUpsert); err != nil {
		return nil, err
	}

	role, err := func() (*db.DatabaseRoleMessage, error) {
		driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, convertToLegacyInstance(instance, environment), "" /* database name */)
		if err != nil {
			return nil, err
		}
		defer driver.Close(ctx)

		role, err := driver.CreateRole(ctx, roleUpsert)
		if err != nil {
			return nil, err
		}

		return role, nil
	}()
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return convertToRole(role, instance), nil
}

// UpdateRole updates an role.
func (s *InstanceRoleService) UpdateRole(ctx context.Context, request *v1pb.UpdateRoleRequest) (*v1pb.DatabaseRole, error) {
	if request.UpdateMask == nil {
		return nil, status.Errorf(codes.InvalidArgument, "update_mask must be set")
	}
	environmentID, instanceID, roleName, err := getEnvironmentInstanceRoleID(request.Role.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	environment, instance, err := s.getEnvironmentAndInstance(ctx, environmentID, instanceID)
	if err != nil {
		return nil, err
	}

	upsert := &db.DatabaseRoleUpsertMessage{
		Name: roleName,
	}
	for _, path := range request.UpdateMask.Paths {
		switch path {
		case "role.title":
			upsert.Name = request.Role.Title
		case "role.password":
			upsert.Password = request.Role.Password
		case "role.connection_limit":
			upsert.ConnectionLimit = request.Role.ConnectionLimit
		case "role.valid_until":
			upsert.ValidUntil = request.Role.ValidUntil
		case "role.attribute":
			upsert.Attribute = &db.DatabaseRoleAttributeMessage{
				SuperUser:   request.Role.Attribute.SuperUser,
				NoInherit:   request.Role.Attribute.NoInherit,
				CreateRole:  request.Role.Attribute.CreateRole,
				CreateDb:    request.Role.Attribute.CreateDb,
				CanLogin:    request.Role.Attribute.CanLogin,
				Replication: request.Role.Attribute.Replication,
				BypassRls:   request.Role.Attribute.BypassRls,
			}
		}
	}
	if err := validateRole(upsert); err != nil {
		return nil, err
	}

	role, err := func() (*db.DatabaseRoleMessage, error) {
		driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, convertToLegacyInstance(instance, environment), "" /* database name */)
		if err != nil {
			return nil, err
		}
		defer driver.Close(ctx)

		role, err := driver.UpdateRole(ctx, roleName, upsert)
		if err != nil {
			return nil, err
		}

		return role, nil
	}()
	if err != nil {
		if common.ErrorCode(err) == common.NotFound {
			return nil, status.Errorf(codes.NotFound, "role %s not found", request.Role.Name)
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return convertToRole(role, instance), nil
}

// DeleteRole deletes an role.
func (s *InstanceRoleService) DeleteRole(ctx context.Context, request *v1pb.DeleteRoleRequest) (*emptypb.Empty, error) {
	environmentID, instanceID, roleName, err := getEnvironmentInstanceRoleID(request.Name)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	environment, instance, err := s.getEnvironmentAndInstance(ctx, environmentID, instanceID)
	if err != nil {
		return nil, err
	}

	if err := func() error {
		driver, err := s.dbFactory.GetAdminDatabaseDriver(ctx, convertToLegacyInstance(instance, environment), "" /* database name */)
		if err != nil {
			return err
		}
		defer driver.Close(ctx)

		return driver.DeleteRole(ctx, roleName)
	}(); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

// UndeleteRole undeletes an role.
func (*InstanceRoleService) UndeleteRole(_ context.Context, _ *v1pb.UndeleteRoleRequest) (*v1pb.DatabaseRole, error) {
	return nil, status.Errorf(codes.Unimplemented, "Undelete role is not supported")
}

func (s *InstanceRoleService) getEnvironmentAndInstance(ctx context.Context, environmentID, instanceID string) (*store.EnvironmentMessage, *store.InstanceMessage, error) {
	environment, err := s.store.GetEnvironmentV2(ctx, &store.FindEnvironmentMessage{
		ResourceID: &environmentID,
	})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, err.Error())
	}
	if environment == nil {
		return nil, nil, status.Errorf(codes.NotFound, "environment %q not found", environmentID)
	}
	// We don't allow access even for the read API because the API will call user instances without using Bytebase metadata.
	if environment.Deleted {
		return nil, nil, status.Errorf(codes.InvalidArgument, "environment %q has been deleted", environmentID)
	}

	instance, err := s.store.GetInstanceV2(ctx, &store.FindInstanceMessage{
		EnvironmentID: &environmentID,
		ResourceID:    &instanceID,
	})
	if err != nil {
		return nil, nil, status.Errorf(codes.Internal, err.Error())
	}
	if instance == nil {
		return nil, nil, status.Errorf(codes.NotFound, "instance %q not found", instanceID)
	}
	if instance.Deleted {
		return nil, nil, status.Errorf(codes.InvalidArgument, "instance %q has been deleted", instanceID)
	}

	return environment, instance, nil
}

func convertToLegacyInstance(instance *store.InstanceMessage, environment *store.EnvironmentMessage) *api.Instance {
	port := ""
	host := ""
	ds := []*api.DataSource{}
	for _, dataSource := range instance.DataSources {
		if dataSource.Type == api.Admin {
			port = dataSource.Port
			host = dataSource.Host
		}
		ds = append(ds, &api.DataSource{
			Type:     dataSource.Type,
			Name:     dataSource.Title,
			Username: dataSource.Username,
			Password: dataSource.Password,
			SslCa:    dataSource.SslCa,
			SslCert:  dataSource.SslCert,
			SslKey:   dataSource.SslKey,
			Host:     dataSource.Host,
			Port:     dataSource.Port,
			Database: dataSource.Database,
		})
	}

	return &api.Instance{
		ID:     instance.UID,
		Engine: instance.Engine,
		Host:   host,
		Port:   port,
		Environment: &api.Environment{
			Name: environment.Title,
		},
		DataSourceList: ds,
	}
}

func convertToRole(role *db.DatabaseRoleMessage, instance *store.InstanceMessage) *v1pb.DatabaseRole {
	return &v1pb.DatabaseRole{
		Name:            fmt.Sprintf("environments/%s/instances/%s/roles/%s", instance.EnvironmentID, instance.ResourceID, role.Name),
		Title:           role.Name,
		ConnectionLimit: &role.ConnectionLimit,
		ValidUntil:      role.ValidUntil,
		Attribute: &v1pb.DatabaseRoleAttribute{
			SuperUser:   role.Attribute.SuperUser,
			NoInherit:   role.Attribute.NoInherit,
			CreateRole:  role.Attribute.CreateRole,
			CreateDb:    role.Attribute.CreateDb,
			CanLogin:    role.Attribute.CanLogin,
			Replication: role.Attribute.Replication,
			BypassRls:   role.Attribute.BypassRls,
		},
	}
}

func validateRole(upsert *db.DatabaseRoleUpsertMessage) error {
	if upsert.Name == "" {
		return status.Errorf(codes.InvalidArgument, "Invalid role name, role name cannot be empty")
	}
	if v := upsert.ConnectionLimit; v != nil && *v < int32(-1) {
		return status.Errorf(codes.InvalidArgument, "Invalid connection limit, it should greater than or equal to -1")
	}
	if v := upsert.ValidUntil; v != nil {
		if _, err := time.Parse(time.RFC3339, *v); err != nil {
			return status.Errorf(codes.InvalidArgument, "Invalid timestamp for valid_until, timestamp should in '2006-01-02T15:04:05+08:00' format.")
		}
	}

	return nil
}
