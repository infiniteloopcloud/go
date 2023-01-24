package bitmask

type Bitmask uint64

func (f Bitmask) HasFlag(flag Bitmask) bool { return f&flag != 0 }

func (f *Bitmask) AddFlag(flag Bitmask) { *f |= flag }

func (f *Bitmask) ClearFlag(flag Bitmask) { *f &= ^flag }

func (f *Bitmask) ToggleFlag(flag Bitmask) { *f ^= flag }
