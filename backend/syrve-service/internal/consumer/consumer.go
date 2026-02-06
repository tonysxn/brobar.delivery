package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"github.com/tonysanin/brobar/pkg/rabbitmq"
	"github.com/tonysanin/brobar/pkg/syrve"
	"github.com/tonysanin/brobar/syrve-service/internal/models"
)

type Consumer struct {
	client       *syrve.Client
	producer     *rabbitmq.Producer
	currentOrder *models.OrderEvent // Not used, maybe for debugging
}

func NewConsumer(client *syrve.Client, producer *rabbitmq.Producer) *Consumer {
	return &Consumer{client: client, producer: producer}
}

func (c *Consumer) Start(uri string) {
	// Initialize RabbitMQ Consumer Wrapper
	rmqConsumer, err := rabbitmq.NewConsumer(uri)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rmqConsumer.Close()

	msgs, err := rmqConsumer.Consume(rabbitmq.QueueSyrve)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	// Consumer 2: Sync Triggers
	syncMsgs, err := rmqConsumer.Consume("syrve.sync.start")
	if err != nil {
		log.Printf("Failed to consume sync triggers: %v", err)
	}

	log.Println("Waiting for Syrve orders and sync triggers...")

	forever := make(chan bool)

	go func() {
		for {
			select {
			case d := <-msgs:
				log.Printf("Received order: %s", d.Body)

				var order models.OrderEvent
				if err := json.Unmarshal(d.Body, &order); err != nil {
					log.Printf("Error unmarshalling order: %v", err)
					continue
				}

				if err := c.processOrder(context.Background(), &order); err != nil {
					log.Printf("Error processing order %s: %v", order.ID, err)
				} else {
					log.Printf("Successfully sent order %s to Syrve", order.ID)
				}
			case d := <-syncMsgs:
				log.Printf("Received sync trigger: %s", d.Body)
				c.handleSync(d)
			}
		}
	}()

	<-forever
}

func (c *Consumer) handleSync(d amqp.Delivery) {
	var payload struct {
		ChatID int64 `json:"chat_id"`
	}
	_ = json.Unmarshal(d.Body, &payload)
	
	// Fetch Stop Lists
	ctx := context.Background()
	tokenResp, err := c.client.GetAccessToken(ctx)
	if err != nil {
		log.Printf("Sync failed (Auth): %v", err)
		return
	}

	stopLists, err := c.client.GetStopLists(ctx, tokenResp.Token, c.client.OrganizationID)
	if err != nil {
		log.Printf("Sync failed (Fetch): %v", err)
		return
	}

	// Flatten
	type StopListEventItem struct {
		ProductID string  `json:"product_id"`
		Balance   float64 `json:"balance"`
	}
	var items []StopListEventItem

	for _, orgList := range stopLists.TerminalGroupStopLists {
		for _, group := range orgList.Items {
			for _, item := range group.Items {
				items = append(items, StopListEventItem{
					ProductID: item.ProductID,
					Balance:   item.Balance,
				})
			}
		}
	}
	
	wrapper := map[string]interface{}{
		"items": items,
		"chat_id": payload.ChatID,
	}

	eventBytes, _ := json.Marshal(wrapper)
	if err := c.producer.SendMessage("syrve.stop_list.updated", string(eventBytes)); err != nil {
		log.Printf("Failed to publish stop list update: %v", err)
	} else {
		log.Printf("Published stop list update with %d items (ChatID: %d)", len(items), payload.ChatID)
	}
}

