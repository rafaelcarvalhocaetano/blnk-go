package blnkgo

import "errors"

func ValidateCreateTransacation(t CreateTransactionRequest) error {
	if t.Source != nil && len(t.Sources) > 0 {
		return errors.New("you can not use both Source and Sources")
	}

	if t.Source == nil && len(t.Sources) == 0 {
		return errors.New("you must use either Source or Sources")
	}

	if t.Destination != nil && len(t.Destinations) > 0 {
		return errors.New("you can not use both Destination and Destinations")
	}

	if t.Destination == nil && len(t.Destinations) == 0 {
		return errors.New("you must use either Destination or Destinations")
	}

	if t.Amount < 0 {
		return errors.New("you can not use a negative amount")
	}

	if len(t.Sources) > 0 {
		err := validateSources(t.Sources, t.Amount)
		if err != nil {
			return err
		}
	}

	if len(t.Destinations) > 0 {
		err := validateSources(t.Destinations, t.Amount)
		if err != nil {
			return err
		}
	}

	return nil
}

func validateSources(sources []Source, amount float64) error {
	//total amount of sources  must be equal to the amount
	total := 0.0
	hasLeft := false
	for _, source := range sources {
		distribution := source.Distribution
		//check if the distribution is valid
		isValid := distribution.IsValid()
		if !isValid {
			return errors.New("invalid distribution")
		}

		switch {
		case distribution.IsPercentage():
			// Get float value from percentage
			percentage := distribution.ToPercentage()
			v := (percentage / 100) * amount
			if v < 0 {
				return errors.New("invalid distribution in source: " + source.Identifier)
			}
			total += v

		case distribution.IsNumber():
			// Get float value from number
			number := distribution.ToNumber()
			if number < 0 {
				return errors.New("invalid distribution in source: " + source.Identifier)
			}
			total += number

		case distribution.IsLeft():
			// Ensure "left" distribution is used only once
			if hasLeft {
				return errors.New("you cannot use left distribution more than once")
			}
			hasLeft = true

		default:
			// Handle invalid or unrecognized distribution
			return errors.New("unknown distribution type in source: " + source.Identifier)
		}
	}

	// If "left" distribution is used, calculate its value and add to total
	if hasLeft {
		left := amount - total
		if left < 0 {
			return errors.New("total amount of sources exceeds the amount")
		}
		total += left
	}

	if total != amount {
		return errors.New("total amount of sources must be equal to the amount")
	}

	return nil
}
