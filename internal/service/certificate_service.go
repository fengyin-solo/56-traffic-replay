package service

import (
	"sort"
	"time"

	"replay/internal/model"
	"replay/pkg/idgen"
)

func (s *Service) CreateCertificate(input model.Certificate) (*model.Certificate, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c := &model.Certificate{
		ID:        idgen.Hex(),
		Name:      input.Name,
		CertPEM:   input.CertPEM,
		KeyPEM:    input.KeyPEM,
		ExpiresAt: input.ExpiresAt,
		CreatedAt: time.Now(),
	}
	if err := s.store.CreateCertificate(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) ListCertificates(filter model.CertificateFilter, page, size int) ([]*model.Certificate, int, error) {
	all := s.store.ListCertificates()
	matched := make([]*model.Certificate, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Certificate{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) GetCertificate(id string) (*model.Certificate, error) {
	return s.store.GetCertificate(id)
}

func (s *Service) UpdateCertificate(id string, input model.Certificate) (*model.Certificate, error) {
	c, err := s.store.GetCertificate(id)
	if err != nil {
		return nil, err
	}
	c.Name = input.Name
	c.CertPEM = input.CertPEM
	c.KeyPEM = input.KeyPEM
	c.ExpiresAt = input.ExpiresAt
	if err := c.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateCertificate(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCertificate(id string) error {
	return s.store.DeleteCertificate(id)
}
