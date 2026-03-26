package git

func (h Handler) GetCurrentBranch() (string, error) {
	res, err := h.exec("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}

	return res, nil
}

func (h Handler) FindMasterBranch() string {
	var res string
	var err error

	// try to get the remotes HEAD
	res, err = h.exec("symbolic-ref", "--short", "refs/remotes/origin/HEAD")
	if err == nil {
		return res
	}

	// TODO: try other remotes if origin isnt available?

	// nothing else worked, get the current branch name
	res, err = h.GetCurrentBranch()
	if err == nil {
		return res
	}

	return ""
}

func (h Handler) BranchExists(name string) bool {
	_, err := h.exec("rev-parse", "--verify", name)
	return err == nil // returns hash on success and error on failure
}