func (c *Consumer) enrichItemWithModifiers(item *syrve.OrderItem, product syrve.MenuItemDTO) {
	// 1. Simple Modifiers (direct children)
	for _, mod := range product.Modifiers {
		if mod.MinAmount > 0 {
			// Check if already present
			found := false
			for _, existing := range item.Modifiers {
				if existing.ProductID == mod.ID {
					found = true
					break
				}
			}
			if !found {
				amount := float64(mod.MinAmount)
				if mod.DefaultAmount != nil {
					amount = float64(*mod.DefaultAmount)
				}
				log.Printf("Auto-adding mandatory modifier %s (Amount: %.0f) for %s", mod.Name, amount, product.Name)
				item.Modifiers = append(item.Modifiers, syrve.OrderModifier{
					ProductID: mod.ID,
					Amount:    amount,
				})
			}
		}
	}

	// 2. Group Modifiers
	for _, group := range product.GroupModifiers {
		if group.MinAmount > 0 {
			// Check current total amount in this group
			currentAmount := 0.0
			childIDs := make(map[string]bool)
			for _, cm := range group.ChildModifiers {
				childIDs[cm.ID] = true
			}

			for _, existing := range item.Modifiers {
				if childIDs[existing.ProductID] {
					currentAmount += existing.Amount
				}
			}

			if currentAmount < float64(group.MinAmount) {
				// Add the first child modifier to satisfy MinAmount
				if len(group.ChildModifiers) > 0 {
					defaultMod := group.ChildModifiers[0]
					amount := float64(group.MinAmount) - currentAmount
					// If defaultMod has a DefaultAmount that's higher, use it? 
					// Actually, group minAmount is usually the constraint.
					if defaultMod.DefaultAmount != nil && float64(*defaultMod.DefaultAmount) > amount {
						amount = float64(*defaultMod.DefaultAmount)
					}

					log.Printf("Auto-adding mandatory group modifier %s from group %s for %s (Amount: %.0f)", defaultMod.Name, group.Name, product.Name, amount)
					item.Modifiers = append(item.Modifiers, syrve.OrderModifier{
						ProductID:      defaultMod.ID,
						ProductGroupID: group.ID,
						Amount:         amount,
					})
				}
			}
		}
	}
}

