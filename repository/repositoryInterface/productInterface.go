package repository

import (
	"context"

	"zog/domain/entity"
)

type ProductInterface interface {
	// Ticket related methods
	CreateTicket(ctx context.Context, ticket *entity.Ticket) error
	GetTicketByID(ctx context.Context, id uint) (*entity.Ticket, error)
	UpdateTicket(ctx context.Context, ticket *entity.Ticket) error
	DeleteTicket(ctx context.Context, ticket *entity.Ticket) error

	// TicketDetails related methods
	CreateTicketDetails(ctx context.Context, details *entity.TicketDetails) error
	GetTicketDetailsByID(ctx context.Context, id uint) (*entity.TicketDetails, error)
	UpdateTicketDetails(ctx context.Context, details *entity.TicketDetails) error
	DeleteTicketDetails(ctx context.Context, details *entity.TicketDetails) error

	// Apparel related methods
	CreateApparel(ctx context.Context, apparel *entity.Apparel) error
	GetApparelByID(ctx context.Context, id uint) (*entity.Apparel, error)
	UpdateApparel(ctx context.Context, apparel *entity.Apparel) error
	DeleteApparel(ctx context.Context, apparel *entity.Apparel) error
}
