package syrve

type NomenclatureRequest struct {
	OrganizationID string `json:"organizationId"`
	StartRevision  *int64 `json:"startRevision,omitempty"`
}

type NomenclatureResponse struct {
	Revision           int64              `json:"revision"`
	CorrelationID      string             `json:"correlationId"`
	StockListGroups    []StockListGroup   `json:"stockListGroups"`
	MenuItemCategories []MenuItemCategory `json:"menuItemCategories"`
	Products           []MenuItem         `json:"products"`
	ItemSizes          []ItemSize         `json:"itemSizes"`
	Groups             []Group            `json:"groups"`
}

type StockListGroup struct {
	ID               string   `json:"id"`
	Code             *string  `json:"code,omitempty"`
	Name             string   `json:"name"`
	Description      *string  `json:"description,omitempty"`
	AdditionalInfo   *string  `json:"additionalInfo,omitempty"`
	Tags             []string `json:"tags,omitempty"`
	IsDeleted        bool     `json:"isDeleted"`
	SEODescription   *string  `json:"seoDescription,omitempty"`
	SEOText          *string  `json:"seoText,omitempty"`
	SEOKeywords      *string  `json:"seoKeywords,omitempty"`
	SEOTitle         *string  `json:"seoTitle,omitempty"`
	ImageLinks       []string `json:"imageLinks"`
	ParentGroup      *string  `json:"parentGroup,omitempty"`
	Order            int      `json:"order"`
	IsIncludedInMenu bool     `json:"isIncludedInMenu"`
	IsGroupModifier  bool     `json:"isGroupModifier"`
}

type MenuItemCategory struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	IsDeleted bool   `json:"isDeleted"`
}

