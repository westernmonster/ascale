//go:build worker
// +build worker

package service

func (s *Service) startSubscriptions() {
	go s.subscriptions()
}