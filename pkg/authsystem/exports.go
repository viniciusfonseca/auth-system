package authsystem

import (
	authsystemroutes "github.com/viniciusfonseca/auth-system/internal/routes"
	authsystem "github.com/viniciusfonseca/auth-system/internal/shared"
)

var NewAuthSystem = authsystem.NewAuthSystem

var AddAuthGroup = authsystemroutes.AddAuthGroup
