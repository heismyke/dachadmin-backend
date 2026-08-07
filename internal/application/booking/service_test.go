package booking

import (
	"context"
	"errors"
	"testing"
	"time"

	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"github.com/google/uuid"
)

type fakeBookingRepo struct {
	assigned repositories.AssignDriverParams
	status   domain.BookingStatus
}

func (f *fakeBookingRepo) Create(ctx context.Context, b *domain.Booking) error {
	b.ID = uuid.New()
	return nil
}
func (f *fakeBookingRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Booking, error) {
	return nil, domain.ErrNotFound
}
func (f *fakeBookingRepo) List(ctx context.Context, filter repositories.BookingFilter) (repositories.ListResult[domain.Booking], error) {
	return repositories.ListResult[domain.Booking]{}, nil
}
func (f *fakeBookingRepo) Update(ctx context.Context, b *domain.Booking) error { return nil }
func (f *fakeBookingRepo) Delete(ctx context.Context, id uuid.UUID) error      { return nil }
func (f *fakeBookingRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status domain.BookingStatus, actorID *uuid.UUID) error {
	f.status = status
	return nil
}
func (f *fakeBookingRepo) AssignDriver(ctx context.Context, params repositories.AssignDriverParams) error {
	f.assigned = params
	return nil
}

func TestCreateBookingRequiresPickupDate(t *testing.T) {
	svc := NewService(&fakeBookingRepo{})
	err := svc.Create(context.Background(), &domain.Booking{PickupAddress: "A", DropoffAddress: "B"})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestAssignDriverDelegatesUseCaseToRepositoryPort(t *testing.T) {
	repo := &fakeBookingRepo{}
	svc := NewService(repo)
	bookingID := uuid.New()
	driverID := uuid.New()
	actorID := uuid.New()
	if err := svc.AssignDriver(context.Background(), bookingID, driverID, &actorID); err != nil {
		t.Fatalf("assign driver: %v", err)
	}
	if repo.assigned.BookingID != bookingID || repo.assigned.DriverID != driverID || *repo.assigned.ActorID != actorID {
		t.Fatalf("assignment params not passed through")
	}
}

func TestCreateBookingDefaultsPending(t *testing.T) {
	repo := &fakeBookingRepo{}
	svc := NewService(repo)
	err := svc.Create(context.Background(), &domain.Booking{
		CustomerID:     uuid.New(),
		ServiceID:      uuid.New(),
		PickupAddress:  "A",
		DropoffAddress: "B",
		PickupDate:     time.Now(),
	})
	if err != nil {
		t.Fatalf("create booking: %v", err)
	}
}