type MenuItem struct {
	ID             string   `json:"id"`
	Code           *string  `json:"code,omitempty"`
	Name           string   `json:"name"`
	Description    *string  `json:"description,omitempty"`
	AdditionalInfo *string  `json:"additionalInfo,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	IsDeleted      bool     `json:"isDeleted"`
	SEODescription *string  `json:"seoDescription,omitempty"`
	SEOText        *string  `json:"seoText,omitempty"`
	SEOKeywords    *string  `json:"seoKeywords,omitempty"`
	SEOTitle       *string  `json:"seoTitle,omitempty"`

	FatAmount           *float64 `json:"fatAmount,omitempty"`
	ProteinsAmount      *float64 `json:"proteinsAmount,omitempty"`
	CarbohydratesAmount *float64 `json:"carbohydratesAmount,omitempty"`
	EnergyAmount        *float64 `json:"energyAmount,omitempty"`

	FatFullAmount           *float64 `json:"fatFullAmount,omitempty"`
	ProteinsFullAmount      *float64 `json:"proteinsFullAmount,omitempty"`
	CarbohydratesFullAmount *float64 `json:"carbohydratesFullAmount,omitempty"`
	EnergyFullAmount        *float64 `json:"energyFullAmount,omitempty"`

	Weight            *float64 `json:"weight,omitempty"`
	GroupID           *string  `json:"groupId,omitempty"`
	ProductCategoryID *string  `json:"productCategoryId,omitempty"`

	Type          *string `json:"type,omitempty"` // dish | good | modifier
	OrderItemType string  `json:"orderItemType"`  // Product | Compound

	ModifierSchemaID   *string `json:"modifierSchemaId,omitempty"`
	ModifierSchemaName *string `json:"modifierSchemaName,omitempty"`
	Splittable         bool    `json:"splittable"`
	MeasureUnit        string  `json:"measureUnit"`

	Prices         []Price         `json:"prices,omitempty"`
	SizePrices     []Price         `json:"sizePrices,omitempty"`
	Modifiers      []Modifier      `json:"modifiers,omitempty"`
	GroupModifiers []ModifierGroup `json:"groupModifiers,omitempty"`

	ImageLinks         []string `json:"imageLinks,omitempty"`
	DoNotPrintInCheque bool     `json:"doNotPrintInCheque"`
	ParentGroup        *string  `json:"parentGroup,omitempty"`
	Order              int      `json:"order"`

	FullNameEnglish   *string `json:"fullNameEnglish,omitempty"`
	UseBalanceForSell bool    `json:"useBalanceForSell"`
	CanSetOpenPrice   bool    `json:"canSetOpenPrice"`
	PaymentSubject    *string `json:"paymentSubject,omitempty"`
}

type MenuItemDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Modifiers      []ModifierDTO      `json:"modifiers"`
	GroupModifiers []ModifierGroupDTO `json:"groupModifiers"`
}

type Group struct {
	ImageLinks       []string `json:"imageLinks"`
	ParentGroup      *string  `json:"parentGroup"`
	Order            int      `json:"order"`
	IsIncludedInMenu bool     `json:"isIncludedInMenu"`
	IsGroupModifier  bool     `json:"isGroupModifier"`
	ID               string   `json:"id"`
	Code             string   `json:"code"`
	Name             string   `json:"name"`
	Description      *string  `json:"description"`
	AdditionalInfo   *string  `json:"additionalInfo"`
	Tags             []string `json:"tags"`
	IsDeleted        bool     `json:"isDeleted"`
	SeoDescription   *string  `json:"seoDescription"`
	SeoText          *string  `json:"seoText"`
	SeoKeywords      *string  `json:"seoKeywords"`
	SeoTitle         *string  `json:"seoTitle"`
}

type Price struct {
	SizeID    *string   `json:"sizeId,omitempty"`
	PriceData PriceData `json:"price"`
}

type PriceData struct {
	Amount       float64 `json:"amount,omitempty"`
	CurrencyCode string  `json:"currencyCode,omitempty"`
}

type ItemSize struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Priority  *int   `json:"priority,omitempty"`
	IsDefault *bool  `json:"isDefault,omitempty"`
	Revision  int64  `json:"revision"`
}

type Modifier struct {
	ID                  string `json:"id"`
	DefaultAmount       *int   `json:"defaultAmount,omitempty"`
	MinAmount           int    `json:"minAmount"`
	MaxAmount           int    `json:"maxAmount"`
	Required            *bool  `json:"required,omitempty"`
	HideIfDefaultAmount *bool  `json:"hideIfDefaultAmount,omitempty"`
	Splittable          *bool  `json:"splittable,omitempty"`
	FreeOfChargeAmount  *int   `json:"freeOfChargeAmount,omitempty"`
}

type ModifierDTO struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	DefaultAmount *int   `json:"defaultAmount,omitempty"`
	Required      *bool  `json:"required,omitempty"`
}

type ModifierGroup struct {
	ID                                   string     `json:"id"`
	MinAmount                            int        `json:"minAmount"`
	MaxAmount                            int        `json:"maxAmount"`
	Required                             *bool      `json:"required,omitempty"`
	ChildModifiersHaveMinMaxRestrictions *bool      `json:"childModifiersHaveMinMaxRestrictions,omitempty"`
	HideIfDefaultAmount                  *bool      `json:"hideIfDefaultAmount,omitempty"`
	DefaultAmount                        *int       `json:"defaultAmount,omitempty"`
	Splittable                           *bool      `json:"splittable,omitempty"`
	FreeOfChargeAmount                   *int       `json:"freeOfChargeAmount,omitempty"`
	ChildModifiers                       []Modifier `json:"childModifiers"`
}

// -----------------------------
// FULL STRUCTURES & DTOs
// -----------------------------

type MenuItemFull struct {
	MenuItem
}

type GroupModifierFull struct {
	ModifierGroup
	Name           string     `json:"name"`
	ChildModifiers []MenuItem `json:"childModifiers"`
}

type ModifierGroupDTO struct {
	ID             string        `json:"id"`
	Name           string        `json:"name"`
	Required       *bool         `json:"required,omitempty"`
	DefaultAmount  *int          `json:"defaultAmount,omitempty"`
	ChildModifiers []ModifierDTO `json:"childModifiers"`
}