func (c *Consumer) processOrder(ctx context.Context, order *models.OrderEvent) error {
	// 1. Authenticate
	tokenResp, err := c.client.GetAccessToken(ctx)
	if err != nil {
		return fmt.Errorf("auth error: %w", err)
	}
	token := tokenResp.Token
	orgID := c.client.OrganizationID

	// 2. Get Organization ID if missing
	if orgID == "" {
		orgs, err := c.client.GetOrganizations(ctx, token, syrve.OrganizationsRequest{})
		if err != nil {
			return err
		}
		if len(orgs.Organizations) == 0 {
			return fmt.Errorf("no organizations found")
		}
		orgID = orgs.Organizations[0].ID
	}

	// 3. Get Terminal ID
	tGroups, err := c.client.GetTerminalGroups(ctx, token, orgID)
	if err != nil {
		return fmt.Errorf("failed to get terminals: %w", err)
	}
	
	log.Printf("Found %d terminal groups", len(tGroups.TerminalGroups))
	if len(tGroups.TerminalGroups) > 0 {
		for i, tg := range tGroups.TerminalGroups {
			log.Printf("Group [%d]: ID=%s, Name=%s, Items=%d", i, tg.ID, tg.Name, len(tg.Items))
		}
	}

	if len(tGroups.TerminalGroups) == 0 || len(tGroups.TerminalGroups[0].Items) == 0 {
		return fmt.Errorf("no terminals found")
	}
	// Logic from legacy: first group, first item ID is used as "terminalId"
	terminalID := tGroups.TerminalGroups[0].Items[0].ID
	log.Printf("Using terminalID: %s", terminalID)

	// 4. Find Restaurant Section "Доставка" and its tables
	availableSections, err := c.client.GetRestaurantSections(ctx, token, []string{terminalID})
	if err != nil {
		return fmt.Errorf("failed to get sections: %w", err)
	}

	var deliveryTableIDs []string
	for _, section := range availableSections.RestaurantSections {
		// Matching name "Доставка" as per legacy logic
		if section.Name == "Доставка" {
			for _, t := range section.Tables {
				deliveryTableIDs = append(deliveryTableIDs, t.ID)
			}
			break
		}
	}
	
	if len(deliveryTableIDs) == 0 {
		return fmt.Errorf("section 'Доставка' not found or has no tables")
	}

	// 5. Find Free Table (Legacy logic: table with oldest last order)
	// We need to fetch active orders for these tables
	dateFrom := time.Now().Add(-24 * time.Hour).Format("2006-01-02T15:04:05")
	dateTo := time.Now().Format("2006-01-02T15:04:05")
	activeOrdersResp, err := c.client.GetOrdersByTables(ctx, token, orgID, deliveryTableIDs, dateFrom, dateTo)
	if err != nil {
		return fmt.Errorf("failed to get active orders: %w", err)
	}
	
	// Map tableID -> last timestamp
	// Default to 0 (very old)
	tableTimestamps := make(map[string]int64)
	for _, tID := range deliveryTableIDs {
		tableTimestamps[tID] = 0
	}

	for _, o := range activeOrdersResp.Orders {
		if len(o.Order.TableIDs) > 0 {
			tID := o.Order.TableIDs[0]
			// We only care if it's in our list
			if _, exists := tableTimestamps[tID]; exists {
				// Keep the latest timestamp for this table
				if o.Timestamp > tableTimestamps[tID] {
					tableTimestamps[tID] = o.Timestamp
				}
			}
		}
	}

	// Find table with minimum timestamp (oldest activity = most free?)
	// Or maybe undefined behavior in legacy logic, mostly works if traffic is low.
	// We want key with min value.
	
	// Sorting logic
	type TableSort struct {
		ID   string
		Time int64
	}
	var sortedTables []TableSort
	for id, t := range tableTimestamps {
		sortedTables = append(sortedTables, TableSort{id, t})
	}
	
	sort.Slice(sortedTables, func(i, j int) bool {
		return sortedTables[i].Time < sortedTables[j].Time
	})
	
	selectedTableID := sortedTables[0].ID

	// 6. Find Order Type "Доставка БРО"
	orderTypesResp, err := c.client.GetOrderTypes(ctx, token, orgID)
	if err != nil {
		return fmt.Errorf("failed to get order types: %w", err)
	}
	
	var orderTypeID string
	// Loop deeply: OrderTypesResponse -> OrderTypes (groups) -> Items
	for _, group := range orderTypesResp.OrderTypes {
		for _, item := range group.Items {
			if item.Name == "Доставка БРО" {
				orderTypeID = item.ID
				break
			}
		}
		if orderTypeID != "" { break }
	}
	if orderTypeID == "" {
		// Fallback or error? Legacy returns null
		return fmt.Errorf("order type 'Доставка БРО' not found")
	}

	// 7. Construct Items
	fullMenu, err := c.client.GetProducts(ctx, token, orgID)
	if err != nil {
		return fmt.Errorf("failed to fetch menu: %w", err)
	}
	
	// Fallback to searching by NAME if GUID fails
	findIDByAny := func(idOrCode, name string) string {
		if idOrCode != "" {
			for _, p := range fullMenu {
				if p.ID == idOrCode || p.Code == idOrCode {
					return p.ID
				}
			}
		}
		// Try by name as last resort
		if name != "" {
			nameUpper := strings.ToUpper(name)
			for _, p := range fullMenu {
				pNameUpper := strings.ToUpper(p.Name)
				// 1. Exact match with prefix removal
				cleanSyrveName := strings.TrimPrefix(pNameUpper, "D ")
				if cleanSyrveName == nameUpper {
					return p.ID
				}
				// 2. Substring match for delivery products
				if strings.Contains(nameUpper, "ДОСТАВКА") && strings.Contains(pNameUpper, nameUpper) {
					log.Printf("Sub-match found for delivery: %s matched %s", name, p.Name)
					return p.ID
				}
				// 3. Fallback for "ДОСТАВКА ТЕСТ" specifically if it's named slightly differently
				if nameUpper == "ДОСТАВКА ТЕСТ" && strings.Contains(pNameUpper, "ДОСТАВКА") && strings.Contains(pNameUpper, "ТЕСТ") {
					log.Printf("Special-match found for %s: matched %s", name, p.Name)
					return p.ID
				}
			}
		}
		return ""
	}

	// Map menu for quick lookup
	menuMap := make(map[string]syrve.MenuItemDTO)
	for _, p := range fullMenu {
		menuMap[p.ID] = p
	}

	var syrveItems []syrve.OrderItem
	
	for _, item := range order.Items {
		// Resolve Main Product
		syrveProductID := findIDByAny(item.ExternalProductID, item.Name)
		
		if syrveProductID == "" {
			log.Printf("Warning: Could NOT resolve item %s (ID: %s) to GUID. Using as is.", item.Name, item.ExternalProductID)
			syrveProductID = item.ExternalProductID
		}

		// Split items into separate lines (Quantity 1) to satisfy Syrve restrictions (Commodity Marks)
		// and simplify modifier logic.
		for i := 0; i < item.Quantity; i++ {
			syrveItem := syrve.OrderItem{
				ProductID: syrveProductID,
				Amount:    1.0, 
				Price:     &item.Price,
				Type:      "Product",
				Modifiers: []syrve.OrderModifier{},
			}

			// Handle Variation/Modifier from event
			if item.ProductVariationExternalID != nil && *item.ProductVariationExternalID != "" {
				vName := ""
				if item.ProductVariationName != nil {
					vName = *item.ProductVariationName
				}
				vID := findIDByAny(*item.ProductVariationExternalID, vName)
				if vID != "" {
					vProduct, exists := menuMap[vID]
					if exists && vProduct.Type == "modifier" {
						// Find if it belongs to a group modifier for the main product
						var foundGroupID string
						defaultAmount := 1.0

						if mainProd, exists := menuMap[syrveItem.ProductID]; exists {
							// 1. Check Group Modifiers
							for _, gm := range mainProd.GroupModifiers {
								for _, cm := range gm.ChildModifiers {
									if cm.ID == vID {
										foundGroupID = gm.ID
										if cm.DefaultAmount != nil {
											defaultAmount = float64(*cm.DefaultAmount)
										} else if cm.MinAmount > 0 {
											defaultAmount = float64(cm.MinAmount)
										}
										break
									}
								}
								if foundGroupID != "" {
									break
								}
							}

							// 2. Check Simple Modifiers if not found in groups
							if foundGroupID == "" {
								for _, m := range mainProd.Modifiers {
									if m.ID == vID {
										if m.DefaultAmount != nil {
											defaultAmount = float64(*m.DefaultAmount)
										} else if m.MinAmount > 0 {
											defaultAmount = float64(m.MinAmount)
										}
										break
									}
								}
							}
						}

						log.Printf("Adding variation %s as MODIFIER to %s (GroupID: %s, Amount: %.0f)", vName, item.Name, foundGroupID, defaultAmount)
						syrveItem.Modifiers = append(syrveItem.Modifiers, syrve.OrderModifier{
							ProductID:      vID,
							ProductGroupID: foundGroupID,
							Amount:         defaultAmount,
						})
					} else if exists {
						log.Printf("Variation %s resolved to Product Type %s, REPLACING main product %s", vName, vProduct.Type, item.Name)
						syrveItem.ProductID = vID
					}
				}
			}

			// Automatic Enrichment: Add Mandatory Modifiers
			if product, exists := menuMap[syrveItem.ProductID]; exists {
				c.enrichItemWithModifiers(&syrveItem, product)
			}
			
			syrveItems = append(syrveItems, syrveItem)
		}
	}
	
	// Delivery Cost
	if order.DeliveryCost > 0 {
		productName := "ДОСТАВКА ТЕСТ"
		pID := findIDByAny("", productName)
		if pID != "" {
			syrveItems = append(syrveItems, syrve.OrderItem{
				ProductID: pID,
				Amount:    1,
				Price:     &order.DeliveryCost,
				Type:      "Product",
			})
		} else {
			log.Printf("Warning: Delivery product '%s' not found", productName)
		}
	}
	
	if order.DeliveryDoor {
		productName := "ДОСТАВКА ДО ДВЕРЕЙ"
		pID := findIDByAny("", productName)
		if pID != "" {
			// Is price fixed? Legacy `getItemPrice($product)`. 
			// We don't have price in event easily for this specific item apart from `DeliveryDoorPrice`.
			price := order.DeliveryDoorPrice
			syrveItems = append(syrveItems, syrve.OrderItem{
				ProductID: pID,
				Amount:    1,
				Price:     &price,
				Type:      "Product",
			})
		}
	}

	// 8. Construct Customer & Order
	customer := &syrve.Customer{
		Name:  order.Name,
		Phone: order.Phone,
		Type:  "regular",
	}

	payload := syrve.CreateOrderRequest{
		OrganizationID: orgID,
		TerminalID:     terminalID,
		Order: syrve.OrderPayload{
			OrderTypeID: orderTypeID,
			TableIDs:    []string{selectedTableID},
			Customer:    customer,
			Phone:       order.Phone,
			Items:       syrveItems,
			Comment:     order.Wishes, // or Address? Legacy puts address in "Delivery Address" fields usually, but here likely just table order.
			// Legacy creates a "TableOrder", so it doesn't pass address in `delivery` block because it uses `createOrder` for TABLE.
			// Legacy `orderObject`: 
			// 'orderTypeId' => $type, 'tableIds' => [$section], 'items' => ..., 'customer' => ...
			// It implies address is NOT sent to Syrve? Or maybe via specific fields?
			// The legacy code showed NO address mapping in `sendOrder`.
			// It assumes kitchen sees order, and courier takes info from elsewhere? 
			// Or waiter sees it.
			// We replicate legacy: Just items and customer on a specific table.
		},
	}
	
	// Add External ID to link them
	payload.Order.ID = order.ID.String()

	resp, err := c.client.CreateOrder(ctx, token, payload)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	log.Printf("Order created in Syrve. ID: %s", resp.OrderInfo.ID)

	return nil
}
