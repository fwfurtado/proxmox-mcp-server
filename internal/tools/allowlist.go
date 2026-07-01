package tools

type Allowlist map[string]struct{}

func NewAllowlist(items []string) Allowlist {
	if len(items) == 0 {
		return nil
	}

	allowlist := make(Allowlist, len(items))
	for _, item := range items {
		if item == "*" {
			return nil
		}
		allowlist[item] = struct{}{}
	}

	return allowlist
}

func (a Allowlist) Allows(name string) bool {
	if a == nil {
		return true
	}

	_, ok := a[name]
	return ok
}
