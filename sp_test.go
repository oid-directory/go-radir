package radir

func bogusSponsor_codecov() error {
	var nilReg *Sponsor = new(Sponsor)

	if err := testBogusSponsorSetters(nilReg); err != nil {
		return err
	}
	if err := testBogusSponsorGetters(nilReg); err != nil {
		return err
	}

	nilReg.Auxiliary()
	nilReg.CN()
	nilReg.CO()
	nilReg.C()
	nilReg.O()
	nilReg.L()
	nilReg.Fax()
	nilReg.Email()
	nilReg.EndTime()
	nilReg.StartTime()
	nilReg.Mobile()
	nilReg.POBox()
	nilReg.PostalAddress()
	nilReg.PostalCode()
	nilReg.ST()
	nilReg.Street()
	nilReg.Tel()
	nilReg.Title()
	nilReg.URI()

	return nil
}

func testBogusSponsorGetters(nilReg *Sponsor) error {
	for _, funk := range []func(GetOrSetFunc) (any, error){
		nilReg.CNGetFunc,
		nilReg.COGetFunc,
		nilReg.CGetFunc,
		nilReg.OGetFunc,
		nilReg.LGetFunc,
		nilReg.FaxGetFunc,
		nilReg.EmailGetFunc,
		nilReg.EndTimeGetFunc,
		nilReg.StartTimeGetFunc,
		nilReg.MobileGetFunc,
		nilReg.POBoxGetFunc,
		nilReg.PostalAddressGetFunc,
		nilReg.PostalCodeGetFunc,
		nilReg.STGetFunc,
		nilReg.StreetGetFunc,
		nilReg.TelGetFunc,
		nilReg.TitleGetFunc,
		nilReg.URIGetFunc,
	} {
		funk(nil)
		val, err := funk(func(v ...any) (any, error) {
			// fake, do whatever
			return v[0], nil
		})

		if err != nil {
			return err
		}

		txt := `testing`
		if assert, ok := val.(string); ok {
			if assert != txt {
				return errorf("Mismatched; want '%s', got '%s'",
					txt, assert)
			}
		} else if sassert, sok := val.([]string); sok {
			if len(sassert) == 1 {
				if sassert[0] != txt {
					return errorf("Mismatched; want '%s', got '%s'",
						txt, sassert[0])
				}
			}
		} else {
			return errorf("Unsupported type '%T'", val)
		}
	}

	return nil
}

func testBogusSponsorSetters(nilReg *Sponsor) error {
	var bogus any = []int{1, 2, 3, 4}

	for _, funk := range []func(...any) error{
		nilReg.SetCN,
		nilReg.SetCO,
		nilReg.SetC,
		nilReg.SetO,
		nilReg.SetL,
		nilReg.SetFax,
		nilReg.SetEmail,
		nilReg.SetEndTime,
		nilReg.SetStartTime,
		nilReg.SetMobile,
		nilReg.SetPOBox,
		nilReg.SetPostalAddress,
		nilReg.SetPostalCode,
		nilReg.SetST,
		nilReg.SetStreet,
		nilReg.SetTitle,
		nilReg.SetTel,
		nilReg.SetURI,
	} {
		for _, err := range []error{
			funk(),
			funk(bogus),
		} {
			if err == nil {
				return mkerr("Expected error, got nothing")
			}
		}

		if err := funk(`testing`); err != nil {
			return err
		}
	}

	var em *Sponsor
	em.marshal(func(any) error {
		return RegistrantValidityErr
	})
	nilReg.marshal(func(any) error {
		return RegistrantValidityErr
	})
	nilReg.marshal(func(any) error {
		return nil
	})
	nilReg.marshal(nil)

	return nil
}
