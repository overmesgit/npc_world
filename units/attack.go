package units

import "time"

type Attack struct {
    IsAttacking      bool
    AttackTimer      time.Time
    CooldownTimer    time.Time
    Message          string
    Range            float64
    Damage           int
    AttackDuration   time.Duration
    CooldownDuration time.Duration
    HasDealtDamage   bool // New field to track if damage has been dealt
}

func NewAttack(atackRange float64) Attack {
    return Attack{
        IsAttacking:      false,
        Message:          "Attack!",
        Range:            atackRange,
        Damage:           20,
        AttackDuration:   500 * time.Millisecond,
        CooldownDuration: 2 * time.Second,
        HasDealtDamage:   false,
    }
}

func (a *Attack) TriggerAttack() bool {
    if time.Now().After(a.CooldownTimer) {
        a.IsAttacking = true
        a.AttackTimer = time.Now().Add(a.AttackDuration)
        a.CooldownTimer = time.Now().Add(a.CooldownDuration)
        a.HasDealtDamage = false
        return true
    }
    return false
}

func (a *Attack) Update() {
    if a.IsAttacking && time.Now().After(a.AttackTimer) {
        a.IsAttacking = false
        a.HasDealtDamage = false // Reset this flag when attack ends
    }
}
