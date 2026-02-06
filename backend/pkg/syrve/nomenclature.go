package syrve

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
)

func (c *Client) GetNomenclature(ctx context.Context, authToken string, req NomenclatureRequest) (*NomenclatureResponse, error) {
	if authToken == "" {
		return nil, errors.New("authorization token is required")
	}
	if req.OrganizationID == "" {
		return nil, errors.New("organizationId is required")
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	reqHTTP, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/nomenclature", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	reqHTTP.Header.Set("Content-Type", "application/json")
	reqHTTP.Header.Set("Authorization", "Bearer "+authToken)

	resp, err := c.HttpClient.Do(reqHTTP)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code: " + resp.Status)
	}

	var nomenResp NomenclatureResponse
	err = json.NewDecoder(resp.Body).Decode(&nomenResp)
	if err != nil {
		return nil, err
	}

	return &nomenResp, nil
}

func (c *Client) GetProducts(ctx context.Context, authToken, organizationID string) ([]MenuItemDTO, error) {
	resp, err := c.GetNomenclature(ctx, authToken, NomenclatureRequest{OrganizationID: organizationID})
	if err != nil {
		return nil, err
	}

	productMap := make(map[string]MenuItem)
	for _, p := range resp.Products {
		if p.Type != nil && *p.Type != "Modifier" {
			productMap[p.ID] = p
		} else if p.Type != nil {
			// log modifiers if needed for debug
		}
	}

	modifiersMap := make(map[string]MenuItem)
	for _, p := range resp.Products {
		if *p.Type == "Modifier" {
			modifiersMap[p.ID] = p
		}
	}

	groupMap := make(map[string]Group)
	for _, group := range resp.Groups {
		groupMap[group.ID] = group
	}

	var result []MenuItemDTO
	for _, p := range productMap {
		code := ""
		if p.Code != nil {
			code = *p.Code
		}
		pType := ""
		if p.Type != nil {
			pType = *p.Type
		}
		full := MenuItemDTO{
			ID:             p.ID,
			Code:           code,
			Name:           p.Name,
			Type:           pType,
			Modifiers:      []ModifierDTO{},
			GroupModifiers: []ModifierGroupDTO{},
		}

		for _, modifier := range p.Modifiers {
			modifierFull := ModifierDTO{
				ID:            modifier.ID,
				Name:          modifiersMap[modifier.ID].Name,
				DefaultAmount: modifier.DefaultAmount,
				MinAmount:     modifier.MinAmount,
				MaxAmount:     modifier.MaxAmount,
				Required:      modifier.Required,
			}

			full.Modifiers = append(full.Modifiers, modifierFull)
		}

		for _, groupModifier := range p.GroupModifiers {
			groupModifierFull := ModifierGroupDTO{
				ID:             groupModifier.ID,
				Name:           groupMap[groupModifier.ID].Name,
				Required:       groupModifier.Required,
				MinAmount:      groupModifier.MinAmount,
				MaxAmount:      groupModifier.MaxAmount,
				DefaultAmount:  groupModifier.DefaultAmount,
				ChildModifiers: []ModifierDTO{},
			}

			for _, cm := range groupModifier.ChildModifiers {
				m := ModifierDTO{
					ID:            cm.ID,
					Name:          modifiersMap[cm.ID].Name,
					DefaultAmount: cm.DefaultAmount,
					MinAmount:     cm.MinAmount,
					MaxAmount:     cm.MaxAmount,
					Required:      cm.Required,
				}
				groupModifierFull.ChildModifiers = append(groupModifierFull.ChildModifiers, m)
			}

			full.GroupModifiers = append(full.GroupModifiers, groupModifierFull)
		}

		result = append(result, full)
	}

	return result, nil
}
