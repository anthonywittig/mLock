package lockengine

// Some tests to consider:
// * do nothing
//   * no managed lock codes
//   * managed lock code but no action to take
//     * codes already exist
//     * codes don't exist and don't need to be created
// * add lock codes
// * remove lock codes
// * add and remove lock codes

/*
func Test_NoLockCodes(t *testing.T) {
	le, dc, dr, pr := newLockEngine()
	err := le.UpdateLocks(ctx)
	assert.Nil(t, err)
}

type mockDeviceController struct{}
type mockDeviceRepository struct{}
type mockPropertyRepository struct{}

func newLockEngine() (*LockEngine, *mockDeviceController, *mockDeviceRepository, *mockPropertyRepository) {
	dc := &mockDeviceController{}
	dr := &mockDeviceRepository{}
	pr := &mockPropertyRepository{}
	le := NewLockEngine(dc, dr, pr)
	return le, dc, dr, pr
}

*/
