package authsystem

import (
	authsystemroutes "github.com/viniciusfonseca/auth-system/internal/routes"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
)

type AuthSystem = authsystem.AuthSystem

var AddAuthGroup = authsystemroutes.AddAuthGroup
