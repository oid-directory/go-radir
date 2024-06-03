package radir

/*
Type returns the string name of the type of [Registration].
*/
func (r ITUT) Type() string { return `ITU-T` }

/*
Type returns the string name of the type of [Registration].
*/
func (r ISO) Type() string { return `ISO` }

/*
Type returns the string name of the type of [Registration].
*/
func (r JointISOITUT) Type() string { return `Joint-ISO-ITU-T` }

/*
Type returns the string name of the type of [Registration].
*/
func (r RootArc) Type() (s string) {
	switch r {
	case ITUTRoot:
		s = `Root ITU-T`
	case ISORoot:
		s = `Root ISO`
	case JointRoot:
		s = `Root Joint-ISO-ITU-T`
	}

	return
}

/*
N returns the [NumberForm] associated with the receiver instance.
*/
func (r ITUT) N() (nf NumberForm) {
	nf, _ = parseNumberForm(r.Common.R_N)
	return
}

/*
N returns the [NumberForm] associated with the receiver instance.
*/
func (r ISO) N() (nf NumberForm) {
	nf, _ = parseNumberForm(r.Common.R_N)
	return
}

/*
N returns the [NumberForm] associated with the receiver instance.
*/
func (r JointISOITUT) N() (nf NumberForm) {
	nf, _ = parseNumberForm(r.Common.R_N)
	return
}

/*
DotNotation returns the [DotNotation] associated with the receiver instance.
*/
func (r ITUT) DotNotation() (dot DotNotation) {
	dot, _ = parseDotNotation(r.X680.R_DNot)
        return
}

/*
DotNotation returns the [DotNotation] associated with the receiver instance.
*/
func (r ISO) DotNotation() (dot DotNotation) {
	dot, _ = parseDotNotation(r.X680.R_DNot)
        return
}

/*
DotNotation returns the [DotNotation] associated with the receiver instance.
*/
func (r JointISOITUT) DotNotation() (dot DotNotation) {
	dot, _ = parseDotNotation(r.X680.R_DNot)
        return
}

/*
ASN1Notation returns the [ASN1Notation] associated with the receiver instance.
*/
func (r ITUT) ASN1Notation() (anot ASN1Notation) {
	anot, _ = parseASN1Notation(r.X680.R_ANot)
        return
}

/*
ASN1Notation returns the [ASN1Notation] associated with the receiver instance.
*/
func (r ISO) ASN1Notation() (anot ASN1Notation) {
	anot, _ = parseASN1Notation(r.X680.R_ANot)
        return
}

/*
ASN1Notation returns the [ASN1Notation] associated with the receiver instance.
*/
func (r JointISOITUT) ASN1Notation() (anot ASN1Notation) {
	anot, _ = parseASN1Notation(r.X680.R_ANot)
        return
}

/*
ASN1Notation returns the [ASN1Notation] associated with the receiver instance.
*/
func (r RootArc) ASN1Notation() (anot ASN1Notation) {
	anot, _ = parseASN1Notation(`{` + r.NameAndNumberForm().String() + `}`)
	return
}

/*
N returns the [NumberForm] associated with the receiver instance.
*/
func (r RootArc) N() (nf NumberForm) {
	switch r {
	case ITUTRoot:
		nf, _ = parseNumberForm(0)
	case ISORoot:
		nf, _ = parseNumberForm(1)
	case JointRoot:
		nf, _ = parseNumberForm(2)
	}

	return
}

/*
Identifier returns the string identifier associated with the receiver instance.
*/
func (r ITUT) Identifier() (id string) {
	id = r.X680.R_Id
	return
}

/*
Identifier returns the string identifier associated with the receiver instance.
*/
func (r ISO) Identifier() (id string) {
	id = r.X680.R_Id
	return
}

/*
Identifier returns the string identifier associated with the receiver instance.
*/
func (r JointISOITUT) Identifier() (id string) {
	id = r.X680.R_Id
	return
}

func (_ *RootArc) SetIdentifier(X any, setfunc ...GetOrSetFunc) error {
	return nil
}

func (r *ITUT) SetIdentifier(X any, setfunc ...GetOrSetFunc) error {
        if len(setfunc) == 0 {
                if assert, ok := X.(string); ok {
			r.X680.R_Id = assert
                        return nil
                }
                return mkerr("Unsupported Identifier type provided without GetOrSetFunc instance")
        }

        v, err := setfunc[0](X, r)
        if err != nil {
                return err
        }
        return r.SetIdentifier(v)
}

func (r *ISO) SetIdentifier(X any, setfunc ...GetOrSetFunc) error {    
        if len(setfunc) == 0 {                                          
                if assert, ok := X.(string); ok {                       
                        r.X680.R_Id = assert                            
                        return nil                                      
                }                                                       
                return mkerr("Unsupported Identifier type provided without GetOrSetFunc instance")
        }                                                               
                                                                        
        v, err := setfunc[0](X, r)                                      
        if err != nil {                                                 
                return err                                              
        }                                                               
        return r.SetIdentifier(v)                                       
}

