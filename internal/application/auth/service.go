package auth

import (
	"context"
	"dach-admin/internal/domain"
	"dach-admin/internal/ports/repositories"

	"golang.org/x/crypto/bcrypt"
)

type TokenIssuer interface {
	Issue(member domain.TeamMember) (string, error)
}

type Service struct {
	team   repositories.TeamRepository
	tokens TokenIssuer
}

func NewService(team repositories.TeamRepository, tokens TokenIssuer) *Service {
	return &Service{team: team, tokens: tokens}
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, *domain.TeamMember, error) {
	member, err := s.team.GetByEmail(ctx, email)
	if err != nil {
		return "", nil, domain.ErrUnauthenticated
	}
	if member.Status != domain.TeamStatusActive {
		return "", nil, domain.ErrUnauthenticated
	}
	if err := bcrypt.CompareHashAndPassword([]byte(member.PasswordHash), []byte(password)); err != nil {
		return "", nil, domain.ErrUnauthenticated
	}
	token, err := s.tokens.Issue(*member)
	if err != nil {
		return "", nil, err
	}
	member.PasswordHash = ""
	return token, member, nil
}
