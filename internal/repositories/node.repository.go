package repositories

import (
	"context"
	"fmt"

	"github.com/dafaath/iot-server/v2/internal/entities"
	"github.com/dafaath/iot-server/v2/internal/helper"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
)

type NodeRepository struct{}

func NewNodeRepository() (NodeRepository, error) {
	return NodeRepository{}, nil
}

func (u *NodeRepository) nodeFieldWithoutId() string {
	return "name, location, id_hardware_node,id_user,is_public,id_hardware_sensor,field_sensor"
}

func (u *NodeRepository) nodeField() string {
	return "id_node, " + u.nodeFieldWithoutId()
}

func (u *NodeRepository) nodePointer(node *entities.Node) []interface{} {
	return []interface{}{&node.IdNode, &node.Name, &node.Location, &node.IdHardwareNode, &node.IdUser, &node.IsPublic, &node.IdHardwareSensor, &node.FieldSensor}
}

func (h *NodeRepository) Create(ctx context.Context, tx helper.Querier, payload *entities.NodeCreate, currentUser *entities.UserRead) (id int, err error) {
	sqlStatement := fmt.Sprintf(`INSERT into NODE 
	(%s) 
	VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id_node`, h.nodeFieldWithoutId())
	err = tx.QueryRow(ctx, sqlStatement,
		payload.Name,
		payload.Location,
		payload.IdHardwareNode,
		currentUser.IdUser,
		payload.IsPublic,
		payload.IdHardwareSensor,
		payload.FieldSensor,
	).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

func (u *NodeRepository) GetAll(ctx context.Context, tx helper.Querier, currentUser *entities.UserRead) (nodes []entities.Node, err error) {
	nodes = []entities.Node{}
	var sqlStatement string
	var rows pgx.Rows
	if currentUser.IsAdmin {
		sqlStatement = fmt.Sprintf(`SELECT %s FROM "node"`, u.nodeField())
		rows, err = tx.Query(ctx, sqlStatement)
		if err != nil {
			return nodes, err
		}
		defer rows.Close()
	} else {
		sqlStatement = fmt.Sprintf(`SELECT %s FROM "node" WHERE id_user=$1 OR is_public=true`, u.nodeField())
		rows, err = tx.Query(ctx, sqlStatement, currentUser.IdUser)
		if err != nil {
			return nodes, err
		}
		defer rows.Close()
	}

	for rows.Next() {
		var node entities.Node
		err := rows.Scan(
			u.nodePointer(&node)...,
		)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nodes, err
	}
	return nodes, nil
}

func (u *NodeRepository) GetById(ctx context.Context, tx helper.Querier, id int) (node entities.Node, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "node" WHERE id_node=$1`, u.nodeField())
	err = tx.QueryRow(ctx, sqlStatement, id).Scan(
		u.nodePointer(&node)...,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return node, fiber.NewError(404, fmt.Sprintf("Node with id %d not found", id))
		}
		return node, err
	}
	return node, nil
}

func (u *NodeRepository) GetHardwareNode(ctx context.Context, tx helper.Querier, hardwareId int) ([]entities.Node, error) {
	nodes := []entities.Node{}
	sqlStatement := fmt.Sprintf(`SELECT %s FROM "node" WHERE id_hardware=$1`, u.nodeField())
	rows, err := tx.Query(ctx, sqlStatement, hardwareId)
	if err != nil {
		return nodes, err
	}
	defer rows.Close()

	for rows.Next() {
		var node entities.Node
		err := rows.Scan(
			u.nodePointer(&node)...,
		)
		if err != nil {
			return nodes, err
		}
		nodes = append(nodes, node)
	}
	if err := rows.Err(); err != nil {
		return nodes, err
	}

	return nodes, nil
}

func (u *NodeRepository) Update(ctx context.Context, tx helper.Querier, node *entities.Node, payload *entities.NodeUpdate) (err error) {
	payload.ChangeSettedFieldOnly(node)

	sqlStatement := `
	UPDATE "node"
	SET 
	name=$1, 
	location=$2,
	id_hardware_node=$3,
	id_hardware_sensor=$4,
	field_sensor=$5
	WHERE id_node=$6`
	res, err := tx.Exec(ctx, sqlStatement, payload.Name, payload.Location, payload.IdHardwareNode, payload.IdHardwareSensor, payload.FieldSensor, node.IdNode)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on update node with id %d", node.IdNode))
	}
	return nil
}

func (u *NodeRepository) Delete(ctx context.Context, tx helper.Querier, id int) (err error) {
	sqlStatement := `DELETE FROM "node" WHERE id_node=$1`
	res, err := tx.Exec(ctx, sqlStatement, id)
	if err != nil {
		return err
	}
	count := res.RowsAffected()
	if count == 0 {
		return fiber.NewError(404, fmt.Sprintf("No row affected on delete with id %d", id))
	}
	return nil
}