func (r *JointISOITUT) SetIdentifier(X any, setfunc ...GetOrSetFunc) error {    
        if len(setfunc) == 0 {                                          
                if assert, ok := X.(string); ok {                       
                        r.X680.R_Id = assert                            
                        return nil                                      
                }                                                       
                return mkerr("Unsupported Identifier type provided without GetOrSetFunc instance")
        }                                                               
                                                                        
        v, err := setfunc[0](X, r)                                      
        if err != nil {                                                 
                return err                                              
        }                                                               
        return r.SetIdentifier(v)                                       
}

func (r *ITUT) setIdentifier(X string) {
	r.R_Id = X
}

func (r *ISO) setIdentifier(X string) {
	r.R_Id = X
}

func (r *JointISOITUT) setIdentifier(X string) {
	r.R_Id = X
}

/*
Identifier returns the string identifier associated with the receiver instance.
*/
func (r RootArc) Identifier() (id string) {
	switch r {
	case ITUTRoot:
		id = `itu-t`
	case ISORoot:
		id = `iso`
	case JointRoot:
		id = `joint-iso-itu-t`
	}

	return
}

/*
NameAndNumberForm returns the [NameAndNumberForm] associated with the receiver instance.
*/
func (r RootArc) NameAndNumberForm() (nanf NameAndNumberForm) {
	nanf, _ = parseNameAndNumberForm(r.Identifier() + `(` + r.N().String() + `)`)
	return
}

/*
NameAndNumberForm returns the [NameAndNumberForm] associated with the receiver instance.
*/
func (r ITUT) NameAndNumberForm() (nanf NameAndNumberForm) {
	nanf, _ = parseNameAndNumberForm(r.Identifier() + `(` + r.N().String() + `)`)
	return
}

/*
NameAndNumberForm returns the [NameAndNumberForm] associated with the receiver instance.
*/
func (r ISO) NameAndNumberForm() (nanf NameAndNumberForm) {
	nanf, _ = parseNameAndNumberForm(r.Identifier() + `(` + r.N().String() + `)`)
	return
}

/*
NameAndNumberForm returns the [NameAndNumberForm] associated with the receiver instance.
*/
func (r JointISOITUT) NameAndNumberForm() (nanf NameAndNumberForm) {
	nanf, _ = parseNameAndNumberForm(r.Identifier() + `(` + r.N().String() + `)`)
	return
}

/*
CreateTime returns a zero [GeneralizedTime] instance. This method only exists
to satisfy Go's interface signature requirements.
*/
func (r RootArc) CreateTime() (ct GeneralizedTime) {
	return
}

/*
CreateTime returns a [GeneralizedTime] instance based on the contents of
the underlying R_CTime field value. If unset, a zero instance is returned.
*/
func (r ITUT) CreateTime() (ct GeneralizedTime) {
	ct, _ = parseGeneralizedTime(r.Extra.R_CTime)
	return
}

/*
CreateTime returns a [GeneralizedTime] instance based on the contents of
the underlying R_CTime field value. If unset, a zero instance is returned.
*/
func (r ISO) CreateTime() (ct GeneralizedTime) {
	ct, _ = parseGeneralizedTime(r.Extra.R_CTime)
	return
}

/*
CreateTime returns a [GeneralizedTime] instance based on the contents of
the underlying R_CTime field value. If unset, a zero instance is returned.
*/
func (r JointISOITUT) CreateTime() (ct GeneralizedTime) {
	ct, _ = parseGeneralizedTime(r.Extra.R_CTime)
	return
}

/*
ModifyTime returns zero slices of [GeneralizedTime] instances. This method
only exists to satisfy Go's interface signature requirements.
*/
func (r RootArc) ModifyTime() (mts []GeneralizedTime) { return }

/*
ModifyTime returns slices of [GeneralizedTime] instances based on the
contents of the underlying R_MTime field value. If unset, a zero instance
is returned.
*/
func (r ITUT) ModifyTime() (mts []GeneralizedTime) {
        for _, slice := range r.Extra.R_MTime {
                if mt, _ := parseGeneralizedTime(slice); !mt.IsZero() {
                        mts = append(mts, mt)
                }
        }

	return
}

/*                                                                      
ModifyTime returns slices of [GeneralizedTime] instances based on the   
contents of the underlying R_MTime field value. If unset, a zero instance
is returned.                                                            
*/                                                                      
func (r ISO) ModifyTime() (mts []GeneralizedTime) {                    
        for _, slice := range r.Extra.R_MTime {                         
                if mt, _ := parseGeneralizedTime(slice); !mt.IsZero() { 
                        mts = append(mts, mt)                           
                }                                                       
        }                                                               
                                                                        
        return                                                          
}

/*                                                                      
ModifyTime returns slices of [GeneralizedTime] instances based on the   
contents of the underlying R_MTime field value. If unset, a zero instance
is returned.                                                            
*/                                                                      
func (r JointISOITUT) ModifyTime() (mts []GeneralizedTime) {                    
        for _, slice := range r.Extra.R_MTime {                         
                if mt, _ := parseGeneralizedTime(slice); !mt.IsZero() { 
                        mts = append(mts, mt)                           
                }                                                       
        }                                                               
                                                                        
        return                                                          
}
