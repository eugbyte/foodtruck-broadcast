package subscription

// Create or update the web push subscription.
//
// Loose coupling between the lambda and business logic in case, e.g. change cloud provider.
func (h *handler) subscribe(sub Subscription) error {
	return h.subRepo.Put(sub)
}
