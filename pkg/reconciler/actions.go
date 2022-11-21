package reconciler

// ReconcileClient to create, update and delete objects
type ReconcileClient interface {
	// Create the object in downstream system
	Create(p interface{}) (*Summary, error)
	// Update the object in downstream system
	Update(p interface{}) (*Summary, error)
	// Delete the object from the downstream system
	Delete(p interface{}) (*Summary, error)
}
