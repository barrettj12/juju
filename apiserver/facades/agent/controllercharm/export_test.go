// Copyright 2023 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package controllercharm

func NewAPI(state backend) *API {
	return &API{state}
}
