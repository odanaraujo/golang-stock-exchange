package entity

type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition
}

type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{},
	}
}

func (i *Investor) AddAssertPosition(assertPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assertPosition)
}

func (i *Investor) UpdateAssertPosition(assertID string, qtdShares int) {
	assertPosition := i.GetAssertPosition(assertID)
	if assertPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assertID, qtdShares))
	} else {
		assertPosition.Shares += qtdShares
	}
}

func (i *Investor) GetAssertPosition(assertID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assertID {
			return assetPosition
		}
	}
	return nil
}

func NewInvestorAssetPosition(assetID string, qtdShares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  qtdShares,
	}
}
