package environment

import "backend/core/constants"

type Environment struct {
	env string
}

func New(env string) *Environment {
	return &Environment{
		env: env,
	}
}

func (e *Environment) Get() string {
	return e.env
}

func (e *Environment) IsProduction() bool {
	return e.env == constants.EnvironmentProduction
}

func (e *Environment) IsDevelopment() bool {
	return e.env == constants.EnvironmentDevelopment
}

func (e *Environment) IsTesting() bool {
	return e.env == constants.EnvironmentTesting
}
