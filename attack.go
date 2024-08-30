package main

import (
	"time"
)

type Attack struct {
	IsAttacking bool
	AttackTimer time.Time
	Message     string
}

func NewAttack() Attack {
	return Attack{
		IsAttacking: false,
		Message:     "Attack!",
	}
}

func (a *Attack) TriggerAttack() {
	a.IsAttacking = true
	a.AttackTimer = time.Now().Add(500 * time.Millisecond)
}

func (a *Attack) Update() {
	if a.IsAttacking && time.Now().After(a.AttackTimer) {
		a.IsAttacking = false
	}
}