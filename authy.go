package radir

/*
StartTime returns a [GeneralizedTime] instance based on the contents
of the underlying R_StartTime field value. If unset, a zero instance
is returned.
*/
func (r CurrentAuthority) StartTime() (st GeneralizedTime) {
        if len(r.R_StartTime) > 0 {
                st, _ = parseGeneralizedTime(r.R_StartTime)
        }

        return
}

/*
EndTime always returns a zero [GeneralizedTime] instance. This method
exists only to satisfy Go's interface signature requirements.  The
concept of an "end timestamp" does not apply to [CurrentAuthority]
instances.
*/
func (r CurrentAuthority) EndTime() (et GeneralizedTime) {
        return
}

/*
StartTime returns a [GeneralizedTime] instance based on the contents
of the underlying R_StartTime field value. If unset, a zero instance
is returned.
*/
func (r FirstAuthority) StartTime() (st GeneralizedTime) {
	if len(r.R_StartTime) > 0 {
		st, _ = parseGeneralizedTime(r.R_StartTime)
	}

	return
}

/*
EndTime returns a [GeneralizedTime] instance based on the contents
of the underlying R_EndTime field value. If unset, a zero instance
is returned.
*/
func (r FirstAuthority) EndTime() (et GeneralizedTime) {
        if len(r.R_EndTime) > 0 {
                et, _ = parseGeneralizedTime(r.R_EndTime)
        }

        return
}

/*
StartTime returns a [GeneralizedTime] instance based on the contents
of the underlying R_StartTime field value. If unset, a zero instance
is returned.
*/
func (r Sponsor) StartTime() (st GeneralizedTime) {
        if len(r.R_StartTime) > 0 {
                st, _ = parseGeneralizedTime(r.R_StartTime)
        }

        return
}

/*
EndTime returns a [GeneralizedTime] instance based on the contents
of the underlying R_EndTime field value. If unset, a zero instance
is returned.
*/
func (r Sponsor) EndTime() (et GeneralizedTime) {
        if len(r.R_EndTime) > 0 {
                et, _ = parseGeneralizedTime(r.R_EndTime)
        }

        return
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *FirstAuthority) IsZero() bool {
	return r == nil
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *CurrentAuthority) IsZero() bool {
	return r == nil
}

/*
IsZero returns a Boolean value indicative of a nil receiver state.
*/
func (r *Sponsor) IsZero() bool {
	return r == nil
}

