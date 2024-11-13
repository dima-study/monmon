package v1

func (x AvailabilityState) String() string {
	switch x {
	case Available:
		return "AVAILABLE"
	case Error:
		return "NOT AVAILABLE"
	case Disabled:
		return "DISABLED"
	case Unspecified:
		fallthrough
	default:
		return x.pbString()
	}
}
